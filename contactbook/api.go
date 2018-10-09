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
	Prefix        string
	maxId         int64
	Routes        []common.Route
	contacts      []apimodels.Contact // temporary virtual database
	ChangeHandler ChangeHandler
}

type ChangeHandler interface {
	HandleChanges()
}

func NewApi(/*changeHandler *ChangeHandler*/) *api {
	service := api{
		Prefix:"users",
		//ChangeHandler: changeHandler
	}
	service.Routes = []common.Route{
		{
			Url:   "",
			Method:  common.Get,
			Handler: service.getList,
		},
		{
			Url:   "",
			Method:  common.Put,
			Handler: service.create,
		},
		{
			Url:   ":id",
			Method:  common.Get,
			Handler: service.getById,
		},
		{
			Url:   ":id",
			Method:  common.Post,
			Handler: service.updateById,
		},
		{
			Url:   ":id",
			Method:  common.Delete,
			Handler: service.deleteById,
		},
	}
	service.contacts = []apimodels.Contact{
		{
			Id:          1,
			Name:        "BeforyDeath",
			Description: `Не забудь сказать: "ты что за хуй?!"`,
			Phone:       "123456",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	service.maxId = 1
	return &service
}

func (service *api) GetContacts() []apimodels.Contact {
	return service.contacts
}

func (service *api) getById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idString := strings.TrimSpace(p.ByName("id"))
	id, err := strconv.ParseInt(idString, 10, 32)
	if idString == "" || err != nil {
		badRequest(w, "Identifier '"+idString+"' can not be parsed")
		return
	}

	log.Println("Resolving contact with id: " + idString)

	for _, contact := range service.contacts {
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

	if len(service.contacts) > 0{
		if err := response(w, service.contacts); err != nil {
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

	for i, contact := range service.contacts {
		if contact.Id == id{
			service.contacts = append(service.contacts[:i], service.contacts[i+1:]...)
		}
	}
	service.ChangeHandler.HandleChanges()
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
		badRequest(w, "Failed to parse body")
		return
	}
	contact.UpdatedAt = time.Now()

	log.Println("Removing contact with id: " + idString)

	for i, c := range service.contacts {
		if c.Id == id{
			c.Description = contact.Description
			c.Name = contact.Name
			c.Phone = contact.Phone
			c.UpdatedAt = time.Now()
			service.contacts = append(append(service.contacts[:i], c), service.contacts[i+1:]...)
		}
	}
	service.ChangeHandler.HandleChanges()
}

func (service *api) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var contact apimodels.Contact
	err := decoder.Decode(&contact)
	if err != nil {
		badRequest(w, "Failed to parse body")
	}
	var id = service.maxId + 1
	service.maxId = id
	contact.Id = id
	contact.CreatedAt = time.Now()
	contact.UpdatedAt = time.Now()

	log.Println("Creating contact.")

	service.contacts = append(service.contacts, contact)
	if err = response(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	service.ChangeHandler.HandleChanges()
}

func badRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
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
