package web

import (

	"Decorations/WebApp/Func/Worker"
	"Decorations/WebApp/Config"
	"Decorations/WebApp/Service/Models"
	"encoding/json"

	"fmt"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)

var  w 	*worker.Worker = worker.InitWorker()

func Login(c *gin.Context){
	var payload map[string]string  
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
	var respJson models.CallBack
	resp, err := w.Post( config.MongoDBApi+"/v1/search", false, gin.H{
		"DataBaseName":"Spider",
		"CollectionName":"LoginUser",
		"Filter" : gin.H{
			"Account": payload["Account"],
		},
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

	
	if len(respJson.Data) == 0 {
		fmt.Fprintln(gin.DefaultWriter, "login error")
		c.JSON( http.StatusNotFound, gin.H{
			"Msg" : "Login error #1",
			"StatusCode" : "404",
		})
		return
	} else if respJson.Data[0]["Password"].(string) == payload["Password"] {
		session := sessions.Default(c)
		session.Set("right", "1")
		session.Save()
		c.JSON( http.StatusOK, gin.H{
			"Msg" : "ok",
			"StatusCode" : "200",
		})
		return
	} else {
		fmt.Fprintln(gin.DefaultWriter, "login error")
		c.JSON( http.StatusNotFound, gin.H{
			"Msg" : "Login error #2",
			"StatusCode" : "404",
		})	
		return		
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
	if _, ok := payload["Service"]; !ok {
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
				"status" : gin.H{
					"$eq": "alive",
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
				"Msg" : respJson["Msg"].(string),
				"StatusCode":"500",
			})
			return
		}

		if respJson["Data"] == nil {
			c.JSON( http.StatusOK, gin.H{
				"Msg" : "ok",
				"StatusCode" : "200",	
				"Data" : []gin.H{},			
			})
			return			
		}

		passBack := []map[string]string{}
		for _,v := range respJson["Data"].([]map[string]string) {
			newOne := map[string]string{
				"Domain"		: v["Domain"],
				"Url"			: v["Url"],
				"CreateTime"	: v["CreateTime"],
				"ReplyAccount"	: v["ReplyAccount"],
				"ReplyKeyWord"	: v["ReplyKeyword"],
			}

			passBack = append( passBack, newOne)
		}

		datalen := strconv.Itoa(len(passBack))

		c.JSON( http.StatusOK, gin.H{
			"Msg" : datalen,
			"StatusCode" : "200",
			"Data" : passBack,
		})
		return 
	case "setting":
		resp, err := w.Post(config.MongoDBApi + "/v1/search", false, gin.H{
			"DataBaseName" : "Spider",
			"CollectionName" : "Config",
			"Filter" : gin.H{
				"status" : gin.H{
					"$eq": "alive",
				},
			},
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
				"Msg" : respJson["Msg"].(string),
				"StatusCode":"404",
			})
			return 	
		}


		c.JSON( http.StatusOK, gin.H{
			"Msg" : "ok",
			"StatusCode" : "200",
			"Data" : respJson["Data"],
		})
		return
	case "keyword":
		resp, err := w.Post( config.MongoDBApi + "/v1/search", false, gin.H{
			"DataBaseName" : "Spider",
			"CollectionName" : "Reply",
			"Filter" : gin.H{
				"status" : gin.H{
					"$eq": "alive",
				},
			},
		})

		if err != nil {
			c.JSON( http.StatusNotFound, gin.H{
				"Msg" : err.Error(),
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
				"Msg" : respJson["Msg"].(string),
				"StatusCode":"404",
			})
			return 	
		}

		if respJson["Data"] == nil {
			c.JSON( http.StatusOK, gin.H{
				"Msg" : "ok",
				"StatusCode" : "200",	
				"Data" : []gin.H{},			
			})
			return
		}



		kvSet := make([]map[string][]string,0)
		for _,v := range respJson["Data"].([]map[string]string){
			kvSet = append(kvSet, map[string][]string{
				v["Keyword"] : []string{
					v["Value"],
					v["Weight"],
				},
			})
		}
		passBack := gin.H{
			"Msg" : "ok",
			"StatusCode" : "200",
			"Data" : kvSet,
		}

		c.JSON( http.StatusOK, passBack)
		return
	case "info":
		resp, err := w.Post( config.MongoDBApi+"/v1/search", false, gin.H{
			"DataBaseName": "Spider",
			"CollectionName": "UserGroup",
			"Filter" : gin.H{
				"status" : gin.H{
					"$eq": "alive",
				},
				"domain" : gin.H{
					"$eq": payload["Domain"], 
				},
			},
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
				"Msg" : respJson["Msg"].(string),
				"StatusCode":"404",
			})
			return 	
		}

		if respJson["Data"] == nil {
			c.JSON( http.StatusOK, gin.H{
				"Msg" : "ok",
				"StatusCode" : "200",	
				"Data" : []gin.H{},			
			})
			return			
		}

		accountUnderDomainMap := map[string]bool{}
		urlUnderAccount := []string{}
		for _, u := range respJson["Data"].([]map[string]string) {
			accountUnderDomainMap[u["Account"]] = true
			if u["Account"] == payload["Account"] {
				urlUnderAccount = append( urlUnderAccount, u["Url"])
			}
		}

		accountUnderDomain := []string{}
		for k,_ := range accountUnderDomainMap {
			accountUnderDomain = append(accountUnderDomain, k)
		}

		c.JSON( http.StatusOK, gin.H{
			"Msg" : "ok",
			"StatusCode" : "200",
			"Data": []gin.H{
				gin.H{
					"Domain" 		: payload["Domain"], 
					"AccountSaved"	: accountUnderDomain,
					"CurrentAccountInfo" : gin.H{
						"Group" : urlUnderAccount,
					},
				},
			},
		})
		return 	
	default:
		c.JSON(http.StatusNotFound, gin.H{
			"Msg" : "Error",
			"StatusCode" : "404",
		})
	}	
}

