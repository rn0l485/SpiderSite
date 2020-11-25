package main

import (
	"log"
	"net/http"
	"time"


	db		"Decorations/DataBase/Service"
	"Decorations/DataBase/Config"
	"Decorations/DataBase/Service/MongoDB"

	"golang.org/x/sync/errgroup"


)

var (
	g errgroup.Group
)

func main() {
	mgo.CreateClients()
	defer mgo.DestroyClients()



	srv1 := &http.Server{
		Addr:		config.Port,
		Handler:	db.R,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		err := srv1.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
	
}
