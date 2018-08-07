package controller

import (
	"fmt"
	"github.com/brave/go-update/extension"
	"github.com/go-chi/chi"
	"github.com/pressly/lg"
	"io/ioutil"
	"net/http"
)

var allExtensions extension.Extensions

// ExtensionsRouter is the router for /extensions endpoints
func ExtensionsRouter(extensions extension.Extensions) chi.Router {
	allExtensions = extensions
	r := chi.NewRouter()
	r.Post("/", UpdateExtensions)
	return r
}

// UpdateExtensions is the handler for updating extensions
func UpdateExtensions(w http.ResponseWriter, r *http.Request) {
	log := lg.Log(r.Context())
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Errorf("Error closing body stream: %v", err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading body: %v", err), http.StatusBadRequest)
		return
	}

	extensions, err := extension.UnmarshalXML(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading body %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/xml")
	w.WriteHeader(http.StatusOK)

	extensions = extension.FilterForUpdates(allExtensions, extensions)
	data, err := extension.MarshalXML(extensions)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in marshal XML %v", err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		log.Errorf("Error writing response: %v", err)
	}
}
