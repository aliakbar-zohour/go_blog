// handler/author_handler: HTTP handlers for author CRUD.
package handler

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/aliakbar-zohour/go_blog/internal/middleware"
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

// List godoc
//
//	@Summary		List authors
//	@Description	Returns all authors.
//	@Tags			authors
//	@Produce		json
//	@Success		200	{object}	response.Body{data=[]model.Author}
//	@Failure		500	{object}	response.Body
//	@Router			/authors [get]
func (h *AuthorHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.List(r.Context())
	if err != nil {
		response.Internal(w, "failed to list authors")
		return
	}
	response.OK(w, list)
}

// Create godoc
//
//	@Summary		Create an author
//	@Description	Creates a new author with name and optional avatar image.
//	@Tags			authors
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name	formData	string	true	"Author name"
//	@Param			avatar	formData	file	false	"Avatar image"
//	@Success		201		{object}	response.Body{data=model.Author}
//	@Failure		400		{object}	response.Body
//	@Router			/authors [post]
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

// GetByID godoc
//
//	@Summary		Get an author by ID
//	@Description	Returns the author with the given ID.
//	@Tags			authors
//	@Produce		json
//	@Param			id	path		int	true	"Author ID"
//	@Success		200	{object}	response.Body{data=model.Author}
//	@Failure		400	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Router			/authors/{id} [get]
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

// Update godoc
//
//	@Summary		Update an author
//	@Description	Updates own profile (name and/or avatar). Requires Authorization: Bearer <token>. You can only update yourself.
//	@Tags			authors
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		Bearer
//	@Param			id		path		int		true	"Author ID"
//	@Param			name	formData	string	false	"Author name"
//	@Param			avatar	formData	file	false	"Avatar image"
//	@Success		200		{object}	response.Body{data=model.Author}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Failure		403		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Router			/authors/{id} [put]
func (h *AuthorHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	loggedID := middleware.GetAuthorID(r.Context())
	if loggedID == 0 || loggedID != uint(id) {
		response.Forbidden(w, "you can only update your own profile")
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

// Delete godoc
//
//	@Summary		Delete an author
//	@Description	Deletes own account. Requires Authorization: Bearer <token>. You can only delete yourself.
//	@Tags			authors
//	@Security		Bearer
//	@Param			id	path	int	true	"Author ID"
//	@Success		204	"No content"
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		403	{object}	response.Body
//	@Failure		500	{object}	response.Body
//	@Router			/authors/{id} [delete]
func (h *AuthorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	loggedID := middleware.GetAuthorID(r.Context())
	if loggedID == 0 || loggedID != uint(id) {
		response.Forbidden(w, "you can only delete your own account")
		return
	}
	if err := h.svc.Delete(r.Context(), uint(id)); err != nil {
		response.Internal(w, "failed to delete author")
		return
	}
	response.NoContent(w)
}
