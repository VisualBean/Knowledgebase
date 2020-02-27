package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
)

const (
	dbName = "knowledgebase"
)

type Server struct {
	DB     *gorm.DB
	Router *chi.Mux
}

func initializeRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.RealIP,
	)

	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/kb", KnowledgebaseRoutes())
	})

	return router
}

func (server *Server) Initialize(dbUser string, dbPassword string, dbAddress string) {
	var err error

	connectionString := fmt.Sprintf("%s:%s@tcp(%s/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbAddress, dbName)
	server.DB, err = gorm.Open("mysql", connectionString)

	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	server.Router = initializeRoutes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(server.Router, walkFunc); err != nil {
		log.Panicln("Logging err: %s", err.Error())
	}
}

func (server *Server) Start(address string) {
	log.Fatal(http.ListenAndServe(address, server.Router))
}
