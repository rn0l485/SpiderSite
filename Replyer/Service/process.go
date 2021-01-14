package Replyer

import (
	//"fmt"
	"regexp"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"Decorations/Replyer/Config"
)

func processing( c *gin.Context ) {
	var payload map[string]string
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

	resp, err := w.Post(config.MongoDBApi+"/v1/search", false, gin.H{
		"DataBaseName" : "Spider",
		"CollectionName" : "Reply",
		"Filter" : gin.H{
			"status" : gin.H{
				"$eq" : "alive",
			},
		},
	})

	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"Msg": err.Error(),
				"StatusCode" : "500",
			},
		)
		return
	}

	var ReplySet map[string]interface{}
	if err := json.Unmarshal( resp.Body, &ReplySet); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"Msg": err.Error(),
				"StatusCode" : "500",
			},
		)
		return
	}


	targetKeyPair := gin.H{
		"Weights" : "0",
	}

	for _,v := range ReplySet["Data"].([]interface{}) {

		matched, err := regexp.MatchString( v.(map[string]interface{})["Keyword"].(string), payload["Data"])
		if err != nil {
			continue
		}
		if matched {
			if v.(map[string]interface{})["Weights"].(string) > targetKeyPair["Weights"].(string) {
				targetKeyPair["Weights"] = v.(map[string]interface{})["Weights"].(string)
				targetKeyPair["Keyword"] = v.(map[string]interface{})["Keyword"].(string)
				targetKeyPair["ReplyStatment"] = v.(map[string]interface{})["ReplyStatment"].(string)
			}
		}
	}

	if targetKeyPair["Weights"].(string) == "0" {
		c.JSON( http.StatusInternalServerError, gin.H{
			"Msg": "no keyword",
			"StatusCode" : "500",
		})
		return
	}
	targetKeyPair["Msg"] = "ok"
	targetKeyPair["StatusCode"] = "200"

	//fmt.Print(targetKeyPair)
	c.JSON(http.StatusOK, targetKeyPair)
}