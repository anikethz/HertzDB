package web

import (
	"log"
	"net/http"
	"os"

	v1Router "github.com/anikethz/HertzDB/src/web/v1"
	"github.com/go-chi/chi/v5"
)

func exportDefaultRouter() *chi.Mux {

	router := chi.NewRouter()
	router.Mount("/v1", v1Router.GetV1Router())
	return router

}

func StartServer() *http.Server {

	PORT := os.Getenv("PORT")
	server := &http.Server{
		Handler: exportDefaultRouter(),
		Addr:    ":" + PORT,
	}

	log.Printf("Sever started at PORT %v", PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
		return nil
	}

	return server
}