func Set(c *gin.Context) {
	var payloadData map[string][]map[string]string 
	if err := c.ShouldBindJSON(&payloadData); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"Msg": err.Error(),
			"StatusCode" : "404",
		})
		return
	}

	payload := payloadData["Data"]

	ErrorMessageSet := [][]string{}

	for _, v := range payload {
		if _, ok := v["Service"]; !ok  { 
			v["Service"] = "Error" 
		} 
		if _, ok := v["Do"]; !ok { 
			v["Do"] = "Error" 
		}

		switch service := v["Service"] ; service {
		case "setting":
			basicInfo := true
			if _,ok := v["Account"]; !ok {
				basicInfo = false
			}
			if _,ok := v["Password"]; !ok {
				if v["Do"] == "add"{
					basicInfo = false
				}
			}
			if _,ok := v["Domain"]; !ok {
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
				resp, err := w.Post( config.MongoDBApi + "/v1/search", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "User",
					"Filter" : gin.H{
						"status" 	: gin.H{
							"$eq"	: "alive",
						},
						"account" 	: gin.H{
							"$eq"	: v["Account"],
						},
						"password"	: gin.H{
							"$eq"	: v["Password"],
						},
						"domain" 	: gin.H{
							"$eq"	: v["Domain"], 
						},
					},				
				})
				if err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["Account"],
						err.Error(),
					})
					continue					
				}
				var respJsonSearch map[string]interface{}
				if err := json.Unmarshal( resp.Body, &respJsonSearch); err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["Account"],
						err.Error(),
					})
					continue
				}

				if respJsonSearch["StatusCode"].(string) != "200" {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["Account"],
						err.Error(),
					})
					continue					
				}

				if respJsonSearch["Data"] != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["Account"],
						"Data exsist",
					})
					continue					
				}

				resp, err = w.Post( config.MongoDBApi + "/v1/add", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "User",
					"Record" : map[string]string{
						"account" 	: v["Account"],
						"password"	: v["Password"],
						"domain" 	: v["Domain"], 
						"status"	: "alive",
					},
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
					"Filter" : gin.H{
						"account" 	: gin.H{
							"$eq"	: v["Account"],
						},
						"domain" 	: gin.H{
							"$eq"	: v["Domain"], 
						},
					},
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
			if _, ok := v["Key"]; !ok {
				basicInfo = false
			}	
			if _, ok := v["Value"]; !ok {
				basicInfo = false 
			}
			if _, ok := v["Weight"]; !ok {
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
				resp, err := w.Post( config.MongoDBApi + "/v1/search", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "Reply",
					"Filter" : gin.H{
						"keyword" 			: gin.H{
							"$eq"			: v["Key"],
						},
						"reply_statment"	: gin.H{
							"$eq"			: v["Value"],
						},
					},				
				})

				if err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Key"],
						v["Value"],
						err.Error(),
					})
					continue					
				}
				var respJsonSearch map[string]interface{}
				if err := json.Unmarshal( resp.Body, &respJsonSearch); err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Key"],
						v["Value"],
						err.Error(),
					})
					continue
				}

				if respJsonSearch["StatusCode"].(string) != "200" {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Key"],
						v["Value"],
						err.Error(),
					})
					continue					
				}

				if respJsonSearch["Data"] != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Key"],
						v["Value"],
						"Data exsist",
					})
					continue					
				}	

				resp, err = w.Post( config.MongoDBApi + "/v1/add", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "Reply",
					"Record" : map[string]string{
						"keyword" 			: v["Key"],
						"reply_statment"	: v["Value"],
						"weights" 			: v["Weight"],
						"status"			: "alive",
					},
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
					},
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
			if _, ok := v["Domain"]; !ok {
				basicInfo = false 
			}
			if _, ok := v["DailyScrapingOn"]; !ok {
				basicInfo = false
			}
			if _, ok := v["SaveForXDays"]; !ok {
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
				resp, err := w.Post( config.MongoDBApi+"/v1/delete", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "Config",
					"Filter" : gin.H{
						"domain" 			: gin.H{
							"$eq" 			: v["Domain"],
						},
						"status" 			: gin.H{
							"$eq" 			: "alive",
						},
					},				
				})

				if err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						"no",
						err.Error(),
					})
					continue					
				}

				var respJsonDelete map[string]string
				if err := json.Unmarshal( resp.Body, &respJsonDelete); err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						"no",
						err.Error(),
					})
					continue
				}

				if respJsonDelete["StatusCode"] != "404" {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						"no",
						"Data exsist",
					})
					continue					
				}

				resp, err = w.Post( config.MongoDBApi + "/v1/add", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "Config",
					"Record" : map[string]string{
						"domain" 				: v["Domain"],
						"daily_scraping_on" 	: v["DailyScrapingOn"],			
						"save_for_x_days" 		: v["SaveForXDays"],
						"status" 				: "alive",
					},
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
		case "info" :
			basicInfo := true
			if _,ok := v["Domain"]; !ok {
				basicInfo = false
			}
			if _,ok := v["CurrentAccount"]; !ok {
				basicInfo = false
			}
			if _,ok := v["AddUrl"]; !ok {
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
				resp, err := w.Post( config.MongoDBApi + "/v1/search", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "UserGroup",
					"Filter" : gin.H{
						"domain" 	: gin.H{
							"$eq"	: v["Domain"],
						},
						"account"	: gin.H{
							"$eq"	: v["CurrentAccount"],    
						},
						"url" 		: gin.H{
							"$eq" 	: v["AddUrl"],  	
						},
						"status" 	: gin.H{
							"$eq"	: "alive",					
						},
					},				
				})

				if err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["CurrentAccount"],
						err.Error(),
					})
					continue					
				}
				var respJsonSearch map[string]interface{}
				if err := json.Unmarshal( resp.Body, &respJsonSearch); err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["CurrentAccount"],
						err.Error(),
					})
					continue
				}

				if respJsonSearch["StatusCode"].(string) != "200" {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["CurrentAccount"],
						err.Error(),
					})
					continue					
				}

				if respJsonSearch["Data"] != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["CurrentAccount"],
						"Data exsist",
					})
					continue					
				}	



				resp, err = w.Post( config.MongoDBApi + "/v1/add", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "UserGroup",
					"Record" : map[string]string{
						"domain" 				: v["Domain"],
						"account"				: v["CurrentAccount"],    
						"url" 					: v["AddUrl"],  	
						"status" 				: "alive",
					},
				})
				if err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["CurrentAccount"],
						err.Error(),
					})
					continue					
				}
				var respJson map[string]string
				if err := json.Unmarshal( resp.Body, &respJson); err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["CurrentAccount"],
						err.Error(),
					})
					continue
				}		
				if respJson["StatusCode"] != "200" {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["CurrentAccount"],
						respJson["Msg"],
					})
					continue
				} 
			} else if v["Do"] == "delete" {
				resp, err := w.Post( config.MongoDBApi + "/v1/delete", false, gin.H{
					"DataBaseName" : "Spider",
					"CollectionName" : "UserGroup",
					"Filter" : gin.H{
						"domain" 	: gin.H{
							"$eq" 	: v["Domain"],
						},
						"account"	: gin.H{
							"$eq" 	: v["CurrentAccount"],    
						},
						"url" 		: gin.H{
							"$eq" 	: v["AddUrl"], 
						},
					},
				})		
				if err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["CurrentAccount"],
						err.Error(),
					})
					continue					
				}		
				var respJson map[string]string
				if err := json.Unmarshal( resp.Body, &respJson); err != nil {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["CurrentAccount"],
						err.Error(),
					})
					continue
				}	
				if respJson["StatusCode"] != "200" {
					ErrorMessageSet = append(ErrorMessageSet, []string{
						v["Domain"],
						v["CurrentAccount"],
						respJson["Msg"],
					})
					continue
				} 
			} else {
				ErrorMessageSet = append(ErrorMessageSet, []string{
					v["Domain"],
					v["CurrentAccount"],
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
			"Msg":"ok",
			"StatusCode":"200",
		})
	}
}
