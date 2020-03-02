package responses

import (
	"net/http"

	"github.com/go-chi/render"
)

func OK(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.JSON(w, r, data)
}
func CREATED(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, 201)
	OK(w, r, data)
}

func ERROR(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	render.Status(r, statusCode)
	render.JSON(w, r, struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}
