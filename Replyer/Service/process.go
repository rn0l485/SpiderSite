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
	replyPair := make(map[int]int)
	for ri, reply := range ReplySet["Data"].([]models.Reply) {
		for pi, post := range PostSet["Data"].([]models.Post) {
			matched, err := regexp.MatchString(reply.Keyword, post.Context)
			if err != nil {
				fmt.Fprintln(gin.DefaultWriter, err.Error())
				continue
			}
			if matched {
				if val, ok := replyPair[post.Id]; ok {
					if ReplySet["Data"][ri].Weights > ReplySet["Data"][val].Weights {
						replyPair[post.Id] = ri
					}
				} else {
					replyPair[post.Id] = ri
				}
			}
		}		
	}

	errSet := []int{}
	for k, v := range replyPair {
		var respJson map[string]string
		resp, err := w.Post(config.MongoDBApi+"/v1/modify", false, gin.H{
			"DataBaseName" : "Spider",
			"CollectionName" : "Post",
			"Filter" :{
				"_id" : k,
			},
			"ChangeField":{
				"ReplyID" 	: ReplySet["Data"].([]models.Reply)[v].Id,
				"Status" 	: "1",
			},
		})
		if err != nil {
			fmt.Fprintln(gin.DefaultWriter, err.Error())
			continue
		}
		if err := json.Unmarshal( resp.Body, &respJson); err != nil {
			fmt.Fprintln(gin.DefaultWriter, err.Error())
			continue
		}

		if respJson["StatusCode"] != "200" {
			errSet = append(errSet, k)
		}		
	}

	c.JSON( http.StatusOK, gin.H{
		"Msg":"ok",
		"StatusCode":"200",
		"ErrorSet" : errSet,
	})
}