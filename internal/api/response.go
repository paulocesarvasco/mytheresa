package api

import (
	"net/http"

	"github.com/go-chi/render"
)

func OKResponse(w http.ResponseWriter, r *http.Request, data any) {
	OKResponseWithStatus(w, r, http.StatusOK, data)
}

func OKResponseWithStatus(w http.ResponseWriter, r *http.Request, status int, data any) {
	render.Status(r, status)
	render.SetContentType(render.ContentTypeJSON)
	render.JSON(w, r, data)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	payload := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}
	render.Status(r, status)
	render.SetContentType(render.ContentTypeJSON)
	render.JSON(w, r, payload)
}
