// handler/post_handler: HTTP handlers for post CRUD and file upload.
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

type PostHandler struct {
	svc *service.PostService
}

func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

// Create godoc
//
//	@Summary		Create a post
//	@Description	Creates a new post with title, body, and optional image/video files.
//	@Tags			posts
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			title	formData	string	true	"Post title"
//	@Param			body	formData	string	false	"Post body"
//	@Param			files	formData	file	false	"Image or video files"
//	@Success		201		{object}	response.Body{data=model.Post}
//	@Failure		400		{object}	response.Body
//	@Router			/posts [post]
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.BadRequest(w, "invalid request format")
		return
	}
	title := r.FormValue("title")
	body := r.FormValue("body")
	authorID := parseOptionalUint(r.FormValue("author_id"))
	categoryID := parseOptionalUint(r.FormValue("category_id"))
	var banner *multipart.FileHeader
	if r.MultipartForm != nil && len(r.MultipartForm.File["banner"]) > 0 {
		banner = r.MultipartForm.File["banner"][0]
	}
	files := r.MultipartForm.File["files"]
	post, err := h.svc.Create(r.Context(), title, body, authorID, categoryID, banner, files)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, post)
}

// GetByID godoc
//
//	@Summary		Get a post by ID
//	@Description	Returns the post with the given ID.
//	@Tags			posts
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	response.Body{data=model.Post}
//	@Failure		400	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Router			/posts/{id} [get]
func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	post, err := h.svc.GetByID(r.Context(), uint(id))
	if err != nil {
		response.NotFound(w, "post not found")
		return
	}
	if post == nil {
		response.NotFound(w, "post not found")
		return
	}
	response.OK(w, post)
}

// List godoc
//
//	@Summary		List posts
//	@Description	Returns a paginated list of posts (limit and offset).
//	@Tags			posts
//	@Produce		json
//	@Param			limit	query		int	false	"Items per page (default 20)"
//	@Param			offset	query		int	false	"Number of items to skip"
//	@Success		200		{object}	response.Body{data=[]model.Post}
//	@Failure		500		{object}	response.Body
//	@Router			/posts [get]
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	categoryID := parseOptionalUint(r.URL.Query().Get("category_id"))
	posts, err := h.svc.List(r.Context(), limit, offset, categoryID)
	if err != nil {
		response.Internal(w, "failed to list posts")
		return
	}
	response.OK(w, posts)
}

// Update godoc
//
//	@Summary		Update a post
//	@Description	Updates a post. Empty fields are left unchanged. New files are optional.
//	@Tags			posts
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		int		true	"Post ID"
//	@Param			title	formData	string	false	"New title"
//	@Param			body	formData	string	false	"New body"
//	@Param			files	formData	file	false	"New files"
//	@Success		200		{object}	response.Body{data=model.Post}
//	@Failure		400		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Router			/posts/{id} [put]
func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	var title, body string
	var authorID, categoryID *uint
	var banner *multipart.FileHeader
	var files []*multipart.FileHeader
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			response.BadRequest(w, "invalid request format")
			return
		}
		title = r.FormValue("title")
		body = r.FormValue("body")
		authorID = parseOptionalUint(r.FormValue("author_id"))
		categoryID = parseOptionalUint(r.FormValue("category_id"))
		if r.MultipartForm != nil {
			if len(r.MultipartForm.File["banner"]) > 0 {
				banner = r.MultipartForm.File["banner"][0]
			}
			files = r.MultipartForm.File["files"]
		}
	} else {
		_ = r.ParseForm()
		title = r.FormValue("title")
		body = r.FormValue("body")
		authorID = parseOptionalUint(r.FormValue("author_id"))
		categoryID = parseOptionalUint(r.FormValue("category_id"))
	}
	post, err := h.svc.Update(r.Context(), uint(id), title, body, authorID, categoryID, banner, files)
	if err != nil {
		response.Internal(w, err.Error())
		return
	}
	if post == nil {
		response.NotFound(w, "post not found")
		return
	}
	response.OK(w, post)
}

// Delete godoc
//
//	@Summary		Delete a post
//	@Description	Deletes the post with the given ID (soft delete).
//	@Tags			posts
//	@Param			id	path	int	true	"Post ID"
//	@Success		204	"No content"
//	@Failure		400	{object}	response.Body
//	@Failure		500	{object}	response.Body
//	@Router			/posts/{id} [delete]
func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), uint(id)); err != nil {
		response.Internal(w, "failed to delete post")
		return
	}
	response.NoContent(w)
}

func parseOptionalUint(s string) *uint {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return nil
	}
	u := uint(n)
	return &u
}
