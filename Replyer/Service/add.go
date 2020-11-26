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


	resp, err := w.POST(config.MongoDBApi+"/v1/add", false, gin.H{
		"DataBaseName" : "Spider",
		"CollectionName" : "Reply",
		"Record" : Reply{
			"Keyword" 		: 	payload	
			"ReplyStatment" :	payload
			"Weights"		: 	payload
		}
	})

	
}