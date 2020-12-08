package web

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

func SetHeader() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("access-control-allow-origin", "http://127.0.0.1:3000")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, origin, Cache-Control, User-Agent, Referer")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}