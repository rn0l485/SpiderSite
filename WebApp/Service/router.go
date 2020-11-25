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

	store := cookie.NewStore([]byte(config.CookieSecret)) // config.CookieSecret

	R.Use(sessions.Sessions("status", store))
	

	R.LoadHTMLGlob("Decorations/WebApp/Service/templates/*.html")
	R.Static("/static", "Decorations/WebApp/Service/static")


	R.GET("/", LoginPage)
	R.GET("/setting", AuthRequired(),SettingPage)

	api := R.Group("/api") 
	{
		v001 := api.Group("/v001")
		{
			
		}
	}

	R.NoRoute(pageNotFound)
	R.NoMethod(pageNotFound)
}
