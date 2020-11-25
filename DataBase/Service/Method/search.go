package methods

import (
	
	"net/http"
	"context"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	
	"Decorations/DataBase/Service/MongoDB/Models"
	"Decorations/DataBase/Service/MongoDB"
)

func Search(c *gin.Context) {
	var payload models.Payload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"Msg": err.Error(),
				"StatusCode" : "500",
			},
		)
		return
	}

	mgoC, err := mgo.GetClient()
	defer mgo.ReturnClient(mgoC)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"Msg": err.Error(),
				"StatusCode" : "500",
			},
		)
		return		
	}	

	targetCollection := mgoC.Database(payload.DataBaseName).Collection( payload.CollectionName )

	filter := bson.M{}
	for k,v := range payload.Filter {
		filter[k] = v
	}

	var targetPostSet []*map[string]interface{}
	cur, err := targetCollection.Find(
		context.TODO(), 
		filter,
	)
	if err != nil { 
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"Msg": err.Error(),
				"StatusCode" : "500",
			},
		)
		return 
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
	// To decode into a struct, use cursor.Decode()
		var result map[string]interface{}
		err := cur.Decode(&result)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"Msg": err.Error(),
					"StatusCode" : "500",
				},
			)
			return 			
		}
		targetPostSet = append(targetPostSet, &result)
	}
	if err := cur.Err(); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"Msg": err.Error(),
				"StatusCode" : "500",
			},
		)
		return 
	}

	c.JSON( http.StatusOK, gin.H{
		"Msg": "success",
		"StatusCode" : "200",
		"Data" : targetPostSet,
	})
}