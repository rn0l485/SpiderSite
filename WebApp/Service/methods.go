package web

import (

	"Decorations/Replyer/Func/Worker"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

var  w 	*worker.Worker = worker.InitWorker()

func Alive(c *gin.Context) {
	c.JSON( http.StatusOK, gin.H{
		"Msg":"ok",
		"StatusCode":"200",
	})
}

func Login(c *gin.Context){
	var payload map[string]string  // tbd
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{
				"Msg": err.Error(),
				"StatusCode" : "404",
			},
		)
		return
	}
	var respJson map[string]interface{}
	resp, err := w.Post( config.MongoDBApi+"/v1/search", false, gin.H{
		"DataBaseName":"Spider",
		"CollectionName":"LoginUser",
		"Filter" : gin.H{
			"account": payload["Account"],
		}
	})
	if err != nil {
		fmt.Fprintln(gin.DefaultWriter, err.Error())
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{
				"Msg": "DataBase Error",
				"StatusCode" : "500",
			},
		)
		return
	}
	if err := json.Unmarshal( resp.Body, &respJson); err != nil {
		fmt.Fprintln(gin.DefaultWriter, err.Error())
		c.AbortWithStatusJSON( http.StatusNotFound, gin.H{
			"Msg": "Error",
			"StatusCode" : "404",
		})
		return
	}

	if  v, ok := respJson["Data"].([]map[string]interface{})[0] ; ok {
		if v["Password"].(string) == payload["Password"] {
			session := sessions.Default(c)
			session.Set("right", "1")
			sessionA.Save()

			c.JSON( http.StatusOK, gin.H{
				"Msg" : "ok",
				"StatusCode" : "200",
			})
			return
		}
	}

	fmt.Fprintln(gin.DefaultWriter, "login error")
	c.JSON( http.StatusNotFound, gin.H{
		"Msg" : "Login error"
		"StatusCode" : "404"
	})
}

