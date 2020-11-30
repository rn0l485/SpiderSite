package web

import (

	"Decorations/Replyer/Func/Worker"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

var 	w 	*worker.Worker = worker.InitWorker()

func Page(c *gin.Context) {
	switch path:=c.FullPath(); path {
	case "/" :
	case "/setting":
	case "/result":     
	}
}

func Login(c *gin.Context){
	var payload models.Payload  // tbd
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{
				"Msg": err.Error(),
				"StatusCode" : "404",
			},
		)
		return
	}
	var respJson map[string]interface{}
	resp, err := w.Post( config.MongoDBApi+"/v1/search", false, gin.H{
		"DataBaseName":"Spider",
		"CollectionName":"User",
		"Filter" : {
			"account": payload.Account,
		}
	})
	if err != nil {
		fmt.Fprintln(gin.DefaultWriter, err.Error())
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{
				"Msg": "Error",
				"StatusCode" : "404",
			},
		)
		return
	}
	if err := json.Unmarshal( resp.Body, &respJson); err != nil {
		fmt.Fprintln(gin.DefaultWriter, err.Error())
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg": "Error",
			"StatusCode" : "404",
		})
		return
	}

	if  v, ok := respJson["Data"].([]map[string]interface{})[0] ; ok {
		if v["Password"].(string) == payload.Password {
			session := sessions.Default(c)
			session.Set("right", "1")
			sessionA.Save()

			c.JSON( http.StatusOK, gin.H{
				"Msg" : "ok",
				"StatusCode" : "200",
			})
			return
		}
	}

	fmt.Fprintln(gin.DefaultWriter, "login error")
	c.JSON( http.StatusNotFound, gin.H{
		"Msg" : "error"
		"StatusCode" : "404"
	})
}
func Set(c *gin.Context) {
	var payload models.Payload  // tbd
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"Msg": err.Error(),
			"StatusCode" : "404",
		})
		return
	}
	if payload.Service == "setting" {

	} else if payload.Service == "keyword" {

	} else {
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg" : "Error",
			"StatusCode" : "404",
		})
	}
}
func Data(c *gin.Context) {
	var payload models.Payload  // tbd
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"Msg": err.Error(),
			"StatusCode" : "404",
		})
		return
	}

		
}