package models


import (
	"time"
)

type Post struct {
	Id 					string 					`bson:"_id" 			json:"Id"`
	Domain	 			string 					`bson:"domain"			json:"Domain"`
	Group	 			string 					`bson:"group" 			json:"Group"`
	Url	 				string 					`bson:"url" 			json:"Url"`
	Client 				string 					`bson:"client"			json:"Client"`
	ClientUrl 			string 					`bson:"client_url"		json:"ClientUrl"`
	Context 			string 					`bson:"context" 		json:"Context"`

	// reply account
	ReplyAccount 		string 					`bson:"reply_account"	json:"ReplyAccount"`
	ReplyID				int 					`bson:"reply_id"		json:"ReplyID"`
	ComeBack 			bool 					`bson:"come_back"		json:"ComeBack"`

	CreateTime 			*time.Time 				`bson:"create_time"		json:"CreateTime"`
	Status 				string 					`bson:"status" 			json:"Status"`
}

type Reply struct {
	Id 					int 					`bson:"_id" 			json:"Id"`
	Keyword 			string 					`bson:"context" 		json:"Keyword"`		
	ReplyStatment 		string 					`bson:"reply_statment" 	json:"ReplyStatment"`
	Weights 			int 	 				`bson:"weights" 		json:"Weights"`
}

type Setting struct {
	Id 					int 					`bson:"_id" 			json:"Id"`
	Name 				string 					`bson:"name" 			json:"Name"`
	Value 				interface{} 			`bson:"value" 			json:"Value"`
	SetTime 			*time.Time 				`bson:"set_time" 		json:"SetTime"`
	PreviousValue 		interface{} 			`bson:"previous_value" 	json:"PreviousValue"`
	Comments 			string 					`bson:"comments" 		json:"Comments"`
}

type User struct {
	Id 					int 					`bson:"_id" 			json:"Id"`
	Account 			string 					`bson:"account" 		json:"Account"`
	Password			string 					`bson:"password" 		json:"Password"`
	Name 				string 					`bson:"name"			json:"Name"`
	Lv 					string 					`bson:"lv" 				json:"Lv"`
	CreateTime 			*time.Time 				`bson:"create_time"		json:"CreateTime"`
}

type Log struct {
	Id 					int 					`bson:"_id" 			json:"Id"`
	IP 					string 					`bson:"ip" 				json:"IP"`
	CreateTime 			*time.Time 				`bson:"create_time"		json:"CreateTime"`
	Context 			string 					`bson:"context" 		json:"Context"`
	Path 				string 					`bson:"path" 			json:"Path"`
	Service 			string 					`bson:"service" 		json:"Service"`
}

type Payload struct {
	Filter 					map[string]interface{} 						`json:"Filter"`
	Record 					interface{} 								`json:"Record"`
	Post	 				Post 										`json:"Post"`
	Reply 					Reply 										`json:"Reply"`
	Setting 				map[string]interface{} 						`json:"Setting"`
	User 					User 										`json:"User"`
	ChangeField 			map[string]interface{}						`json:"ChangeField"`
	DataBaseName			string 										`json:"DataBaseName"`
	CollectionName 			string 										`json:"CollectionName"`
}


