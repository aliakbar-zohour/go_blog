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

// List godoc
//
//	@Summary		List categories
//	@Description	Returns all categories.
//	@Tags			categories
//	@Produce		json
//	@Success		200	{object}	response.Body{data=[]model.Category}
//	@Failure		500	{object}	response.Body
//	@Router			/categories [get]
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.List(r.Context())
	if err != nil {
		response.Internal(w, "failed to list categories")
		return
	}
	response.OK(w, list)
}

// Create godoc
//
//	@Summary		Create a category
//	@Description	Creates a new category with the given name.
//	@Tags			categories
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			name	formData	string	true	"Category name"
//	@Success		201		{object}	response.Body{data=model.Category}
//	@Failure		400		{object}	response.Body
//	@Router			/categories [post]
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

// GetByID godoc
//
//	@Summary		Get a category by ID
//	@Description	Returns the category with the given ID.
//	@Tags			categories
//	@Produce		json
//	@Param			id	path		int	true	"Category ID"
//	@Success		200	{object}	response.Body{data=model.Category}
//	@Failure		400	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Router			/categories/{id} [get]
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

// Update godoc
//
//	@Summary		Update a category
//	@Description	Updates the category name.
//	@Tags			categories
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			id		path		int		true	"Category ID"
//	@Param			name	formData	string	false	"Category name"
//	@Success		200		{object}	response.Body{data=model.Category}
//	@Failure		400		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Router			/categories/{id} [put]
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

// Delete godoc
//
//	@Summary		Delete a category
//	@Description	Deletes the category with the given ID.
//	@Tags			categories
//	@Param			id	path	int	true	"Category ID"
//	@Success		204	"No content"
//	@Failure		400	{object}	response.Body
//	@Failure		500	{object}	response.Body
//	@Router			/categories/{id} [delete]
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
