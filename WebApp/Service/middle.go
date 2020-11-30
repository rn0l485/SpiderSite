package router

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	//"etronErp/Config"
)


func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if status := session.Get("right"); status == nil {
			//c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
			c.AbortWithStatusJSON( http.StatusInternalServerError, gin.H{
				"Msg":"right error",
				"StatusCode" : "500",
			})
		} 
		c.Next()
	}
}