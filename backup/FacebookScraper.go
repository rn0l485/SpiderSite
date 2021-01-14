package Crawler

import (
	//"fmt"
	//"log"
	"time"
	"net/http"
	//"regexp"
	//"errors"
	//"encoding/json"

	"github.com/gin-gonic/gin"

	//"github.com/chromedp/chromedp"
	//"github.com/chromedp/cdproto/cdp"
	//"github.com/chromedp/cdproto/dom"	
	//"github.com/chromedp/cdproto/input"


	"Decorations/Scraper/Service/Models"
	//"Decorations/Scraper/Config"
)

func FacebookScraper(c *gin.Context) {
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

	if _, ok := FacebookActionChan["facebook"+(*payload.Account)] ; !ok{
		c.JSON( http.StatusNotFound, gin.H{
			"Msg" : "#2 No account",
			"StatusCode" : "404",
		})
		return
	} 

	select {
	case FacebookActionChan["facebook"+(*payload.Account)+"ErrorGate"] <- []string{"ok"} :
	case <- time.After(5*time.Second):
		c.JSON( http.StatusNotFound, gin.H{
			"Msg" : "#2 still working",
			"StatusCode" :"404",
		})
		return
	}



	actionPayload := []string{
		(*payload.Method),
	}

	actionPayload = append(actionPayload, (*payload.Data)...)

	select{
	case FacebookActionChan["facebook"+(*payload.Account)] <- actionPayload :
		c.JSON(http.StatusOK, gin.H{
			"Msg" : "ok",
			"StatusCode":"200",
		})
		return
	case <- time.After( 10 * time.Second):
		c.JSON(http.StatusNotFound, gin.H{
			"Msg" : "action chan blocked",
			"StatusCode" : "404",
		})
		return
	}
}