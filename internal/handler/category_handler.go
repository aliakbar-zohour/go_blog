// handler/category_handler: HTTP handlers for category CRUD.
package handler

import (
	"net/http"
	"strconv"

	"github.com/aliakbar-zohour/go_blog/internal/service"
	"github.com/aliakbar-zohour/go_blog/pkg/response"
	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.List(r.Context())
	if err != nil {
		response.Internal(w, "failed to list categories")
		return
	}
	response.OK(w, list)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	name := r.FormValue("name")
	c, err := h.svc.Create(r.Context(), name)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, c)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	c, err := h.svc.GetByID(r.Context(), uint(id))
	if err != nil {
		response.NotFound(w, "category not found")
		return
	}
	if c == nil {
		response.NotFound(w, "category not found")
		return
	}
	response.OK(w, c)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	_ = r.ParseForm()
	name := r.FormValue("name")
	c, err := h.svc.Update(r.Context(), uint(id), name)
	if err != nil {
		response.Internal(w, err.Error())
		return
	}
	if c == nil {
		response.NotFound(w, "category not found")
		return
	}
	response.OK(w, c)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), uint(id)); err != nil {
		response.Internal(w, "failed to delete category")
		return
	}
	response.NoContent(w)
}
