package v1Router

import (
	v1Search "github.com/anikethz/HertzDB/src/web/v1/search"
	"github.com/go-chi/chi/v5"
)

func GetV1Router() *chi.Mux {

	v1Router := chi.NewRouter()

	v1Router.Get("/{index}/search", v1Search.SearchHandler)

	return v1Router

}
