package contactbook

import (
	"html/template"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/MISTikus/contactbook/apimodels"
	"github.com/MISTikus/contactbook/common"
	"github.com/julienschmidt/httprouter"
)

type view struct {
	Prefix       string
	Api          *api
	Common       *common.Api
	Routes       []common.Route
	DefaultRoute common.Route
	tagMap       map[string]map[string]string
	funcMap      map[string]interface{}
	changed      bool
}

func NewView() *view {
	// ToDo: initialize data from api
	service := view{
		Prefix: "",
		tagMap: getmaps(apimodels.Contact{}),
	}
	service.funcMap = service.getFuncMap()
	service.Routes = []common.Route{
		{
			Url:     "list",
			Method:  common.Get,
			Handler: service.index,
		},
		{
			Url:     "update/:id",
			Method:  common.Get,
			Handler: service.update,
		},
		{
			Url:     "delete/:id",
			Method:  common.Get,
			Handler: service.delete,
		},
		{
			Url:     "create",
			Method:  common.Get,
			Handler: service.create,
		},
		{
			Url:     "hasChanges",
			Method:  common.Get,
			Handler: service.hasChanges,
		},
	}
	service.DefaultRoute = service.Routes[0]
	return &service
}
func (v *view) getFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"getupdatelink":     v.getupdatelink,
		"getdeletelink":     v.getdeletelink,
		"getcreatelink":     v.getcreatelink,
		"getimagelink":      v.getimagelink,
		"formatdatetime":    v.formatdatetime,
		"gethaschangeslink": v.gethaschangeslink,
		"gettag":            v.gettag,
		"getupdateurl":      v.getupdateurl,
		"getcreateurl":      v.getcreateurl,
		"getdeleteurl":      v.getdeleteurl,
	}
}

func (v *view) index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !strings.Contains(r.RequestURI, "list") {
		http.Redirect(w, r, common.BuildUrl("view", v.Prefix, "list"), 307)
		return
	}

	t := template.Must(template.New("list").Funcs(v.funcMap).ParseGlob("views/" + "*.t*"))

	data := v.Api.GetContacts()
	v.changed = false;
	t.Execute(w, data)
}

func (v *view) update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idString := strings.TrimSpace(p.ByName("id"))
	id, err := strconv.ParseInt(idString, 10, 32)
	if idString == "" || err != nil {
		badRequest(w, "Identifier '"+idString+"' can not be parsed")
		return
	}

	t := template.Must(template.New("update").Funcs(v.funcMap).ParseGlob("views/" + "*.t*"))

	for _, contact := range v.Api.contacts {
		if contact.Id == id {
			t.Execute(w, contact)
			break
		}
	}
}

func (v *view) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t := template.Must(template.New("create").Funcs(v.funcMap).ParseGlob("views/" + "*.t*"))

	t.Execute(w, nil)
}

func (v *view) delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idString := strings.TrimSpace(p.ByName("id"))
	id, err := strconv.ParseInt(idString, 10, 32)
	if idString == "" || err != nil {
		badRequest(w, "Identifier '"+idString+"' can not be parsed")
		return
	}

	t := template.Must(template.New("delete").Funcs(v.funcMap).ParseGlob("views/" + "*.t*"))

	for _, contact := range v.Api.contacts {
		if contact.Id == id {
			t.Execute(w, contact)
			break
		}
	}
}

func (v *view) hasChanges(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if v.changed {
		w.WriteHeader(201)
	} else {
		w.WriteHeader(200)
	}
}

func (v *view) HandleChanges() {
	v.changed = true
}

func (v *view) gettag(field string, tag string) string {
	return v.tagMap[field][tag]
}

func (v *view) getupdatelink(id int64) string {
	return common.BuildUrl("view", v.Prefix, "update", strconv.FormatInt(id, 10))
}

func (v *view) getdeletelink(id int64) string {
	return common.BuildUrl("view", v.Prefix, "delete", strconv.FormatInt(id, 10))
}

func (v *view) getcreatelink() string {
	return common.BuildUrl("view", v.Prefix, "create")
}

func (v *view) gethaschangeslink() string {
	return common.BuildUrl("view", v.Prefix, "hasChanges")
}

func (v *view) getimagelink(key string) string {
	return common.BuildUrl("api", v.Common.Prefix, "images", key)
}

func (v *view) getupdateurl(id int64) string {
	return common.BuildUrl("api", v.Api.Prefix, strconv.FormatInt(id, 10))
}

func (v *view) getcreateurl() string {
	return common.BuildUrl("api", v.Api.Prefix)
}

func (v *view) getdeleteurl(id int64) string {
	return common.BuildUrl("api", v.Api.Prefix, strconv.FormatInt(id, 10))
}

func (v *view) formatdatetime(value time.Time) string {
	return value.Format(time.RFC3339)
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
