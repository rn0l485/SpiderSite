package main

import (
	"log"
	"net/http"
	"time"


	db		"Decorations/Models/DataBase"
	sp 		"Decorations/Models/Scraper"
	web 	"Decorations/Models/WebApp"


	"golang.org/x/sync/errgroup"


)

var (
	g errgroup.Group
)

func main() {
	defer sp.Cancel()




	srv1 := &http.Server{
		Addr:		":9000",
		Handler:	web.R,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	srv2 := &http.Server{
		Addr:		":9001",
		Handler:	sp.R,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}	

	srv3 := &http.Server{
		Addr:		":9002",
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

	g.Go(func() error {
		err := srv2.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	g.Go(func() error {
		err := srv3.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})
	

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
	
}
