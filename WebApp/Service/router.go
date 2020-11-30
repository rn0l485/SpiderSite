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
	

	R.LoadHTMLGlob("Decorations/WebApp/Service/templates/*.html")
	R.Static("/static", "Decorations/WebApp/Service/static")


	R.GET("/", Page)
	R.GET("/setting", 	AuthRequired(), Page)
	R.GET("/result",	AuthRequired(), Page)

	api := R.Group("/api") 
	{
		v001 := api.Group("/v001")
		{
			v001.POST( "/login"		, Login )
			v001.POST( "/setting"	, AuthRequired()	, Set)
			v001.POST( "/data"		, AuthRequired()	, Data)
			v001.POST( "/alive"		, AuthRequired()	, AliveCheck)
		}
	}

	R.NoRoute(pageNotFound)
	R.NoMethod(pageNotFound)

	go Start()
}
