package Crawler

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"regexp"
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"	
	"github.com/chromedp/cdproto/input"


	"Decorations/Scraper/Service/Models"
	"Decorations/Scraper/Config"
)

func FBReplyer( c *gin.Context) {
	select {
	case FBWorking <- true:
	case <- time.After(10*time.Second):
	}
}


func fbReply() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(c context.Context) error {

			err := chromedp.Navigate(`https://www.facebook.com`).Do(c)
			if err != nil { return err }

			err = chromedp.WaitReady(`#facebook`).Do(c)
			if err != nil { return err }

			err = chromedp.AttributeValue(`#facebook`, "class", &attrValue, &ok).Do(c)
			if err != nil { return err }
			if ok && attrValue == "" {
				err = chromedp.WaitVisible(`#email`).Do(c)
				if err != nil { return err }
				err = chromedp.Click(`#email`).Do(c)
				if err != nil { return err }
				err = chromedp.SendKeys(`#email`, setting.Value.(map[string]interface{})["Account"].(string)).Do(c)
				if err != nil { return err }
				err = chromedp.Click(`#pass`).Do(c)
				if err != nil { return err }
				err = chromedp.SendKeys(`#pass`, setting.Value.(map[string]interface{})["Password"].(string)).Do(c)
				if err != nil { return err }
				err = chromedp.Click(`[name="login"]`).Do(c)
				if err != nil { return err }
				err = chromedp.WaitReady(`#facebook`).Do(c)
				if err != nil { return err }
			}
			return nil 
		}),
		chromedp.ActionFunc(func(c context.Context) error{
			resp, err := w.Post(config.MongoDBApi+"/v1/search", false, gin.H{
				"DataBaseName"	: "Spider",
				"CollectionName" : "Post",
				"Filter" : {
					"domain" : "facebook",
					"status" : "0",
					"account" : setting.Value.(map[string]interface{})["Account"].(string),
				},
			})
			if err != nil { return err }
			var respJson map[string]interface{}
			var errRecord []map[string]interface{}
			if err := json.Unmarshal( resp, &respJson); err != nil { return  err }
			if respJson["StatusCode"].(string) == "200" {
				for _, v := range respJson["Data"].([]map[string]interface{}){
					err = chromedp.Navigate(v["Url"].(string)).Do(c)
					if err != nil {
						errRecord = append(errRecord, v)
						continue
					}			
					err = chromedp.WaitVisible(`div[role="textbox"]`, chromedp.ByQuery).Do(c)
					if err != nil { 
						errRecord = append(errRecord, v)
						continue
					}						

					err = chromedp.Click(`div[role="textbox"]`, chromedp.ByQuery).Do(c)
					if err != nil {
						errRecord = append(errRecord, v) 
						continue
					}


					rs, err := w.Post(config.MongoDBApi+"/v1/search", false, gin.H{
						"DataBaseName"	: "Spider",
						"CollectionName" : "Reply",
						"Filter" : {
							"_id" : v["ReplyID"].(int),
						},
					})
					if err != nil {
						errRecord = append(errRecord, v) 
						continue
					}
					var rsjson map[string]interface{}
					if err := json.Unmarshal( rs, &rsjson); err != nil { 
						errRecord = append(errRecord, v) 
						continue
					} 




					err = input.DispatchKeyEvent(input.KeyChar).WithText(rsjson.ReplyStatment.(string)).Do(c)
					if err != nil { 
						errRecord = append(errRecord, v) 
						continue 
					}

					err = input.DispatchKeyEvent(input.KeyRawDown).WithWindowsVirtualKeyCode(13).Do(c)
					if err != nil {
						errRecord = append(errRecord, v) 
						continue
					}
					_ = chromedp.Sleep( 3 * time.Second)
				}
			} else {
				return error.New("No post received")
			}
			for _, v := range errRecord {
				_, err := w.Post(config.MongoDBApi+"/v1/modify", false, gin.H{
					"DataBaseName"	: "Spider",
					"CollectionName" : "Post",
					"Filter" : {
						"_id" : v["Id"].(int),
					},
					"ChangeField" : {
						"status" : "-1",
					}
				})
			}
			return nil
		}),
	}
}

