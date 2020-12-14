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

func FacebookScraper(c *gin.Context) {
	var payload map[string]string
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

	

}