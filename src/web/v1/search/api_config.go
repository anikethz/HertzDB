package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ApiConfig struct {
	Index         string
	Filename      string
	Json_Filename string
}

func (apiConfig *ApiConfig) Handler(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		apiConfig.Index = chi.URLParam(r, "index")
		apiConfig.Filename = apiConfig.Index + ".hz"
		apiConfig.Json_Filename = apiConfig.Index + ".json"
		f(w, r)
	}

}
