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


	//"Decorations/Scraper/Service/Models"
	"Decorations/Scraper/Config"
)

var FacebookActionChan map[string](chan []string) = make(map[string](chan []string))

func FacebookInit(c *gin.Context) {
	var payload map[string]string
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

	if _, ok := FacebookActionChan[payload["Domain"]+payload["Account"]]; ok {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{
				"Msg" : "#1 Account exsit",
				"StatusCode": "404",
			},
		)
		return
	}

	FacebookActionChan[payload["Domain"]+payload["Account"]] = make(chan []string, 1)
	if _,ok := FacebookActionChan[payload["Domain"]+payload["Account"]+"ErrorGate"]; !ok {
		FacebookActionChan[payload["Domain"]+payload["Account"]+"ErrorGate"] = make(chan []string, 1)
	} else {
		<- FacebookActionChan[payload["Domain"]+payload["Account"]+"ErrorGate"]
	}


	go InitFacebookAccount(FacebookActionChan[payload["Domain"]+payload["Account"]], FacebookActionChan[payload["Domain"]+payload["Account"]+"ErrorGate"], payload["Account"], payload["Password"])
	

	c.JSON(http.StatusOK, gin.H{
		"Msg"			: "ok",
		"StatusCode" 	: "200",
	})
}

func InitFacebookAccount( groupURL, accountErrorGate chan []string, account, password string) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("disable-notifications", true),
		chromedp.Flag("mute-audio", false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Safari/537.36`),
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:], options...)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	// create chrome instance
	ctxA, cancelA := chromedp.NewContext( allocCtx )	
	defer cancelA()

	ctx, cancel := context.WithTimeout( ctxA, 600 * time.Second )
	defer cancel()

	err := chromedp.Run( ctx, FacebookLogin(account, password))
	if err != nil {
		fmt.Fprintln(gin.DefaultWriter, "#0 "+err.Error())
		delete(FacebookActionChan, "facebook"+account )

		return
	}
	fmt.Fprintln(gin.DefaultWriter, "facebook"+account)

	for {
		select {
		case target := <- groupURL: 
			if target[0] == "Stop" {
				fmt.Fprintln(gin.DefaultWriter, "#1 Stop the account")
				fmt.Fprintln(gin.DefaultWriter, "facebook over "+account)
				delete(FacebookActionChan, "facebook"+account )
				return
			} else if target[0] == "Group" {
				for _,v := range target[1:]{
				
					err = chromedp.Run( ctx, FacebookGroupScraping(v, account) )
					if err != nil {
						if eMsg := err.Error(); eMsg == "#1"{
							fmt.Fprintln(gin.DefaultWriter, eMsg+"Group error : "+v)							
							<- accountErrorGate
							continue							
						}
						fmt.Fprintln(gin.DefaultWriter, "#2 "+err.Error())
						fmt.Fprintln(gin.DefaultWriter, "facebook over "+account+" // "+v)
						delete(FacebookActionChan, "facebook"+account ) 
						<- accountErrorGate
						continue
					}
				}
				<- accountErrorGate
			} else if target[0] == "Reply" {				
				for i,_ := range target[1:]{

					if i%2 == 0 {
						continue
					}
					err = chromedp.Run( ctx, FacebookReply(target[i],target[i+1]))
					if err != nil {
						if eMsg := err.Error(); eMsg == "#1"{
							fmt.Fprintln(gin.DefaultWriter, eMsg+"Group error : "+ target[i])							
							<- accountErrorGate
							continue							
						}
						fmt.Fprintln(gin.DefaultWriter, "#3 "+err.Error())
						fmt.Fprintln(gin.DefaultWriter, "facebook over "+account)
						delete(FacebookActionChan, "facebook"+account )
						<- accountErrorGate
						continue
					}	
								
				}
				<- accountErrorGate
			} 
		}
	}
	return
}

func FacebookLogin(account, password string) chromedp.Tasks {
	var attrValue 	string
	var ok 			bool
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			c := chromedp.FromContext(ctx)

			err := chromedp.Navigate(`https://www.facebook.com`).Do(cdp.WithExecutor(ctx, c.Target))
			if err != nil {
				fmt.Fprintln(gin.DefaultWriter, err.Error())
				return err
			}

			err = chromedp.WaitReady(`#facebook`).Do(cdp.WithExecutor(ctx, c.Target))
			if err != nil {
				fmt.Fprintln(gin.DefaultWriter, err.Error())
				return err
			}

			err = chromedp.AttributeValue(`#facebook`, "class", &attrValue, &ok).Do(cdp.WithExecutor(ctx, c.Target))
			if err != nil {
				fmt.Fprintln(gin.DefaultWriter, err.Error())
				return err
			}
			if ok && attrValue == "" {

				err = chromedp.WaitVisible(`#email`).Do(cdp.WithExecutor(ctx, c.Target))
				if err != nil {
					fmt.Fprintln(gin.DefaultWriter, err.Error())
					return err
				}
				err = chromedp.Click(`#email`).Do(cdp.WithExecutor(ctx, c.Target))
				if err != nil {
					fmt.Fprintln(gin.DefaultWriter, err.Error())
					return err
				}
				err = chromedp.SendKeys(`#email`, account).Do(cdp.WithExecutor(ctx, c.Target))
				if err != nil {
					fmt.Fprintln(gin.DefaultWriter, err.Error())
					return err
				}
				err = chromedp.Click(`#pass`).Do(cdp.WithExecutor(ctx, c.Target))
				if err != nil {
					fmt.Fprintln(gin.DefaultWriter, err.Error())
					return err
				}
				err = chromedp.SendKeys(`#pass`, password).Do(cdp.WithExecutor(ctx, c.Target))
				if err != nil {
					fmt.Fprintln(gin.DefaultWriter, err.Error())
					return err
				}
				err = chromedp.Click(`[name="login"]`).Do(cdp.WithExecutor(ctx, c.Target))
				if err != nil {
					fmt.Fprintln(gin.DefaultWriter, err.Error())
					return err
				}
				err = chromedp.WaitReady(`#facebook`).Do(cdp.WithExecutor(ctx, c.Target))
				if err != nil {
					fmt.Fprintln(gin.DefaultWriter, err.Error())
					return err
				}
			}
			return nil 
		}),
		chromedp.WaitReady(`#facebook`),
	}
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

func FacebookReply( url, reply string) chromedp.Tasks {
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
