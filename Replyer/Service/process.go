package Replyer

import (
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


	targetKeyPair := map[string]string{
		"Weights" : "0",
	}

	for _,v := range ReplySet["Data"].([]interface{}) {

		matched, err := regexp.MatchString( v.(map[string]interface{})["Keyword"].(string), payload["Data"])
		if err != nil {
			continue
		}
		if matched {
			if v.(map[string]interface{})["Weights"].(string) > targetKeyPair["Weights"] {
				targetKeyPair["Weights"] = v.(map[string]interface{})["Weights"].(string)
				targetKeyPair["Keyword"] = v.(map[string]interface{})["Keyword"].(string)
				targetKeyPair["ReplyStatment"] = v.(map[string]interface{})["ReplyStatment"].(string)
			}
		}
	}

	if _,ok := targetKeyPair["Keyword"] ; !ok {
		c.JSON( http.StatusInternalServerError, gin.H{
			"Msg": "no keyword",
			"StatusCode" : "500",
		})
	}
	targetKeyPair["Msg"] = "ok"
	targetKeyPair["StatusCode"] = "200"

	c.JSON(http.StatusOK, targetKeyPair)
}