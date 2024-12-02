package v1Router

import (
	v1Search "github.com/anikethz/HertzDB/src/web/v1/search"
	"github.com/go-chi/chi/v5"
)

func GetV1Router() *chi.Mux {

	v1Router := chi.NewRouter()

	apiConfig := v1Search.ApiConfig{}

	v1Router.Get("/{index}/search", apiConfig.Handler(apiConfig.SearchHandler))

	return v1Router

}
