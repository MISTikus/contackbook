package contactbook

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MISTikus/contactbook/apimodels"
	"github.com/MISTikus/contactbook/common"
	"github.com/julienschmidt/httprouter"
)

type api struct {
	Routes   []common.Route
	Contacts []apimodels.Contact // temporary virtual database
}

func NewApi() *api {
	service := api{}
	service.Routes = []common.Route{
		{
			Url:   "contact/:id",
			Method:  common.Get,
			Handler: service.getById,
		},
		{
			Url:   "contact/:id/delete",
			Method:  common.Get,
			Handler: service.deleteById,
		},
		{
			Url:   "contact",
			Method:  common.Post,
			Handler: service.add,
		},

		// Alternate routes
		{
			Url:   "users",
			Method:  common.Get,
			Handler: service.getList,
		},
		{
			Url:   "user",
			Method:  common.Put,
			Handler: service.create,
		},
		{
			Url:   "user/:id",
			Method:  common.Get,
			Handler: service.getById,
		},
		{
			Url:   "user/:id",
			Method:  common.Post,
			Handler: service.updateById,
		},
		{
			Url:   "user/:id",
			Method:  common.Delete,
			Handler: service.deleteById,
		},
	}
	service.Contacts = []apimodels.Contact{}
	return &service
}

func (service *api) getById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idString := strings.TrimSpace(p.ByName("id"))
	id, err := strconv.ParseInt(idString, 10, 32)
	if idString == "" || err != nil {
		badRequest(w, "Identifier '"+idString+"' can not be parsed")
		return
	}

	log.Println("Resolving contact with id: " + idString)

	for _, contact := range service.Contacts {
		if contact.Id == id{
			if err = response(w, contact); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				return
			}
		}
	}
}

func (service *api) getList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("Resolving contacts list")

	if len(service.Contacts) > 0{
		if err := response(w, service.Contacts); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (service *api) deleteById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idString := strings.TrimSpace(p.ByName("id"))
	id, err := strconv.ParseInt(idString, 10, 32)
	if idString == "" || err != nil {
		badRequest(w, "Identifier '"+idString+"' can not be parsed")
		return
	}

	log.Println("Removing contact with id: " + idString)

	for i, contact := range service.Contacts {
		if contact.Id == id{
			service.Contacts = append(service.Contacts[:i], service.Contacts[i+1:]...)
		}
	}
}

func (service *api) updateById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idString := strings.TrimSpace(p.ByName("id"))
	id, err := strconv.ParseInt(idString, 10, 32)
	if idString == "" || err != nil {
		badRequest(w, "Identifier '"+idString+"' can not be parsed")
		return
	}
	decoder := json.NewDecoder(r.Body)
	var contact apimodels.Contact
	err = decoder.Decode(&contact)
	if err != nil {
		badRequest(w, "Failed to parse request body")
	}

	log.Println("Removing contact with id: " + idString)

	for i, c := range service.Contacts {
		if c.Id == id{
			c.Description = contact.Description
			c.Name = contact.Name
			c.Phone = contact.Phone
			c.UpdateAt = time.Now()
			service.Contacts = append(service.Contacts[:i], append(service.Contacts[i+1:], c)...)
		}
	}
}

func (service *api) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var contact apimodels.Contact
	err := decoder.Decode(&contact)
	if err != nil {
		badRequest(w, "Failed to parse request body")
	}

	log.Println("Creating contact.")

	contact.CreatedAt = time.Now()
	contact.UpdateAt = time.Now()

	service.Contacts = append(service.Contacts, contact)
	if err = response(w, nil); err != nil {
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
