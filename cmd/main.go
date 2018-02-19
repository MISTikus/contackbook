package main

import (
	"log"
	"net/http"

	"github.com/MISTikus/contactbook/contactbook"
	"github.com/julienschmidt/httprouter"
)

func main() {
	apiservice := contactbook.NewApi()
	viewservice := contactbook.NewView()

	router := httprouter.New()
	for _, r := range apiservice.Routes {
		router.Handle(r.Method, "/api/"+r.Route, r.Handler)
	}
	for _, r := range viewservice.Routes {
		router.Handle(r.Method, "/view/"+r.Route, r.Handler)
	}

	router.Handle(viewservice.DefaultRoute.Method, "/", viewservice.DefaultRoute.Handler)

	log.Println("Service started listening on 'http://localhost:9091' ...")

	log.Fatal(http.ListenAndServe(":9091", router))

	log.Println("Service stopped ...")
}
