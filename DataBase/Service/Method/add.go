package methods

import (
	"net/http"
	"context"
	"github.com/gin-gonic/gin"

	

	"Decorations/DataBase/Service/MongoDB/Models"
	"Decorations/DataBase/Service/MongoDB"
)

func Add(c *gin.Context) {
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

	targetCollection := mgoC.Database( payload.DataBaseName).Collection( payload.CollectionName )

	_, err = targetCollection.InsertOne(context.TODO(), payload.Record )
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