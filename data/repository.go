package data

import (
	"errors"
	"github.com/MISTikus/contactbook/apimodels"
	"strconv"
	"time"
)

type Repository struct {
	file     *FileWorker
	maxId    int64
	contacts []apimodels.Contact
}

func NewRepo(file *FileWorker) *Repository {
	repo := Repository{
		file:     file,
		contacts: []apimodels.Contact{},
	}
	file.Read(&repo.contacts)

	if len(repo.contacts) == 0 {
		repo.contacts = append(repo.contacts, apimodels.Contact{
			Id:          1,
			Name:        "Byd",
			Description: "Some whoi",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Phone:       "no phone",
		})
		repo.save()
		repo.maxId = 1
	}

	for _, contact := range repo.contacts {
		if contact.Id > repo.maxId {
			repo.maxId = contact.Id
		}
	}

	return &repo
}

func (repo *Repository) GetAll() []apimodels.Contact {
	return repo.contacts
}

func (repo *Repository) Get(id int64) (apimodels.Contact, error) {
	for _, contact := range repo.contacts {
		if contact.Id == id {
			return contact, nil
		}
	}
	return apimodels.Contact{}, errors.New("Contact with id '" + strconv.FormatInt(id, 10) + "' not found.")
}

func (repo *Repository) Add(contact apimodels.Contact) {
	contact.Id = repo.getNewId()
	repo.contacts = append(repo.contacts, contact)
	repo.save()
}

func (repo *Repository) Delete(id int64) {
	for i, contact := range repo.contacts {
		if contact.Id == id {
			repo.contacts = append(repo.contacts[:i], repo.contacts[i+1:]...)
			repo.save()
			return
		}
	}
}

func (repo *Repository) Update(contact apimodels.Contact) {
	for i, c := range repo.contacts {
		if c.Id == contact.Id {
			repo.contacts = append(append(repo.contacts[:i], contact), repo.contacts[i+1:]...)
			repo.save()
			return
		}
	}
}

func (repo *Repository) save() {
	repo.file.Write(repo.contacts)
}

func (repo *Repository) getNewId() int64 {
	id := repo.maxId + 1
	repo.maxId = id
	return id
}
