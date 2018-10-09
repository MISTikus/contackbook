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
	changed      bool
}

func NewView(/*apiService *api, commonService *common.Api*/) *view {
	// ToDo: initialize data from api
	service := view{
		Prefix: "",
		//Api: apiService,
		//Common: commonService,
		tagMap: getmaps(apimodels.Contact{}),
	}
	service.Routes = []common.Route{
		{
			Url:   "list",
			Method:  common.Get,
			Handler: service.index,
		},
		{
			Url:   "update/:id",
			Method:  common.Get,
			Handler: service.update,
		},
		{
			Url:   "delete/:id",
			Method:  common.Get,
			Handler: service.delete,
		},
		{
			Url:   "create",
			Method:  common.Get,
			Handler: service.create,
		},
		{
			Url: "hasChanges",
			Method: common.Get,
			Handler: service.hasChanges,
		},
	}
	service.DefaultRoute = service.Routes[0]
	return &service
}

func (v *view) index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !strings.Contains(r.RequestURI, "list") {
		http.Redirect(w, r, common.BuildUrl("view", v.Prefix, "list"), 307)
		return
	}

	funcMap := template.FuncMap{
		"getupdatelink": v.getupdatelink,
		"getdeletelink": v.getdeletelink,
		"getcreatelink": v.getcreatelink,
		"getimagelink": v.getimagelink,
		"formatdatetime": v.formatdatetime,
		"gethaschangeslink": v.gethaschangeslink,
		"gettag":        v.gettag,
	}

	t := template.Must(template.New("list").Funcs(funcMap).Parse(indexTemplate))

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

	funcMap := template.FuncMap{
		"getupdateurl": v.getupdateurl,
		"formatdatetime": v.formatdatetime,
		"gettag":        v.gettag,
	}

	t := template.Must(template.New("update").Funcs(funcMap).Parse(updateTemplate))

	for _, contact := range v.Api.contacts {
		if contact.Id == id {
			t.Execute(w, contact)
			break
		}
	}
}

func (v *view) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	funcMap := template.FuncMap{
		"getcreateurl": v.getcreateurl,
		"gettag":        v.gettag,
	}

	t := template.Must(template.New("create").Funcs(funcMap).Parse(createTemplate))

	t.Execute(w, nil)
}

func (v *view) delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idString := strings.TrimSpace(p.ByName("id"))
	id, err := strconv.ParseInt(idString, 10, 32)
	if idString == "" || err != nil {
		badRequest(w, "Identifier '"+idString+"' can not be parsed")
		return
	}

	funcMap := template.FuncMap{
		"getdeleteurl": v.getdeleteurl,
		"gettag":        v.gettag,
	}

	t := template.Must(template.New("delete").Funcs(funcMap).Parse(deleteTemplate))

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
	v.changed = true;
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
	return  common.BuildUrl("view", v.Prefix, "hasChanges")
}

func (v *view) getimagelink(key string) string{
	return common.BuildUrl("api", v.Common.Prefix, "images", key)
}

func (v *view) getupdateurl(id int64) string{
	return common.BuildUrl("api", v.Api.Prefix, strconv.FormatInt(id, 10))
}

func (v *view) getcreateurl() string{
	return common.BuildUrl("api", v.Api.Prefix)
}

func (v *view) getdeleteurl(id int64) string{
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

func equals(a, b []apimodels.Contact) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false;
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}


const indexTemplate = `<html>
	<head>
		<script>
`+httpClientScript+`
`+longPoolingScript+`
		</script>
	</head>
	<h1 align="center">Контакты</h1>
	<table width="100%" border="1" cellspacing="0">
		<thead>
			<th>{{gettag "Name" "desc"}}</th>
			<th>{{gettag "Phone" "desc"}}</th>
			<th>{{gettag "CreatedAt" "desc"}}</th>
			<th>{{gettag "UpdatedAt" "desc"}}</th>
			<th>{{gettag "Description" "desc"}}</th>
			<th>Действия</th>
		</thead>
		<tbody>
			{{range .}}
				<tr>
					<td>{{.Name}}</td>
					<td align="center">{{.Phone}}</td>
					<td align="center">{{.CreatedAt | formatdatetime}}</td>
					<td align="center">{{.UpdatedAt | formatdatetime}}</td>
					<td>{{.Description}}</td>
					<td>
						<a href="{{.Id | getupdatelink}}" target="_blank">Изменить</a>
						<br/>
						<a href="{{.Id | getdeletelink}}" target="_blank">Удалить</a>
					</td>
				</tr>
			{{end}}
		</tbody>
	</table>
	<p align="right">
		<a href="{{getcreatelink}}"  target="_blank">
			<img width="60" src="{{"plus" | getimagelink}}" />
		</a>
	</p>
</html>`

