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


var (
	FacebookActionChan map[string]ScraperModels.Page = make(map[string]ScraperModels.Page)
)

func FacebookInit(c *gin.Context) {
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

	if payload.Account == nil || payload.Password == nil || payload.Domain == nil {
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg" : "#1 No basic info",
			"StatusCode":"404",
		})
		return
	}

	if _, ok := FacebookActionChan[(*payload.Domain)+(*payload.Account)]; ok {
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg":"#2 Account exsit",
			"StatusCode":"404",
		})
		return
	}

	ctx, cancel := context.WithCancel( allocCtx )
	ctx, _ := chromedp.NewContext(ctx)	

	
	if err := chromedp.Run( ctx, FacebookLogin( (*payload.Account), (*payload.Password))); err != nil {
		cancel()
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg":"#3 "+err.Error(),
			"StatusCode":"404",
		})
		return
	}

	FacebookActionChan[(*payload.Domain)+(*payload.Account)] := ScraperModels.Page{
		CTX : ctx,
		Cancel : cancel,
	}

	c.JSON( http.StatusOK, gin.H{
		"Msg" : "ok",
		"StatusCode":"200",
	})
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