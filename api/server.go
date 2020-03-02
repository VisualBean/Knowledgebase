package api

import (
	"fmt"
	"log"
	"net/http"

	"knowledgebase/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	dbName = "KB"
)

type Server struct {
	Database *gorm.DB
	Router   *chi.Mux
}

func initializeRoutes(server *Server) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.RealIP,
	)

	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/kb", InitializeKnowledgebaseApi(server))
	})

	return router
}

func (server *Server) Initialize(dbUser string, dbPassword string, dbAddress string) {
	var err error

	connectionString := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbName)
	server.Database, err = gorm.Open("mysql", connectionString)

	if err != nil {
		log.Fatal("Cannot connect to database: " + err.Error())
	}

	models.DBMigrate(server.Database)

	server.Router = initializeRoutes(server)

	printRoutes := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(server.Router, printRoutes); err != nil {
		log.Panicln("Logging err: %s", err.Error())
	}
}

func (server *Server) Start(address string) {
	fmt.Printf("Starting server at: %s", address)
	log.Fatal(http.ListenAndServe(address, server.Router))
}
