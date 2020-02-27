package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	. "knowledgebase/models"
)

// Routes of this set of endpoints.
func KnowledgebaseRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", getAllDocuments)
	router.Get("/{id}", getNewestDocumentByID)
	router.Get("/{id}/version/{version}", getDocumentVersion)

	router.Put("/{id}", updateDocument)

	router.Post("/", createDocument)

	router.Delete("/{id}/version/{version}", deleteVersion)
	return router
}

func getAllDocuments(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "test")
}

func getNewestDocumentByID(w http.ResponseWriter, r *http.Request) {
	// documentID, err := getIDParam(r)

	//if err != nil {
	//	render.Status(r, 400)
	//}

	document := Entry{
		ID: 123,
		Elements: []Element{
			Element{Text: "This is a test text", Type: "H1"},
		},
		Version: 1,
	}
	render.JSON(w, r, document)
}
func getDocumentVersion(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "test")
}
func updateDocument(w http.ResponseWriter, r *http.Request) {
	// Disregard version and ID fields
}

func createDocument(w http.ResponseWriter, r *http.Request) {
	// Disregard version and id field.
	render.JSON(w, r, "test")
}

func deleteVersion(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "test")
}

func getIDParam(r *http.Request) (int64, error) {
	output := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(output, 10, 0)
	return i, err
}
