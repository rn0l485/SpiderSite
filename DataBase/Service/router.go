package DBServer

import (
	"os"
	"io"

	"net/http"
	"github.com/gin-gonic/gin"

	"Decorations/DataBase/Service/Method"
	
)

var R *gin.Engine 

func init(){
	gin.SetMode(gin.ReleaseMode)

	var f *os.File
	if _, err := os.Stat("./DB.log"); err == nil {
		f,_ = os.OpenFile("./DB.log", os.O_RDWR|os.O_CREATE, 0755)
	} else if os.IsNotExist(err) {
		f,_ = os.Create("./DB.log")
	} else {
		f,_ = os.OpenFile("./DB.log", os.O_RDWR|os.O_CREATE, 0755)
	}

	gin.DefaultWriter = io.MultiWriter(f)

	R = gin.Default()

	R.GET("/", alive)
	v1 := R.Group("/v1")
	{
		v1.POST("/add"		, methods.Add)
		v1.POST("/search"	, methods.Search)
		v1.POST("/modify"	, methods.Modify)
		v1.POST("/delete"	, methods.Delete)
		v1.POST("/clear"	, methods.Clear)
	}


	R.NoRoute(pageNotFound)
	R.NoMethod(pageNotFound)
}

func alive( c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Msg":"Service alive",
		"StatusCode":"200",
	})
}

func pageNotFound(c *gin.Context){
	c.JSON( http.StatusNotFound, gin.H{
		"Msg":"Path error",
		"StatusCode":"404",
	})
}