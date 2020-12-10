package web

import (


	"os"
	"io"
	"net/http"
	"time"
	"encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"


	"golang.org/x/crypto/bcrypt"
)


var R *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)

	var f *os.File
	if _, err := os.Stat("./Web.log"); err == nil {
		f,_ = os.OpenFile("./Web.log", os.O_RDWR|os.O_CREATE, 0755)
	} else if os.IsNotExist(err) {
		f,_ = os.Create("./Web.log")
	} else {
		f,_ = os.OpenFile("./Web.log", os.O_RDWR|os.O_CREATE, 0755)
	}

	gin.DefaultWriter = io.MultiWriter(f)

	R = gin.Default()



	R.Use(SetHeader())

	R.GET("/", Alive)

	api := R.Group("/api") 
	{
		v001 := api.Group("/v001")
		{
			v001.GET(  "/token"				, Token )
			v001.POST( "/login"				, Login )
			v001.POST( "/setting"			, AuthRequired()	, Set)
			v001.POST( "/data"				, AuthRequired()	, Data)
			v001.POST( "/alive"				, AuthRequired()	, AliveCheck)
			//v001.GET(  "/ticker/:domain"	, AuthRequired()	, Ticker)
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

func pageNotFound(c *gin.Context) {
	c.JSON( http.StatusOK, gin.H{
		"Msg":"Error",
		"StatusCode":"404",
	})
}

func Token(c *gin.Context) {
	newToke := GenerateToken(time.Now().String())

	c.JSON( http.StatusOK, gin.H{
		"Msg" : newToke,
		"StatusCode":"200",
	})
}

func GenerateToken(inputKey string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(inputKey), bcrypt.DefaultCost) //CompareHashAndPassword(hashedPassword, passwordNotCheck) 
	if err != nil {
		return err.Error() + " #1"
	}

	return base64.StdEncoding.EncodeToString(hash)
}

func AliveCheck(c *gin.Context){
	c.JSON( http.StatusOK, gin.H{
		"Msg" : "ok",
		"StatusCode":"200",
	})
} 

type Claims struct {
	Account 	string 			`json:"Account"`
	Role 		string 			`json:"Role"`
	jwt.StandardClaims
}