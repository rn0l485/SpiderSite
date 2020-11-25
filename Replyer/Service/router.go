package Replyer

import (
	"github.com/gin-gonic/gin"

	"Decorations/Func/Worker"
	"Decorations/Replyer/Service/Models"
)
var (
	R 				*gin.Engine
	w 				*worker.Worker = worker.InitWorker()
)

func init() {
	R = gin.Default()

	R.GET("/", alive)
	
	v1 := R.Group("/v1")
	{
		v1.POST("/processing", processing)

	}



	R.NoRoute(pageNotFound)
	R.NoMethod(pageNotFound)
}

func alive(c *gin.Context) {
	c.JSON( http.StatusNotFound, gin.H{
		"Msg":"success",
		"StatusCode":"200",
	})	
}

func pageNotFound(c *gin.Context){
	c.JSON( http.StatusNotFound, gin.H{
		"Msg":"Path error",
		"StatusCode":"404",
	})
}
