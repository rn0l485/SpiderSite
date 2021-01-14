package ScraperModels


import (
	"time"
)
type Resp struct {
	Msg 				string 					`json:"Msg"`
	StatusCode 			string 					`json:"StatusCode"`
	Data				[]interface{}			`json:"Data"`
}


type Payload struct {
	Account 			*string 				`json:"Account"`
	Password 			*string 				`json:"Password"`
	Msg 				*string 				`json:"Msg"`
	Url 				*string 				`json:"Url"`
	Domain 				*string 				`json:"Domain"`
}

type Page struct {
	CTX					context.Context
	Cancel 				context.CancelFunc
}