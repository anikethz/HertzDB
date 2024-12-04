package v1

import (
	"github.com/anikethz/HertzDB/src/web/types"
	"github.com/anikethz/HertzDB/src/web/v1/ingest"
	"github.com/anikethz/HertzDB/src/web/v1/search"
	"github.com/go-chi/chi/v5"
)

func GetV1Router() *chi.Mux {

	v1Router := chi.NewRouter()

	apiConfig := types.ApiConfig{}

	v1Router.Get("/{index}/search", apiConfig.Handler(search.SearchHandler))
	v1Router.Post("/{index}/ingest", apiConfig.Handler(ingest.IngestionHandler))

	return v1Router

}
