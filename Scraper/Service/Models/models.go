package ScraperModels


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
	Reply 				Reply 					`bson:"reply"			json:"Reply"`
	ComeBack 			bool 					`bson:"come_back"		json:"ComeBack"`

	CreateTime 			*time.Time 				`bson:"create_time"		json:"CreateTime"`
	Status 				string 					`bson:"status" 			json:"Status"`
}


type Reply struct {
	ReplyAccount 		string 					`bson:"reply_account"	json:"ReplyAccount"`
	Context 			string 					`bson:"context" 		json:"Context"`		
	ReplyTime 			*time.Time 				`bson:"reply_time"		json:"ReplyTime"`
}

type Setting struct {
	Id 					int 					`bson:"_id" 			json:"Id"`
	Name 				string 					`bson:"name" 			json:"Name"`
	Value 				interface{} 			`bson:"value" 			json:"Value"`
	SetTime 			*time.Time 				`bson:"set_time" 		json:"SetTime"`
	PreviousValue 		interface{} 			`bson:"previous_value" 	json:"PreviousValue"`
	Comments 			string 					`bson:"comments" 		json:"Comments"`
}

type Resp struct {
	Msg 				string 					`json:"Msg"`
	StatusCode 			string 					`json:"StatusCode"`
	Data				[]interface{}			`json:"Data"`
}

type Payload struct {
	Account 			string 					`json:"Account"`
	Method 				string 					`json:"Method"`
	Data 				[]string 				`json:"Data"`
}