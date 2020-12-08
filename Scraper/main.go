package main

import (
	"log"
	"net/http"
	"time"

	sp 		"Decorations/Models/Scraper"
	"Decorations/Scraper/Config"

	"golang.org/x/sync/errgroup"



)

var (
	g errgroup.Group
)

func main() {
	srv1 := &http.Server{
		Addr:		config.Port,
		Handler:	sp.R,
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
