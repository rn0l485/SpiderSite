package methods

import (
	"time"
	"net/http"
	"context"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"

	"Decorations/DataBase/Service/MongoDB/Models"
	"Decorations/DataBase/Service/MongoDB"
)

func Delete(c *gin.Context) {
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

	filter := bson.M{}
	for k,v := range payload.Filter {
		filter[k] = v
	}

	targetCollection := mgoC.Database(payload.DataBaseName).Collection( payload.CollectionName )

	changingPart := []bson.E{} 
	changingPart = append(changingPart, bson.E{ Key: "status", Value: time.Now().Format("2006-01-02 15:04:05") })
	update := bson.D{{"$set", changingPart}}
	opts := options.FindOneAndUpdate().SetUpsert(true)

	record := targetCollection.FindOneAndUpdate(context.TODO(), filter, update, opts)

	err = record.Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"Msg": err.Error(),
					"StatusCode" : "500",
				},
			)		
		} else {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"Msg": err.Error(),
					"StatusCode" : "500",
				},
			)	
		}
		return
	} else {
		c.JSON( 
			http.StatusOK, 
			gin.H{
				"Msg": "success",
				"StatusCode" : "200",
			},
		)
	}

}

func Clear(c *gin.Context) {
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

	filter := bson.M{"status": bson.M{"$ne": "alive"}}

	_, err = targetCollection.DeleteMany(context.TODO(), filter)

	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"Msg": err.Error(),
				"StatusCode" : "500",
			},
		)	
		return
	} else {
		c.JSON( 
			http.StatusOK, 
			gin.H{
				"Msg": "success",
				"StatusCode" : "200",
			},
		)
	}

}