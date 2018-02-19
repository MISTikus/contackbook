package common

import "github.com/julienschmidt/httprouter"

const Get = "GET"
const Post = "POST"
const Put = "PUT"
const Patch = "PATCH"
const Delete = "DELETE"

type Route struct {
	Method  string
	Route   string
	Handler httprouter.Handle
}
