package apimodels

import "time"

type Contact struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name" desc:"Имя" length:"255"`
	Phone       string    `json:"phone" desc:"Телефон" length:"150" mask:"\d+"`
	Description string    `json:"desc" desc:"Примечание" length:"510"`
	CreatedAt   time.Time `json:"createdAt" desc:"Создано"`
	UpdatedAt   time.Time `json:"updatedAt" desc:"Изменено"`
}
