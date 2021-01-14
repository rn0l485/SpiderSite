package Crawler

import (
	"fmt"
	//"log"
	"strings"
	"context"
	"time"
	"net/http"
	"regexp"
	"errors"
	"strconv"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"	
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/cdproto/runtime"


	"Decorations/Scraper/Service/Models"
	"Decorations/Scraper/Config"
)

func FacebookReplyer( c *gin.Context) {
	var payload ScraperModels.Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{
				"Msg": "#0" + err.Error(),
				"StatusCode" : "404",
			},
		)
		return
	}

	if payload.Account == nil || payload.Domain == nil || payload.Url == nil || payload.Msg == nil {
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg" : "#1 No basic info",
			"StatusCode":"404",
		})
		return
	}	

	if _,ok := FacebookActionChan[(*payload.Domain)+(*payload.Account)]; !ok {
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg" : "#2 No Account",
			"StatusCode": "404",
		})
		return
	}


	ctxSub, cancel() := chromedp.NewContext(FacebookActionChan[(*payload.Domain)+(*payload.Account)].CTX)
	defer cancel()
	ctxSub, _ = context.WithTimeout( ctxSub, 300 * time.Second)


	if err := chromedp.Run( ctxSub, FacebookReplying( (*payload.Url), (*payload.Msg))); err != nil {
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg":"#3 "+err.Error(),
			"StatusCode":"404",
		})
		return
	}

	c.JSON( http.StatusOK, gin.H{
		"Msg" : "ok",
		"StatusCode":"200",
	})
	return
}


func FacebookReplying( url, reply string) chromedp.Tasks {
	return chromedp.Tasks{ 
		chromedp.ActionFunc(func(c context.Context) error {
			err := chromedp.Navigate(url).Do(c)
			if err != nil {
				return errors.New("#1A" + err.Error())
			}

			err = chromedp.WaitVisible(`div[aria-label="留言"]`, chromedp.ByQuery).Do(c)
			if err != nil {
				return errors.New("#1B" + err.Error())
			}

			_ = chromedp.Sleep( 1 * time.Second)	

			err = chromedp.Click(`div[aria-label="留言"]`, chromedp.ByQuery).Do(c)
			if err != nil {
				return errors.New("#1C" + err.Error())
			}

			_ = chromedp.Sleep( 1 * time.Second)

			err = input.DispatchKeyEvent(input.KeyChar).WithText(reply).Do(c)
			if err != nil {
				return err
			}

			_ = chromedp.Sleep( 1 * time.Second)

			err = input.DispatchKeyEvent(input.KeyRawDown).WithWindowsVirtualKeyCode(13).Do(c)
			if err != nil {
				return err
			}

			_ = chromedp.Sleep( 3 * time.Second)

			respOK, err := w.Post( config.MongoDBApi+"/v1/modify", false, gin.H{
				"DataBaseName" : "Spider",
				"CollectionName" : "Post",
				"Filter" : gin.H{
					"Url" : gin.H{
						"$eq" : url,
					},
					"status" : gin.H{
						"$eq" : "unreply",
					},
				},
				"ChangeField" : gin.H{
					"status" : "alive",
				},
			})
			if err != nil {
				return err
			}

			var respJson map[string]string
			if err := json.Unmarshal( respOK.Body, &respJson); err != nil {
				return err
			}

			if respJson["StatusCode"] != "200" {
				return errors.New("#reply saving error")
			}

			return nil 
		}),
	}
}