func Set(c *gin.Context) {
	var payload []map[string]string 
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"Msg": err.Error(),
			"StatusCode" : "404",
		})
		return
	}

	ErrorMessageSet := make([][]string) 

	for i, v := range payload {
		if s, ok := v["Service"]; !ok{ 
			v["Service"] = "Error" 
		} 
		if d, ok := v["Do"]; !ok { 
			v["Do"] = "Error" 
		}



		switch service := v["Service"] ; service {
		case "setting":
			basicInfo := true
			if a,ok := v["Account"]; !ok {
				basicInfo = false
			}
			if p,ok := v["Password"]; !ok {
				if v["Do"] == "add"{
					basicInfo = false
				}
			}
			if d,ok := v["Domain"]; !ok {
				basicInfo = false
			}
			if !basicInfo {
				ErrorMessageSet = append(ErrorMessageSet, []string{
					"Basic Info Error",
					"Basic Info Error",
					"Basic Info Error",
				})
				continue
			}

			if v["Do"] == "add" {
				resp, err := w.Post( config.MongoDBApi + "/v1/add", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "User",
					"Record" : map[string]string{
						"account" 	: v["Account"],
						"password"	: v["Password"],
						"domain" 	: v["Domain"], 
						"service" 	: v["Service"],
					}
				})

				if err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["Account"],
						err.Error(),
					})
					continue
				}
				var respJson map[string]string
				if err := json.Unmarshal( resp.Body, &respJson); err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["Account"],
						err.Error(),
					})
					continue
				}

				if respJson["StatusCode"] != "200" {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["Account"],
						respJson["Msg"],
					})
					continue
				} 

			} else if v["Do"] == "delete" {
				resp, err := w.Post( config.MongoDBApi + "/v1/delete", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "User",
					"Filter" : map[string]string{
						"account" 	: v["Account"],
						"domain" 	: v["Domain"], 
						"service" 	: v["Service"],
					}
				})

				if err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["Account"],
						err.Error(),
					})
					continue
				}
				var respJson map[string]string
				if err := json.Unmarshal( resp.Body, &respJson); err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["Account"],
						err.Error(),
					})
					continue
				}

				if respJson["StatusCode"] != "200" {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["Account"],
						respJson["Msg"],
					})
					continue
				} 		
			} else {
				ErrorMessageSet = append(ErrorMessageSet, []string{
					v["Domain"],
					v["Account"],
					"Operation error",
				})
				continue
			}
		case "keyword":
			basicInfo := true
			if keyword, ok := v["Key"]; !ok {
				basicInfo = false
			}	
			if rs, ok := v["Value"]; !ok {
				basicInfo = false 
			}
			if w, ok := v["Weight"]; !ok {
				if v["Do"] == "add"{
					basicInfo = false
				}
			}
			if !basicInfo {
				ErrorMessageSet = append(ErrorMessageSet, []string{
					"Basic Info Error",
					"Basic Info Error",
					"Basic Info Error",
				})
				continue
			}


			if v["Do"] == "add" {
				resp, err := w.Post( config.MongoDBApi + "/v1/add", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "Reply",
					"Record" : map[string]string{
						"keyword" 			: v["Key"],
						"reply_statment"	: v["Value"],
						"weights" 			: v["Weight"],
						"service" 			: v["Service"],
					}
				})

				if err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Key"],
						v["Value"],
						err.Error(),
					})
					continue
				}
				var respJson map[string]string
				if err := json.Unmarshal( resp.Body, &respJson); err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Key"],
						v["Value"],
						err.Error(),
					})
					continue
				}

				if respJson["StatusCode"] != "200" {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Key"],
						v["Value"],
						respJson["Msg"],
					})
					continue
				}
			} else if v["Do"] == "delete" {
				resp, err := w.Post( config.MongoDBApi + "/v1/delete", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "Reply",
					"Filter" : map[string]string{
						"keyword" 			: v["Key"],
						"reply_statment"	: v["Value"],
						"service" 			: v["Service"],
					}
				})

				if err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Key"],
						v["Value"],
						err.Error(),
					})
					continue
				}
				var respJson map[string]string
				if err := json.Unmarshal( resp.Body, &respJson); err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Key"],
						v["Value"],
						err.Error(),
					})
					continue
				}

				if respJson["StatusCode"] != "200" {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Key"],
						v["Value"],
						respJson["Msg"],
					})
					continue
				}			
			} else {
				ErrorMessageSet = append(ErrorMessageSet, []string{
					v["Key"],
					v["Value"],
					"Operation error",
				})
				continue
			}
		case "select" :
			basicInfo := true 
			if d, ok := v["Domain"]; !ok {
				basicInfo = false 
			}
			if ca, ok := v["CurrentAccount"]; !ok {
				basicInfo = false
			}
			if ds, ok := v["DailyScrapingOn"]; !ok {
				basicInfo = false
			}
			if sfd, ok := v["SaveForXDays"]; !ok {
				basicInfo = false
			}
			if !basicInfo {
				ErrorMessageSet = append(ErrorMessageSet, []string{
					"Basic Info Error",
					"Basic Info Error",
					"Basic Info Error",
				})
				continue
			}			

			if v["Do"] == "add" {	
				resp, err := w.Post( config.MongoDBApi + "/v1/add", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "Config",
					"Record" : map[string]string{
						"domain" 				: v["Domain"],
						"service" 				: v["Service"],
						"current_account"		: v["CurrentAccount"],      	
						"daily_scraping_on" 	: v["DailyScrapingOn"],			
						"save_for_x_days" 		: v["SaveForXDays"],
					}
				})

				if err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						"no",
						err.Error(),
					})
					continue
				}
				var respJson map[string]string
				if err := json.Unmarshal( resp.Body, &respJson); err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						"no",
						err.Error(),
					})
					continue
				}

				if respJson["StatusCode"] != "200" {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						"no",
						respJson["Msg"],
					})
					continue
				}

			} else if v["Do"] == "delete" {
				resp, err := w.Post( config.MongoDBApi + "/v1/delete", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "Config",
					"Filter" : map[string]string{
						"domain" 				: v["Domain"],
						"service" 				: v["Service"],
						"current_account"		: v["CurrentAccount"],      	
						"daily_scraping_on" 	: v["DailyScrapingOn"],			
						"save_for_x_days" 		: v["SaveForXDays"],
					}
				})

				if err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						"no",
						err.Error(),
					})
					continue
				}
				var respJson map[string]string
				if err := json.Unmarshal( resp.Body, &respJson); err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						"no",
						err.Error(),
					})
					continue
				}

				if respJson["StatusCode"] != "200" {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						"no",
						respJson["Msg"],
					})
					continue
				}
			} else {
				ErrorMessageSet = append(ErrorMessageSet, []string{
					v["Domain"],
					"no",
					"Operation error",
				})
				continue
			}
		default:
			ErrorMessageSet = append(ErrorMessageSet, []string{
				"Option Error:" + service,
				"no",
				"Operation error",
			})
			continue
		}
	}

	if errNum := len(ErrorMessageSet); errNum > 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Msg" : "Error",
			"StatusCode" : "404",
			"Data" : ErrorMessageSet,
		})
	} else {
		c.JSON( http.StatusOK, gin.H{
			"Msg" :　"ok",
			"StatusCode" : "200",
		})
	}
}


