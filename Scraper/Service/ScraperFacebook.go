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

func FacebookScraper( c *gin.Context) {
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

	if payload.Account == nil || payload.Domain == nil || payload.Url == nil {
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


	if err := chromedp.Run( ctxSub, FacebookGroupScraping( (*payload.Url), (*payload.Account))); err != nil {
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

func FacebookGroupScraping( url, account string)  chromedp.Tasks {
	var nodes 		[]*cdp.Node
	var pnodes		[]*cdp.Node
	var nodeListLen int = 0
	var lastLen 	int = 0
	var keep		bool = true
	var GroupName 	string

	return chromedp.Tasks{
		chromedp.ActionFunc(func(c context.Context) error{
			respForGroupName, err := w.Post( config.MongoDBApi+"/v1/search", false, gin.H{
				"DataBaseName"	: "Spider",
				"CollectionName" : "UserGroup",
				"Filter" : gin.H{
					"status": gin.H{
						"$eq":"alive",
					},
					"Url" :	gin.H{
						"$eq" : url,
					},
				},	
			})

			if err != nil {
				return err
			}
			var respForGroupNameJson map[string]interface{}
			if err := json.Unmarshal( respForGroupName.Body, &respForGroupNameJson); err != nil { return  err }
			if respForGroupNameJson["StatusCode"].(string) != "200" {
				return errors.New("saving error: "+respForGroupNameJson["Msg"].(string))
			}
			if respForGroupNameJson["Data"] == nil {
				return errors.New("No Group Data")
			}

			if gn, ok := respForGroupNameJson["Data"].([]interface{})[0].(map[string]interface{})["GroupName"].(string); ok{
				GroupName = gn
				return nil
			} else {
				return errors.New("parsing error")
			}
		}),
		chromedp.ActionFunc(func(c context.Context) error {
			err := chromedp.Navigate(`https://www.facebook.com`).Do(c)
			if err != nil {
				return errors.New("#1")
			}

			err = chromedp.WaitReady(`#facebook`).Do(c)
			if err != nil {
				return errors.New("#1")
			}			

			err = chromedp.Navigate(url).Do(c)
			if err != nil {
				return errors.New("#1")
			}

			err = chromedp.WaitReady(`div[role="feed"]`, chromedp.ByQuery).Do(c)
			if err != nil {
				return errors.New("#1")
			}

			_ = chromedp.Sleep( 3 * time.Second).Do(c)
			return nil		
		}),
		chromedp.ActionFunc(func(c context.Context) error{
			stratTime := time.Now().Unix()


			_, exp, err := runtime.Evaluate(`window.scrollTo(0,document.body.scrollHeight);`).Do(c)
			if err != nil { return err }
			if exp != nil { return exp }

			for true {
				err = chromedp.WaitReady(`div[role="feed"]`, chromedp.ByQuery).Do(c)
				if err != nil { return err }

				nodes = nil 
				err = chromedp.Nodes(`//*[@id="mount_0_0"]/div/div[1]/div[1]/div[3]/div/div/div[1]/div[1]/div[4]/div/div/div/div/div[1]/div[2]/div/div`, &nodes, chromedp.BySearch).Do(c)
				if err != nil { return err }
				lastLen = len(nodes)
				err = chromedp.Nodes(`/html/body/div[1]/div/div[1]/div[1]/div[3]/div/div/div[1]/div[1]/div[4]/div/div/div/div/div[1]/div[2]/div/div/div/div/div/div/div/div/div/div/div/div[2]/div/div[2]/div/div[2]/div/div[2]/span/span/span[2]/span/a`, &pnodes,  chromedp.BySearch).Do(c)
				if  err != nil { return err }	
				for i, node := range pnodes {
					if i < nodeListLen {
						keep = false
						continue
					}
					err = dom.ScrollIntoViewIfNeeded().WithNodeID(node.NodeID).Do(c)
					if err != nil { err.Error() }
					_ = chromedp.Sleep(3 * time.Second)
					err = MouseMoveNode(node).Do(c)
					if err != nil { 
						fmt.Printf(err.Error() )
						continue
					}
						
					
					nodeListLen = i
					keep = true
				}
				
				if len(nodes)-3 > config.SearchLimitationPostNum { break }
				if nowTime := time.Now().Unix(); nowTime - stratTime > 180 { break }
				if keep {
					_, exp, err := runtime.Evaluate(`window.scrollTo(0,document.body.scrollHeight);`).Do(c)
					if err != nil { return err }
					if exp != nil { return exp }
				}
			}
			return nil 
		}),	
		chromedp.Sleep( 5 * time.Second ),
		chromedp.ActionFunc(func(c context.Context) error {
			rPost, _ := regexp.Compile(`href=\"https://www.facebook.com/groups/([0-9]+?)/permalink/([\s\S]+?)/\?`)
			rUser, _ := regexp.Compile(`href=\"/groups/([0-9]+?)/user/([\s\S]+?)\"`)

			rClass, _ := regexp.Compile(`class=\"([\s\S]+?)\"`)

			reCaptial, _ := regexp.Compile("\\<[\\S\\s]+?\\>")


			//去除 STYLE
			reStyle, _ := regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
    		//去除 SCRIPT
			reScript, _ := regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
		
			for _, node := range nodes {			
				innerHtml, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(c)					
				if err != nil {	continue }
				postURL := rPost.FindString(innerHtml)
				userURL := rUser.FindString(innerHtml)
				innerHtml = reCaptial.ReplaceAllStringFunc(innerHtml, 	strings.ToLower)
				innerHtml = rClass.ReplaceAllString( innerHtml, 	"")
				innerHtml = reStyle.ReplaceAllString(innerHtml, 	"")
				innerHtml = reScript.ReplaceAllString(innerHtml, 	"")

				if postURL == "" { 
					continue
				}


				respKey, err := w.Post( config.ReplyApi+ "/v1/do", false, gin.H{
					"Data" : innerHtml,
				})
				if err != nil {
					return err
				}
				var respKeyJson map[string]string
				keywordReplying := ""
				statmentReplying := ""
				if err := json.Unmarshal( respKey.Body, &respKeyJson); err != nil { return  err }
				if respKeyJson["StatusCode"] != "200" {
					if respKeyJson["Msg"] != "no keyword" {
						return errors.New("error: "+respKeyJson["Msg"])
					}
				} else {
					keywordReplying = respKeyJson["Keyword"]
					statmentReplying = respKeyJson["ReplyStatment"]
				}

				newPost := gin.H{
					"Domain"		: "facebook",
					"Group"			: url,
					"Url"			: postURL,
					"GroupName" 	: GroupName,
					"ClientUrl"		: userURL,
					//"Context" 		: innerHtml,
					"ReplyKeyword"	: keywordReplying, 
					"ReplyStatment" : statmentReplying,
					"CreateTime" 	: strconv.FormatInt(time.Now().Unix(), 10),
					"ReplyAccount" 	: account,
					"status" 		: "unreply",
				}

				resp, err := w.Post( config.MongoDBApi+"/v1/add", false, gin.H{
					"DataBaseName"	: "Spider",
					"CollectionName" : "Post",
					"Record" : newPost,
				})
				if err != nil { return err }

				var respJson map[string]string
				if err := json.Unmarshal( resp.Body, &respJson); err != nil { return  err }
				if respJson["StatusCode"] != "200" {
					return errors.New("saving error: "+respJson["Msg"])
				}

			}
			return nil 
		}),
	}
}
