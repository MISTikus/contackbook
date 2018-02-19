package contactbook

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/MISTikus/contactbook/apimodels"
	"github.com/MISTikus/contactbook/common"
	"github.com/julienschmidt/httprouter"
)

type api struct {
	Routes []common.Route
}

func NewApi() api {
	service := api{}
	service.Routes = []common.Route{
		{
			Route:   "contact/:id",
			Method:  common.Get,
			Handler: service.getById,
		},
		{
			Route:   "contact",
			Method:  common.Post,
			Handler: service.add,
		},
	}
	return service
}

func (service api) getById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idString := strings.TrimSpace(p.ByName("id"))
	_, err := strconv.ParseInt(idString, 10, 32)
	if idString == "" || err != nil {
		badRequest(w, "Identifier '"+idString+"' can not be parsed")
		return
	}

	log.Println("Resolving contact with id: " + idString)

	contact := apimodels.Contact{Name: "SomeName", Phone: "13251221", Description: "Нигадяй..."}

	if err = response(w, contact); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func badRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func (service api) add(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println(r)
	log.Println(p)
	w.Write([]byte("new task"))
}

func response(w http.ResponseWriter, obj interface{}) error {
	js, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	return err
}
