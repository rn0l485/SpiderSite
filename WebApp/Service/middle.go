package web

import (
	"strings"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"

	"Decorations/WebApp/Config"

)
var jwtSecret = []byte(config.CookieSecret)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		token := strings.Split(auth, "Bearer ")
		payToken := token[len(token)-1]



		tokenClaims, err := jwt.ParseWithClaims(payToken, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
			return jwtSecret, nil
		})

		if err != nil {
			var message string
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors & jwt.ValidationErrorMalformed != 0 {
					message = "token is malformed"
				} else if ve.Errors & jwt.ValidationErrorUnverifiable != 0{
					message = "token could not be verified because of signing problems"
				} else if ve.Errors & jwt.ValidationErrorSignatureInvalid != 0 {
					message = "signature validation failed"
				} else if ve.Errors & jwt.ValidationErrorExpired != 0 {
					message = "token is expired"
				} else if ve.Errors & jwt.ValidationErrorNotValidYet != 0 {
					message = "token is not yet valid before sometime"
				} else {
					message = "can not handle this token"
				}
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"Msg": message,
				"StatusCode" : "401",
			})
			c.Abort()
			return
		}
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			c.Set("account", claims.Account)
			c.Set("role", claims.Role)
			c.Next()
		} else {
			c.Abort()
			return
		}		
	}
}

func SetHeader() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("access-control-allow-origin", "http://54.150.155.212")
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