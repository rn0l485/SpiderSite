package models


type Payload struct {
	Filter 					map[string]interface{} 						`json:"Filter"`
	Record 					interface{} 								`json:"Record"`
	ChangeField 			map[string]interface{}						`json:"ChangeField"`
	DataBaseName			string 										`json:"DataBaseName"`
	CollectionName 			string 										`json:"CollectionName"`
}