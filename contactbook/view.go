package contactbook

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/MISTikus/contactbook/apimodels"

	"github.com/MISTikus/contactbook/common"
	"github.com/julienschmidt/httprouter"
)

type view struct {
	Routes       []common.Route
	DefaultRoute common.Route
}

func NewView() view {
	// ToDo: initialize data from api
	service := view{}
	service.Routes = []common.Route{
		{
			Route:   "index",
			Method:  common.Get,
			Handler: service.index,
		},
	}
	service.DefaultRoute = service.Routes[0]
	return service
}

func (c *view) index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !strings.Contains(r.RequestURI, "index") {
		http.Redirect(w, r, "/view/index", 301)
		return
	}

	funcMap := template.FuncMap{
		"getupdatelink": getupdatelink,
		"getdeletelink": getdeletelink,
	}

	t := template.Must(template.New("index").Funcs(funcMap).Parse(viewTemplate))

	data := []apimodels.Contact{
		apimodels.Contact{
			Id:          1,
			Name:        "BeforyDeath",
			Description: "Не забудь сказать: ты что за хуй?!",
			Phone:       "123456",
		},
	}

	t.Execute(w, data)
}

func getupdatelink(id int64) string {
	return "/api/contact/delete/" + strconv.FormatInt(id, 10)
}

func getdeletelink(id int64) string {
	return "/view/contact/edit/" + strconv.FormatInt(id, 10)
}

const viewTemplate = `<html>
	<h1 align="center">Справочник</h1>
	<table>
		<thead>
			<th>ФИО</th>
			<th>Номер телефона</th>
			<th>Примечание</th>
			<th>Действия</th>
		</thead>
		<tbody>
			{{range .}}
				<tr>
					<td>{{.Name}}</td>
					<td>{{.Phone}}</td>
					<td>{{.Description}}</td>
					<td>
						<a href="{{.Id | getupdatelink}}">Изменить</a>
						<br/>
						<a href="{{.Id | getdeletelink}}">Удалить</a>
					</td>
				</tr>
			{{end}}
		</tbody>
	</table>
</html>`
