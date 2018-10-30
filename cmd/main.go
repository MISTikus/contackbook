package main

import (
	"github.com/MISTikus/contactbook/common"
	"github.com/MISTikus/contactbook/contactbook"
	"github.com/MISTikus/contactbook/data"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Args[1]

	file := data.FileWorker{FileName: "data.json"}
	repo := data.NewRepo(&file)

	viewService := contactbook.NewView(repo)
	apiService := contactbook.NewApi(repo)
	commonService := common.NewApi()

	viewService.Api = apiService
	viewService.Common = commonService
	apiService.ChangeHandler = viewService

	router := httprouter.New()

	var routes []string
	for _, r := range apiService.Routes {
		route := common.BuildUrl("api", apiService.Prefix, r.Url)
		routes = append(routes, r.Method+": "+route)
		router.Handle(r.Method, route, r.Handler)
	}
	for _, r := range viewService.Routes {
		route := common.BuildUrl("view", viewService.Prefix, r.Url)
		routes = append(routes, r.Method+": "+route)
		router.Handle(r.Method, route, r.Handler)
	}
	for _, r := range commonService.Routes {
		route := common.BuildUrl("api", commonService.Prefix, r.Url)
		routes = append(routes, r.Method+": "+route)
		router.Handle(r.Method, route, r.Handler)
	}

	router.Handle(viewService.DefaultRoute.Method, "/list", viewService.DefaultRoute.Handler)
	router.Handle(viewService.DefaultRoute.Method, "/", viewService.DefaultRoute.Handler)

	log.Println("Service started listening on 'http://localhost:" + port + "' ...")
	log.Println("Currently listening routes:")
	for _, route := range routes {
		log.Println("\t" + route)
	}

	log.Fatal(http.ListenAndServe(":"+port+"", router))

	log.Println("Service stopped ...")
}
