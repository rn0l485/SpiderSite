package Replyer

import (
	"fmt"
	"regexp"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"Decorations/Replyer/Service/Models"
	"Decorations/Replyer/Func/Worker"
	"Decorations/Replyer/Config"
)

func add(c *gin.Context) {
	var payload models.Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"Msg": err.Error(),
				"StatusCode" : "500",
			},
		)
		return
	}

	var respJson map[string]string
	resp, err := w.POST(config.MongoDBApi+"/v1/add", false, gin.H{
		"DataBaseName" : "Spider",
		"CollectionName" : "Reply",
		"Record" : models.Reply{
			"Keyword" 		: 	payload.Setting["Keyword"],
			"ReplyStatment" :	payload.Setting["ReplyStatment"],
			"Weights"		: 	payload.Setting["Weights"],
		}
	})
	if err != nil {
		fmt.Fprintln(gin.DefaultWriter, err.Error())
		return
	}

	if err := json.Unmarshal( resp.Body, &respJson); err != nil {
		fmt.Fprintln(gin.DefaultWriter, err.Error())
		return
	}

	if respJson["StatusCode"] != "200" {
		fmt.Fprintln(gin.DefaultWriter, err.Error())
		return
	}
}