func Data(c *gin.Context) {
	var payload map[string]string 
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"Msg": err.Error(),
			"StatusCode" : "404",
		})
		return
	}
	if serv, ok := payload["Service"]; !ok {
		c.JSON( http.StatusNotFound, gin.H{
			"Msg" 			: "Basic Info Error",
			"StatusCode" 	: "404",
		})
		return
	}


	switch service := payload["Service"]; service {
	case "data":

		fromDate, err := strconv.Atoi(payload["From"])
		if err != nil {
			c.JSON( http.StatusNotFound, gin.H{
				"Msg":err.Error(),
				"StatusCode":"404",
			})
		}
		toDate, err := strconv.Atoi(payload["To"])
		if err != nil {
			c.JSON( http.StatusNotFound, gin.H{
				"Msg":err.Error(),
				"StatusCode":"404",
			})
		}		

		resp, err := w.Post( config.MongoDBApi + "/v1/search", false, gin.H{
			"DataBaseName" : "Spider",
			"CollectionName" : "Post",
			"Filter" : gin.H{
				"create_time" : gin.H{
					"$gte" 	: fromDate,
					"$lte"	: toDate,
				},
			},
		})
		if err != nil {
			c.JSON( http.StatusNotFound, gin.H{
				"Msg": err.Error(),
				"StatusCode":"404",
			})
			return
		}
		var respJson map[string]interface{}
		if err := json.Unmarshal( resp.Body, &respJson); err != nil {
			c.JSON( http.StatusNotFound, gin.H{
				"Msg" : err.Error(),
				"StatusCode":"404",
			})
			return 	
		}		
		if respJson["StatusCode"].(string) != "200" {
			c.JSON( http.StatusNotFound, gin.H{
				"Msg" : "Internal Error",
				"StatusCode":"500",
			})
			return
		}

		passBack := make([]map[string]string)
		for _,v := respJson["Data"].([]models.Post) {
			passBack = append( passBack, map[string]string{
				"Domain"		: v.Domain,
				"Url"			: v.Url,
				"CreateTime"	: strconv.Itoa(v.CreateTime),
				"ReplyAccount"	: v.ReplyAccount,
				"ReplyKeyWord"	: v.ReplyKeyword,
			})
		}

		datalen, err := strconv.Atoi(len(passBack))
		if err != nil {
			c.JSON( http.StatusNotFound, gin.H{
				"Msg":err.Error(),
				"StatusCode":"404",
			})
		}

		c.JSON( http.StatusOK, gin.H{
			"Msg" : datalen,
			"StatusCode" : "200",
			"Data" : passBack,
		})
		return 
	case "setting":
		resp, err := w.Post(config.MongoDBApi + "/v1/search", false, gin.H{
			"DataBaseName" : "Spider",
			"CollectionName" : "User",
			"Filter" : {},
		})
		if err != nil {
			c.JSON( http.StatusNotFound, gin.H{
				"Msg":err.Error(),
				"StatusCode":"404",
			})			
		}
		var respJson map[string]interface{}
		if err := json.Unmarshal( resp.Body, &respJson); err != nil {
			c.JSON( http.StatusNotFound, gin.H{
				"Msg" : err.Error(),
				"StatusCode":"404",
			})
			return 	
		}

		if respJson["StatusCode"].(string) != "200" {
			c.JSON( http.StatusNotFound, gin.H{
				"Msg" : "Internal Error",
				"StatusCode":"404",
			})
			return 	
		}

		AccountSaved := []string{}
		for _, user := respJson["Data"].([]models.User) {
			AccountSaved = append(AccountSaved, user.Account)
		}






		




	case "keyword":
	default:
	}
		
}

func AliveCheck(c *gin.Context){
	c.JSON( http.StatusOK, gin.H{
		"Msg" : "ok",
		"StatusCode":"200",
	})
}