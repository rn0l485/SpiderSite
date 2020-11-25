package Setting

import (
	"github.com/gin-gonic/gin"

	"Decorations/Func/Worker"
	"Decorations/Scraper/Service/Models"
)

func init() {
	R = gin.Default()

	R.GET("/", alive)
	
	info := R.Group("/")

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