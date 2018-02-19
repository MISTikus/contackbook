package common

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/MISTikus/gotalotoftime/common"
	"github.com/julienschmidt/httprouter"
)

type api struct {
	Routes []common.Route
}

func NewApi() api {
	service := api{}
	service.Routes = []common.Route{
		{
			Route:   "common/images/:imageName",
			Method:  common.Get,
			Handler: service.getImage,
		},
	}
	return service
}

func (service api) getImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idString := strings.TrimSpace(p.ByName("imageName"))
	if idString == "" {
		badRequest(w, "Identifier '"+idString+"' can not be parsed")
		return
	}

	log.Println("Resolving image with id: " + idString)

	fImg1, _ := os.Open("common/" + idString + ".png")
	defer fImg1.Close()
	img, _ := png.Decode(fImg1)

	writeImage(w, &img)
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

func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := png.Encode(w, *img); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}
