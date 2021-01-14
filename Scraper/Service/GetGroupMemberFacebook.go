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

func GetGroupMemberFacebook( c *gin.Context) {
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


	if err := chromedp.Run( ctxSub, FacebookMemberGetting( (*payload.Url))); err != nil {
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

func FacebookMemberGetting( url string) chromedp.Tasks {
	return chromedp.Tasks{ 
		chromedp.Navigate(url),
		chromedp.WaitVisible(`/html/body/div[1]/div/div[1]/div[1]/div[3]/div/div/div[1]/div[1]/div[4]/div/div/div/div/div/div/div/div/div/div/div[1]`, chromedp.BySearch)
		chromedp.ActionFunc(func(c context.Context) error {

		},	
}