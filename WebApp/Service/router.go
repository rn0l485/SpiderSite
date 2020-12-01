package web

import (
	"Decorations/Config"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)


var R *gin.Engine

func init() {
	R = gin.Default()

	store := cookie.NewStore([]byte(config.CookieSecret))

	R.Use(sessions.Sessions("status", store))

	R.GET("/", Alive)

	api := R.Group("/api") 
	{
		v001 := api.Group("/v001")
		{
			v001.POST( "/login"		, PathTest) //Login )
			v001.POST( "/setting"	, AuthRequired()	, PathTest) //Set)
			v001.POST( "/data"		, AuthRequired()	, PathTest) //Data)
			v001.POST( "/alive"		, AuthRequired()	, PathTest) //AliveCheck)
		}
	}

	R.NoRoute(pageNotFound)
	R.NoMethod(pageNotFound)
}

func Alive(c *gin.Context) {
	c.JSON( http.StatusOK, gin.H{
		"Msg":"ok",
		"StatusCode":"200",
	})
}