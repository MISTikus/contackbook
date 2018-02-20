package contactbook

import (
	"html/template"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/MISTikus/contactbook/apimodels"
	"github.com/MISTikus/contactbook/common"
	"github.com/julienschmidt/httprouter"
)

type view struct {
	Routes       []common.Route
	DefaultRoute common.Route
	TagMap       map[string]map[string]string
}

func NewView() view {
	// ToDo: initialize data from api
	service := view{
		TagMap: getmaps(apimodels.Contact{}),
	}
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

func (v *view) index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !strings.Contains(r.RequestURI, "index") {
		http.Redirect(w, r, "/view/index", 301)
		return
	}

	funcMap := template.FuncMap{
		"getupdatelink": getupdatelink,
		"getdeletelink": getdeletelink,
		"gettag":        v.gettag,
	}

	t := template.Must(template.New("index").Funcs(funcMap).Parse(viewTemplate))

	data := []apimodels.Contact{
		apimodels.Contact{
			Id:          1,
			Name:        "BeforyDeath",
			Description: `Не забудь сказать: "ты что за хуй?!"`,
			Phone:       "123456",
		},
	}

	t.Execute(w, data)
}

func getupdatelink(id int64) string {
	return "/view/contact/" + strconv.FormatInt(id, 10) + "/edit"
}

func getdeletelink(id int64) string {
	return "/api/contact/" + strconv.FormatInt(id, 10) + "/delete"
}

func (v *view) gettag(field string, tag string) string {
	return v.TagMap[field][tag]
}

func getmaps(obj interface{}) map[string]map[string]string {
	t := reflect.TypeOf(obj)
	result := map[string]map[string]string{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		m := map[string]string{}
		for _, element := range strings.Fields(string(f.Tag)) {
			split := strings.Split(element, ":")
			if len(split) > 1 {
				m[split[0]] = strings.Replace(split[1], `"`, "", -1)
			} else {
				m[split[0]] = ""
			}
		}
		result[f.Name] = m
	}
	return result
}

const viewTemplate = `<html>
	<h1 align="center">Контакты</h1>
	<table width="100%" border="1" cellspacing="0">
		<thead>
			<th>{{gettag "Name" "desc"}}</th>
			<th>{{gettag "Phone" "desc"}}</th>
			<th>{{gettag "Description" "desc"}}</th>
			<th>Действия</th>
		</thead>
		<tbody>
			{{range .}}
				<tr>
					<td>{{.Name}}</td>
					<td align="center">{{.Phone}}</td>
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
	<p align="right">
		<a href="/view/contact/add">
			<img width="60" src="/api/common/images/plus" />
		</a>
	</p>
</html>`
