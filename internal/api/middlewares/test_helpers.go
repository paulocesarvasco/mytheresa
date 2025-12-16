package middlewares

import (
	"net/http"

	"github.com/go-chi/render"
)

type fakeHandler struct {
	respCode int
	respBody any
}

func NewFakeHandler(code int, body any) *fakeHandler {
	return &fakeHandler{
		respCode: code,
		respBody: body,
	}
}

func (f *fakeHandler) WriteJSON(w http.ResponseWriter, r *http.Request) {
	render.Status(r, f.respCode)
	render.SetContentType(render.ContentTypeJSON)
	render.JSON(w, r, f.respBody)
}
