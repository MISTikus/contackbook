package domain

type contact struct {
	Name        string `desc:"Имя" length:"255"`
	Phone       string `desc:"Телефон" length:"150" mask:"\d+"`
	Description string `desc:"Примечание" length:"510"`
}
