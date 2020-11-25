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


var (
	FBWorking chan bool = make(chan bool)
)


func FBCrawler(c *gin.Context) {
	select {
	case FBWorking <- true:
		// todo, to start the work

		resp, err := w.Post(config.MongoDBApi+"/v1/search", false, gin.H{
			"DataBaseName"	: "Spider",
			"CollectionName" : "Config",
			"Filter" : {
				"name":"FACEBOOK_BASIC_INFO",
				"comments":"recent",
			},
		}) 
		if err != nil { 
			panic(err)
		}
		if err := json.Unmarshal( resp, &setting); err != nil {
			panic( err )
		}





		err := chromedp.Run( CTX,fb())
		if err != nil { 
			<- FBWorking
			log.Fatal(err)
		}
	case <- time.After(10*time.Second):
		c.JSON(http.StatusOK, gin.H{
			"Msg":"busy",
			"StatusCode":"200",
		})
	}
}

func fb() chromedp.Tasks{
	var attrValue 	string
	var ok 			bool
	//var innerHtml 	string	
	var nodes 		[]*cdp.Node
	var pnodes		[]*cdp.Node
	var nodeListLen int = 0
	var lastLen 	int = 0
	var keep		bool = true

	// Get setting

	/*var account 	string = "rm0l485@hotmail.com.tw"
	var password 	string = "ChengHuan@1992"
	var GroupURL 	string = "https://www.facebook.com/groups/752971828086206?sorting_setting=CHRONOLOGICAL"
	var SearchLimitationLastPost 	string = "https://www.facebook.com/groups/DataScienceGroup/permalink/3726016267460236/"
	var SearchLimitationPostNum 	int = 50*/


	return chromedp.Tasks{
		chromedp.ActionFunc(func(c context.Context) error {

			err := chromedp.Navigate(`https://www.facebook.com`).Do(c)
			if err != nil {
				return err
			}

			err = chromedp.WaitReady(`#facebook`).Do(c)
			if err != nil {
				return err
			}

			err = chromedp.AttributeValue(`#facebook`, "class", &attrValue, &ok).Do(c)
			if err != nil {
				return err
			}
			if ok && attrValue == "" {

				err = chromedp.WaitVisible(`#email`).Do(c)
				if err != nil {
					return err
				}
				err = chromedp.Click(`#email`).Do(c)
				if err != nil {
					return err
				}
				err = chromedp.SendKeys(`#email`, setting.Value.(map[string]interface{})["Account"].(string)).Do(c)
				if err != nil {
					return err
				}
				err = chromedp.Click(`#pass`).Do(c)
				if err != nil {
					return err
				}
				err = chromedp.SendKeys(`#pass`, setting.Value.(map[string]interface{})["Password"].(string)).Do(c)
				if err != nil {
					return err
				}
				err = chromedp.Click(`[name="login"]`).Do(c)
				if err != nil {
					return err
				}
				err = chromedp.WaitReady(`#facebook`).Do(c)
				if err != nil {
					return err
				}
			}
			return nil 
		}),
		chromedp.WaitReady(`#facebook`),
		chromedp.ActionFunc(func(c context.Context) error {

			err := chromedp.Navigate(setting.Value.(map[string]interface{})["GroupURL"].(string)).Do(c)
			if err != nil {
				return err
			}
			err = chromedp.WaitReady(`div[role="feed"]`, chromedp.ByQuery).Do(c)
			if err != nil {
				return err
			}
			_ = chromedp.Sleep( 3 * time.Second).Do(c)
			return nil
		}),

		chromedp.ActionFunc(func(c context.Context) error {
			

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
				

				if  len(nodes)-3 > setting.Value.(map[string]interface{})["SearchLimitationPostNum"].(int) { break }
				if keep {
					_, exp, err := runtime.Evaluate(`window.scrollTo(0,document.body.scrollHeight);`).Do(c)
					if err != nil { return err }
					if exp != nil { return exp }
				}
			}
			return nil 
		}),	
		chromedp.ActionFunc(func(c context.Context) error {
			rPost, _ := regexp.Compile(`href=\"https://www.facebook.com/groups/([0-9]+?)/permalink/([\s\S]+?)\"`)
			rUser, _ := regexp.Compile(`href=\"/groups/([0-9]+?)/user/([\s\S]+?)\"`)

			for i, node := range nodes {			
				innerHtml, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(c)					
				if err != nil {	return err }

				postURL := rPost.FindString(innerHtml)
				userURL := rUser.FindString(innerHtml)
				if postURL == setting.Value.(map[string]interface{})["SearchLimitationLastPost"].(string)  { break }

				newPost := Post{
					Domain: 		"facebook",
					Group:			setting.Value.(map[string]interface{})["GroupURL"].(string),
					Url: 			postURL,
					ClientUrl: 		userURL,
					Context: 		innerHtml,
					CreateTime: 	time.Now(),
					ReplyAccount: 	setting.Value.(map[string]interface{})["Account"].(string),
					Status: 		"0",
				}

				resp, err := w.Post( config.MongoDBApi+"/v1/add", false, gin.H{
					"DataBaseName"	: "Spider",
					"CollectionName" : "Post",
					"Record" : newPost,
				})
				if err != nil { return err }

				var respJson map[string]string
				if err := json.Unmarshal( resp, &respJson); err != nil { return  err }
				if respJson["StatusCode"] == "500" {
					return errors.New("saving error")
				}

			}
			return nil 
		}),
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

