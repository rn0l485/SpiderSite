package Crawler

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"regexp"
	"errors"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"	
	"github.com/chromedp/cdproto/input"


	"Decorations/Scraper/Service/Models"
	"Decorations/Scraper/Config"
)

func FacebookReplyAll (c *gin.Context) {
	var payload ScraperModels.Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{
				"Msg": "#0 " + err.Error(),
				"StatusCode" : "404",
			},
		)
		return
	}
	if payload.Account == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Msg" : "#1 Basic info",
			"StatusCode" : "404",
		})
		return
	} else if payload.Method == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Msg" : "#1 Basic info",
			"StatusCode" : "404",
		})
		return	
	}
	
	resp, err := w.Post( config.MongoDBApi+"/v1/search", false, gin.H{
		"DataBaseName":"Spider",
		"CollectionName":"Post",
		"Filter" : gin.H{
			"Domain": gin.H{
				"$eq": "facebook",
			},
			"status": gin.H{
				"$eq": "unreply",
			},
		},
	})

	if err != nil {
		c.JSON( http.StatusNotFound, gin.H{
			"Msg" : "#2 searching erro",
			"StatusCode":"404",
		})
		return 
	}

	var respJson ScraperModels.Resp
	if err := json.Unmarshal( resp.Body, &respJson); err != nil {
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg": "#3 Error",
			"StatusCode" : "404",
		})
		return
	}

	if respJson.StatusCode != "200" {
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg": "#4 Error",
			"StatusCode" : "404",
		})
		return		
	}

	dataSet := map[string][]string{}
	for _,record := range respJson.Data {
		if rc, ok := record.(map[string]interface{}); !ok {
			c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
				"Msg": "#5",
				"StatusCode" : "404",
				"Data" : rc,
			})
			return			
		} else {
			resp, err := w.Post( config.ReplyApi+ "/v1/do", false, gin.H{
				"Data" : rc["Context"].(string),
			})
			if err != nil {
				c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
					"Msg": "#6",
					"StatusCode" : "404",
				})
				return					
			}
			var respJsonReply map[string]string
			if err := json.Unmarshal( resp.Body, &respJsonReply); err != nil {
				c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
					"Msg": "#7 "+err.Error(),
					"StatusCode" : "404",
				})
				return
			}
			if respJsonReply["StatusCode"] != "200"{
				c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
					"Msg": "#7 Error",
					"StatusCode" : "404",
				})
				return					
			}

			if _, ok := dataSet["facebook"+rc["ReplyAccount"].(string)]; !ok {

				dataSet["facebook"+rc["ReplyAccount"].(string)] := []string{"Reply", rc["Url"].(string), respJsonReply["ReplyStatment"]} 
			} else {
				dataSet["facebook"+rc["ReplyAccount"].(string)] = append(dataSet["facebook"+rc["ReplyAccount"].(string)], rc["Url"].(string), respJsonReply["ReplyStatment"])
			}
		}
	}

	for k,v := range dataSet {
		if _, ok := FacebookActionChan[k]; !ok {
			c.JSON(http.StatusNotFound, gin.H{
				"Msg" : "#8 no account: "+ k,
				"StatusCode":"404",
			})
			return
		} 
		FacebookActionChan[k] <- v
		time.Sleep(2*time.Second)
	}

	c.JSON( http.StatusOK, gin.H{
		"Msg" : "ok",
		"StatusCode":"200",
	})
}