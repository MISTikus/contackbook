package contactbook

import (
	"encoding/json"
	"github.com/MISTikus/contactbook/data"
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
	repository    *data.Repository
	ChangeHandler ChangeHandler
}

type ChangeHandler interface {
	HandleChanges()
}

func NewApi(repository *data.Repository) *api {
	service := api{
		Prefix:     "users",
		repository: repository,
	}
	service.Routes = []common.Route{
		{
			Url:     "",
			Method:  common.Get,
			Handler: service.getList,
		},
		{
			Url:     "",
			Method:  common.Put,
			Handler: service.create,
		},
		{
			Url:     ":id",
			Method:  common.Get,
			Handler: service.getById,
		},
		{
			Url:     ":id",
			Method:  common.Post,
			Handler: service.updateById,
		},
		{
			Url:     ":id",
			Method:  common.Delete,
			Handler: service.deleteById,
		},
	}

	return &service
}

func (service *api) GetContacts() []apimodels.Contact {
	return service.repository.GetAll()
}

func (service *api) getById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idString := strings.TrimSpace(p.ByName("id"))
	id, err := strconv.ParseInt(idString, 10, 32)
	if idString == "" || err != nil {
		badRequest(w, "Identifier '"+idString+"' can not be parsed")
		return
	}

	log.Println("Resolving contact with id: " + idString)
	c, _ := service.repository.Get(id)

	if err = response(w, c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		return
	}
}

func (service *api) getList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("Resolving contacts list")

	if err := response(w, service.repository.GetAll()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	service.repository.Delete(id)
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
	contact.Id = id
	contact.UpdatedAt = time.Now()

	log.Println("Removing contact with id: " + idString)

	service.repository.Update(contact)
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

	service.repository.Add(contact)
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
