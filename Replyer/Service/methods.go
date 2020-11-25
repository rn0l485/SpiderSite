package Replyer

import (
	"regexp"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"Decorations/Replyer/Service/Models"
	"Decorations/Replyer/Func/Worker"
	"Decorations/Replyer/Config"
)

func processing( c *gin.Context ) {

	resp, err := w.Post(config.MongoDBApi+"/v1/search", false, gin.H{
		"DataBaseName" : "Spider",
		"CollectionName" : "Post",
		"Filter" : {
			"domain" : "facebook",
			"status" : "0",
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
	var PostSet map[string]interface{}
	if err := json.Unmarshal( resp.Body, &PostSet); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"Msg": err.Error(),
				"StatusCode" : "500",
			},
		)
		return
	}

	resp, err = w.Post(config.MongoDBApi+"/v1/search", false, gin.H{
		"DataBaseName" : "Spider",
		"CollectionName" : "Reply",	
		"Filter" :{},	
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

	
	keep := true
	for _, reply := range ReplySet["Data"].([]models.Reply) {
		for _, post := range PostSet["Data"].([]models.Post) {
			matched, err := regexp.MatchString(reply.Keyword, post.Context)
			if err != nil {
				continue
			}
			if matched {

				keep = false
				break
			}
		}
		if !keep {
			
			break
		}
		
	}
}