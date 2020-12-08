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

var ActionChan map[string](chan []string) = make(map[string](chan []string),0)

func PersonalCrawler(c *gin.Context){
	var payload map[string]interface{}
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



	workingAccount := payload["Domain"].(string)+payload["Account"].(string)
	if v,ok := ActionChan[workingAccount]; !ok {
		ActionChan[workingAccount] = make(chan []string)
		ActionChan[workingAccount] <- payload["UrlQue"].([]string)
		err := InitFacebookAccount(ActionChan[workingAccount], payload["Account"].(string), payload["Password"].(string))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Msg" : err.Error(),
				"StatusCode" : "404",
			})
			return 
		}
	} else {
		select {
		case targetQue := <- v:
			if op := payload["UrlQue"].([]string)[0]; op == "Stop"{
				
			}
			targetQue = append( targetQue, payload["UrlQue"].([]string)...)
			v <- targetQue
		case <- time.After( 30* time.Second) :
			c.JSON( http.StatusNotFound, gin.H{
				"Msg" : "Que error #1",
				"StatusCode": "404",
			})
		}
	}
}








func InitFacebookAccount( groupURL chan []string, account, password string) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Safari/537.36`),
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:], options...)
	allocCtx, Cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(logger), // todo, can set the request log
	)	
	defer cancel()

	ctx, cancel = context.WithTimeout( ctx, 300 * time.Second)
	defer cancel()

	err := chromedp.Run( ctx, FacebookLogin(account, password))
	if err != nil {
		fmt.Fprintln(gin.DefaultWriter, err.Error())
	}
	for {
		select {
		case target := <- groupURL: 
			if target[0] == "Stop" {
				return 
			}
			err = chromedp.Run( ctx, )
			if err != nil {
				fmt.Fprintln(gin.DefaultWriter, err.Error())
			}			
		}
	}
}

func FacebookLogin(account, password string) chromedp.Tasks {
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
				err = chromedp.SendKeys(`#email`, account).Do(c)
				if err != nil {
					return err
				}
				err = chromedp.Click(`#pass`).Do(c)
				if err != nil {
					return err
				}
				err = chromedp.SendKeys(`#pass`, password).Do(c)
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
	}
}

func FacebookGroupScraping( url string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(c context.Context) error {

		}),
	}
}














func fb() chromedp.Tasks{
	var attrValue 	string
	var ok 			bool

	var nodes 		[]*cdp.Node
	var pnodes		[]*cdp.Node
	var nodeListLen int = 0
	var lastLen 	int = 0
	var keep		bool = true

	// Get setting

	return chromedp.Tasks{

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
