package web

import (
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

	resp, err := w.Post( config.MongoDBApi+"/v1/search", false, gin.H{
		"DataBaseName":"Spider",
		"CollectionName":"User",
		"Filter" : {
			"account": payload.
		}
	})


}
func Set(c *gin.Context) {}
func Data(c *gin.Context) {}