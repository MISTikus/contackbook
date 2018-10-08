package main

import (
	"log"
	"net/http"
	"os"

	"github.com/MISTikus/contactbook/common"
	"github.com/MISTikus/contactbook/contactbook"
	"github.com/julienschmidt/httprouter"
)

func main() {
	port := os.Args[1]

	apiservice := contactbook.NewApi()
	commonservice := common.NewApi()
	viewservice := contactbook.NewView()

	router := httprouter.New()
	for _, r := range apiservice.Routes {
		router.Handle(r.Method, "/api/"+r.Url, r.Handler)
	}
	for _, r := range viewservice.Routes {
		router.Handle(r.Method, "/view/"+r.Url, r.Handler)
	}
	for _, r := range commonservice.Routes {
		router.Handle(r.Method, "/api/"+r.Url, r.Handler)
	}

	router.Handle(viewservice.DefaultRoute.Method, "/", viewservice.DefaultRoute.Handler)

	log.Println("Service started listening on 'http://localhost:"+port+"' ...")

	log.Fatal(http.ListenAndServe(":"+port+"", router))

	log.Println("Service stopped ...")
}
