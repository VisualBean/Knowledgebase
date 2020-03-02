package api

import (
	"encoding/json"
	. "knowledgebase/models"
	"knowledgebase/responses"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var (
	database *gorm.DB
)

func InitializeKnowledgebaseApi(server *Server) *chi.Mux {
	database = server.Database

	router := chi.NewRouter()
	router.Get("/", getAllDocuments)
	router.Get("/{id}", getNewestDocumentByID)
	router.Get("/{id}/versions/{version}", getDocumentVersion)
	router.Get("/{id}/versions", getDocumentVersions)

	router.Put("/{id}", updateDocument)

	router.Post("/", createDocument)

	router.Delete("/{id}/versions/{version}", deleteVersion)
	return router
}

func getAllDocuments(w http.ResponseWriter, r *http.Request) {
	entry := &Entry{}

	entries, err := entry.GetAllEntries(database)
	if err != nil {
		responses.ERROR(w, r, http.StatusInternalServerError, err)
		return
	}

	responses.OK(w, r, entries)
}

func getDocumentVersions(w http.ResponseWriter, r *http.Request) {
	uuid, _, err := getAndParseParameters(r)
	if err != nil {
		responses.ERROR(w, r, http.StatusBadRequest, nil)
		return
	}

	entry := &Entry{}
	entries, err := entry.GetVersions(database, uuid)

	if err != nil {
		if err == ErrNotFound {
			responses.ERROR(w, r, http.StatusNotFound, err)
			return
		}
		responses.ERROR(w, r, http.StatusInternalServerError, nil)
		return
	}

	responses.OK(w, r, entries)
}

func getNewestDocumentByID(w http.ResponseWriter, r *http.Request) {
	uuid, _, err := getAndParseParameters(r)
	if err != nil {
		responses.ERROR(w, r, http.StatusBadRequest, nil)
		return
	}

	entry := &Entry{}
	entry, err = entry.GetNewestVersion(database, uuid)
	if err != nil {
		responses.ERROR(w, r, http.StatusInternalServerError, nil)
		return
	}

	responses.OK(w, r, entry)
}

func getDocumentVersion(w http.ResponseWriter, r *http.Request) {
	uuid, version, err := getAndParseParameters(r)
	if err != nil {
		responses.ERROR(w, r, http.StatusBadRequest, nil)
		return
	}

	entry := &Entry{}
	entry, err = entry.GetVersion(database, uuid, version)

	if err != nil {
		if err == ErrNotFound {
			responses.ERROR(w, r, http.StatusNotFound, err)
			return
		}

		responses.ERROR(w, r, http.StatusInternalServerError, nil)
		return
	}

	responses.OK(w, r, entry)
}

func updateDocument(w http.ResponseWriter, r *http.Request) {
	uuid, _, err := getAndParseParameters(r)
	if err != nil {
		responses.ERROR(w, r, http.StatusBadRequest, nil)
		return
	}

	entry := &Entry{}
	json.NewDecoder(r.Body).Decode(&entry)

	if err := entry.Validate(); err != nil {
		responses.ERROR(w, r, http.StatusBadRequest, err)
		return
	}

	entry.UUID = uuid

	entry, err = entry.UpdateEntry(database, uuid)
	if err != nil {
		if err == ErrNotFound {
			responses.ERROR(w, r, http.StatusNotFound, err)
			return
		}

		responses.ERROR(w, r, http.StatusInternalServerError, err)
		return
	}

	responses.OK(w, r, entry)
}

func createDocument(w http.ResponseWriter, r *http.Request) {
	entry := &Entry{}
	json.NewDecoder(r.Body).Decode(&entry)

	if err := entry.Validate(); err != nil {
		responses.ERROR(w, r, http.StatusBadRequest, err)
		return
	}

	created, err := entry.CreateEntry(database)
	if err != nil {
		responses.ERROR(w, r, http.StatusInternalServerError, nil)
		return
	}

	responses.CREATED(w, r, &created)
}

func deleteVersion(w http.ResponseWriter, r *http.Request) {
	uuid, version, err := getAndParseParameters(r)
	if err != nil {
		responses.ERROR(w, r, http.StatusBadRequest, nil)
		return
	}

	entry := &Entry{}
	_, err = entry.DeleteVersion(database, uuid, version)

	if err != nil {
		if err == ErrNotFound {
			responses.ERROR(w, r, http.StatusNotFound, err)
			return
		}
		responses.ERROR(w, r, http.StatusInternalServerError, nil)
		return
	}

	responses.OK(w, r, nil) // Could also be not found.
}

func getAndParseParameters(r *http.Request) (uuid.UUID, uint16, error) {
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))

	converted, err := strconv.ParseUint(chi.URLParam(r, "version"), 16, 16)

	var version uint16
	version = uint16(converted)

	return uuid, version, err
}
