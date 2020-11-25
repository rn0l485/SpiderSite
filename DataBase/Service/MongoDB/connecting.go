package mgo


import (
	"errors"
	"time"
	//"context"
	"go.mongodb.org/mongo-driver/mongo"

	"Decorations/DataBase/Config"
)

const mgoNum int = config.MongoDBConnectionNum
const mgoURL string = config.MongoDBURL //"mongodb://viewBot:weel99699@35.203.141.223:27017/UserList?connect=direct"
var clients chan *mongo.Client = make(chan *mongo.Client, mgoNum )


func CreateClients(){
	var c *mongo.Client
	var err error
	for i:=0 ; i<mgoNum ; i++ {
		c, err = connecting(mgoURL)
		if err != nil { 
			panic(err) 
		}
		if err := ping(c); err != nil {
			panic(err)
		} else {
			clients <- c
		}
	}
}

func GetClient() (*mongo.Client, error) {
	select {
		case c := <- clients :
			return c, nil
		case <- time.After(time.Second * 10) : 
			return nil, errors.New("Client waiting timeout.")
	}
}

func ReturnClient(c *mongo.Client) error {
	select {
		case clients <- c:
			return nil
		case <- time.After(time.Second * 10) :
			return errors.New("Client return error.")
	}
}

func DestroyClients() error {
	for {
		select {
		case c := <- clients :
			err := disconnecting(c)
			if err != nil {
				return err
			}
		case <- time.After(10 * time.Second):
			return nil 
		}
	}
}