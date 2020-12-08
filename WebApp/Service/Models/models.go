package models


type CallBack struct {
	Msg 				string  					`json:"Msg"`
	StatusCode 			string 						`json:"StatusCode"`
	Data 				[]map[string]interface{}	`json:"Data"`
}