// handler/author_handler: HTTP handlers for author CRUD.
package handler

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/aliakbar-zohour/go_blog/internal/service"
	"github.com/aliakbar-zohour/go_blog/pkg/response"
	"github.com/go-chi/chi/v5"
)

type AuthorHandler struct {
	svc *service.AuthorService
}

func NewAuthorHandler(svc *service.AuthorService) *AuthorHandler {
	return &AuthorHandler{svc: svc}
}

func (h *AuthorHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.List(r.Context())
	if err != nil {
		response.Internal(w, "failed to list authors")
		return
	}
	response.OK(w, list)
}

func (h *AuthorHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.BadRequest(w, "invalid request format")
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	var avatar *multipart.FileHeader
	if r.MultipartForm != nil && len(r.MultipartForm.File["avatar"]) > 0 {
		avatar = r.MultipartForm.File["avatar"][0]
	}
	a, err := h.svc.Create(r.Context(), name, avatar)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, a)
}

func (h *AuthorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	a, err := h.svc.GetByID(r.Context(), uint(id))
	if err != nil {
		response.NotFound(w, "author not found")
		return
	}
	if a == nil {
		response.NotFound(w, "author not found")
		return
	}
	response.OK(w, a)
}

func (h *AuthorHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	var name string
	var avatar *multipart.FileHeader
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			response.BadRequest(w, "invalid request format")
			return
		}
		name = strings.TrimSpace(r.FormValue("name"))
		if r.MultipartForm != nil && len(r.MultipartForm.File["avatar"]) > 0 {
			avatar = r.MultipartForm.File["avatar"][0]
		}
	} else {
		_ = r.ParseForm()
		name = strings.TrimSpace(r.FormValue("name"))
	}
	a, err := h.svc.Update(r.Context(), uint(id), name, avatar)
	if err != nil {
		response.Internal(w, err.Error())
		return
	}
	if a == nil {
		response.NotFound(w, "author not found")
		return
	}
	response.OK(w, a)
}

func (h *AuthorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), uint(id)); err != nil {
		response.Internal(w, "failed to delete author")
		return
	}
	response.NoContent(w)
}
