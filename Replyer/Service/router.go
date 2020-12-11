package Replyer

import (
	"os"
	"io"
	"net/http"
	"github.com/gin-gonic/gin"
	"Decorations/Replyer/Func/Worker"
)
var (
	R 				*gin.Engine
	w 				*worker.Worker = worker.InitWorker()
)

func init() {
	gin.SetMode(gin.ReleaseMode)

	var f *os.File
	if _, err := os.Stat("./Replyer.log"); err == nil {
		f,_ = os.OpenFile("./Replyer.log", os.O_RDWR|os.O_CREATE, 0755)
	} else if os.IsNotExist(err) {
		f,_ = os.Create("./Replyer.log")
	} else {
		f,_ = os.OpenFile("./Replyer.log", os.O_RDWR|os.O_CREATE, 0755)
	}

	gin.DefaultWriter = io.MultiWriter(f)

	R = gin.Default()

	R.GET("/", alive)
	
	v1 := R.Group("/v1")
	{
		v1.POST("/do", 	processing)
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
