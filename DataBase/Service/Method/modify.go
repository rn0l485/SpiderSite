package methods

import (
	
	"net/http"
	"context"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"

	"Decorations/DataBase/Service/MongoDB/Models"
	"Decorations/DataBase/Service/MongoDB"
)

func Modify(c *gin.Context) {
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
	opts := options.FindOneAndUpdate().SetUpsert(true)

	filter := bson.M{}
	for k,v := range payload.Filter {
		filter[k] = v
	}

	var changingPart []bson.E = make([]bson.E, 0)
	for k,v := range payload.ChangeField{
		changingPart = append(changingPart, bson.E{ Key: k, Value:v })
	}
	update := bson.D{{"$set", changingPart}}

	record := targetCollection.FindOneAndUpdate(context.TODO(), filter, update, opts)

	err = record.Err()


	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"Msg": err.Error(),
					"StatusCode" : "500",
				},
			)
			return	
		}
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"Msg": err.Error(),
				"StatusCode" : "500",
			},
		)
		return	
	} else {
		c.JSON( http.StatusOK, gin.H{
			"Msg": "success",
			"StatusCode" : "200",
		})
	}




}