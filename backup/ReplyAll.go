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

func FacebookReplyAll(c *gin.Context) {
	var payload map[string]string
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{
				"Msg": err.Error(),
				"StatusCode" : "404",
			},
		)
		return
	}

	resp, err := w.Post( config.MongoDBApi+"/v1/search", false, gin.H{
		"DataBaseName":"Spider",
		"CollectionName":"Post",
		"Filter" : gin.H{
			"status" 	: gin.H{
				"$eq" 	: "unreply",
			},
			"Domain" 	: gin.H{
				"$eq"	: "facebook",
			},
		},
	})

	if err != nil {
		fmt.Fprintln(gin.DefaultWriter, err.Error())
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{
				"Msg": "DataBase Error",
				"StatusCode" : "500",
			},
		)
		return
	}
	var respJson ScraperModels.Resp
	if err := json.Unmarshal( resp.Body, &respJson); err != nil {
		fmt.Fprintln(gin.DefaultWriter, err.Error())
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg": "Error",
			"StatusCode" : "404",
		})
		return
	}

	if respJson.Msg != "200" {
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg" : respJson.Msg,
			"StatusCode": respJson.StatusCode,
		})
		return
	}

	if len(respJson.Data) == 0 {
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg" : "No data",
			"StatusCode":"404",
		})
		return		
	}

	questionQue := map[string][]string{}
	ErrorIds := map[int]string
	for i,v := range respJson.Data {
		respReply, err := w.Post( config.ReplyApi + "/v1/do", false, gin.H{
			"Data" : v.(map[string]interface{})["context"].(string), 							// back with keyword
		})
		if err != nil {
			ErrorIds[v["_id"].(int)] = err.Error()
			continue
		}

		var callback map[string]string
		if err := json.Unmarshal( respReply.Body, &callback); err != nil {
			ErrorIds[v["_id"].(int)] = err.Error()
			continue
		}

		if callback["StatusCode"] != "200"{
			ErrorIds[v["_id"].(int)] = "code error"
			continue
		}



		if _,ok := ActionChan["facebook"+v["ReplyAccount"].(string)]; !ok {
			ErrorIds[v["_id"].(int)] = "no goroutine"
			continue			
		} 

		if qq,ok:= questionQue["facebook"+v["ReplyAccount"].(string)]; !ok{
			questionQue["facebook"+v["ReplyAccount"].(string)] = []string{"Reply", ""}
		}


		messionQue := []string{"Reply"}
		for 
		messionQue = append( messionQue, )








		ActionChan["facebook"+v["ReplyAccount"].(string)] <- 







		respReslut, err := w.Post( config.MongoDBApi + "/v1/modify", false, gin.H{
			"DataBaseName" : "Spider",
			"CollectionName" : "Post",
			"Filter" : gin.H{
				"_id" : gin.H{
					"$eq" 	:  v["_id"].(int),
				},
			},
			"ChangeField" : gin.H{
				"Status" 		: "alive",
				"ReplyKeyWord" 	:  callback["Key"],
			}
		})
	}
}