const updateTemplate = `<html>
	<head>
		<script>
`+httpClientScript+`
document.addEventListener('DOMContentLoaded', function(){
	var form = document.querySelector('form');
	form.addEventListener('submit', function onSubmit(event) {
    	event.preventDefault();
    	http.post(form.action, toJSONString(form), {
    		onsuccess: function() { window.close() },
      		onerror: function() { alert(error); window.close(); }
    	});
	});
});
		</script>
	</head>
	<h1>Изменить</h1>
	<form action="{{.Id | getupdateurl}}" method="POST">
		<input type="hidden" name="createdAt" value="{{.CreatedAt | formatdatetime}}" />
		<label>{{gettag "Name" "desc"}}</label><input name="name" type="text" value="{{.Name}}" /><br />
		<label>{{gettag "Phone" "desc"}}</label><input name="phone" type="text" value="{{.Phone}}" /><br />
		<label>{{gettag "Description" "desc"}}</label><textarea name="desc" cols="40" rows="3">{{.Description}}</textarea><br />
		<input type="submit" value="Submit" />
	</form>
</html>`

const createTemplate = `<html>
	<head>
		<script>
`+httpClientScript+`
document.addEventListener('DOMContentLoaded', function(){
	var form = document.querySelector('form');
	form.addEventListener('submit', function onSubmit(event) {
    	event.preventDefault();
		http.put(form.action, toJSONString(form), {
	  		onsuccess: function() { window.close() },
	    	onerror: function() { alert(error); window.close(); }
	    });
	});
});
		</script>
	</head>
	<h1>Создать</h1>
	<form action="{{getcreateurl}}"  method="PUT">
		<label>{{gettag "Name" "desc"}}</label><input name="name" type="text" value="{{.Name}}" /><br />
		<label>{{gettag "Phone" "desc"}}</label><input name="phone" type="text" value="{{.Phone}}" /><br />
		<label>{{gettag "Description" "desc"}}</label><textarea name="desc" cols="40" rows="3">{{.Description}}</textarea><br />
		<input type="submit" value="Submit" />
	</form>
</html>`

const deleteTemplate = `<html>
	<head>
		<script>
`+httpClientScript+`
document.addEventListener('DOMContentLoaded', function(){
	var form = document.querySelector('form');
	form.addEventListener('submit', function onSubmit(event) {
	    event.preventDefault();
	    http.delete(form.action, {
	    	onsuccess: function() { window.close() },
	      	onerror: function() { alert(error); window.close(); }
    	});
	});
});
		</script>
	</head>
	<h1>Удалить</h1>
	<h3>Уверены, что хотите удалить контакт с идентиифкатором {{.Id}}?</h3>
	<form action="{{.Id | getdeleteurl}}" method="DELETE">
		<input type="submit" value="Delete" /><br />
		<input type="button" value="Cancel" onclick="window.close()" />
	</form>
</html>`

const httpClientScript = `/// HTTP CLIENT ///

var http = (function () {
  function HttpClient() {
    this.xhr = new XMLHttpRequest();
    this.async = true;
  }

  function sendRequest(method, url, body, callbacks) {
    var self = this;
    this.xhr.open(method, url, this.async);

    if (method === 'POST') {
      this.xhr.setRequestHeader('Content-type', 'application/json');
    }

    this.xhr.send(body);
    this.xhr.onreadystatechange = function onRequestStateChange() {
      var done = 4;
      var ok = 200;
      if (self.xhr.readyState != done) {
        return;
      }
      if (self.xhr.status !== ok && self.xhr.status !== 201) {
        callbacks.onerror({ status: self.xhr.status, statusText: self.xhr.statusText });
        return;
      } else {
        callbacks.onsuccess(self.xhr.response ? JSON.parse(self.xhr.response) : null, self.xhr.status);
      }
    }
  }

  HttpClient.prototype.get = function sendGet(url, callbacks) {
    sendRequest.call(this, 'GET', url, null, callbacks);
  }

  HttpClient.prototype.post = function sendPost(url, body, callbacks) {
    sendRequest.call(this, 'POST', url, body, callbacks);
  }

  HttpClient.prototype.delete = function sendDelete(url, callbacks) {
    sendRequest.call(this, 'DELETE', url, null, callbacks);
  }

  HttpClient.prototype.put = function sendPut(url, body, callbacks) {
    sendRequest.call(this, 'PUT', url, body, callbacks);
  }

  return new HttpClient();
})();

var toJSONString = (function ( form ) {
	var obj = {};
	var elements = form.querySelectorAll( "input, select, textarea" );
	for( var i = 0; i < elements.length; ++i ) {
		var element = elements[i];
		var name = element.name;
		var value = element.value;

		if( name ) {
			obj[ name ] = value;
		}
	}
	return JSON.stringify( obj );
});`

const longPoolingScript = `/// Long pooling ///
document.addEventListener('DOMContentLoaded', function(){
	var pooling = (function(){
		window.setTimeout(function(){
			http.get({{gethaschangeslink}}, {
				onsuccess: function(obj, status) {
					if (status === 201){
						location.reload();
					} else {
						console.log("No changes, continue...")
						pooling();
					}
				},
				onerror: function(err){ console.log(err); }
			});
		},1000);
	});
	pooling();
});`