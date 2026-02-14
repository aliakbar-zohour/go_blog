// handler/comment_handler: HTTP handlers for comments (under a post).
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/aliakbar-zohour/go_blog/internal/middleware"
	"github.com/aliakbar-zohour/go_blog/internal/service"
	"github.com/aliakbar-zohour/go_blog/pkg/response"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type CommentHandler struct {
	svc *service.CommentService
}

func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

// ListByPostID godoc
//
//	@Summary		List comments for a post
//	@Description	Returns all comments for the given post ID.
//	@Tags			comments
//	@Produce		json
//	@Param			postId	path		int	true	"Post ID"
//	@Success		200		{object}	response.Body{data=[]model.Comment}
//	@Failure		400		{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Router			/posts/{postId}/comments [get]
func (h *CommentHandler) ListByPostID(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseUint(chi.URLParam(r, "postId"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid post id")
		return
	}
	list, err := h.svc.ListByPostID(r.Context(), uint(postID))
	if err != nil {
		response.Internal(w, "failed to list comments")
		return
	}
	response.OK(w, list)
}

// Create godoc
//
//	@Summary		Create a comment
//	@Description	Creates a new comment on the given post. Requires Authorization: Bearer <token>. Comment is linked to the logged-in author.
//	@Tags			comments
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Security		Bearer
//	@Param			postId		path		int		true	"Post ID"
//	@Param			body		formData	string	true	"Comment body"
//	@Param			author_name	formData	string	false	"Optional display name override"
//	@Success		201			{object}	response.Body{data=model.Comment}
//	@Failure		400			{object}	response.Body
//	@Failure		401			{object}	response.Body
//	@Failure		404			{object}	response.Body
//	@Router			/posts/{postId}/comments [post]
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseUint(chi.URLParam(r, "postId"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid post id")
		return
	}
	authorID := middleware.GetAuthorID(r.Context())
	if authorID == 0 {
		response.Unauthorized(w, "authorization required to create a comment")
		return
	}
	_ = r.ParseForm()
	body := r.FormValue("body")
	authorName := r.FormValue("author_name")
	c, err := h.svc.Create(r.Context(), uint(postID), body, authorName, &authorID)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	if c == nil {
		response.NotFound(w, "post not found")
		return
	}
	response.Created(w, c)
}

// Update godoc
//
//	@Summary		Update a comment
//	@Description	Updates own comment. Requires Authorization: Bearer <token>. You can only edit your own comment.
//	@Tags			comments
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Security		Bearer
//	@Param			id		path		int		true	"Comment ID"
//	@Param			body	formData	string	false	"Comment body"
//	@Success		200		{object}	response.Body{data=model.Comment}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Failure		403		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Router			/comments/{id} [put]
func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	authorID := middleware.GetAuthorID(r.Context())
	if authorID == 0 {
		response.Unauthorized(w, "authorization required to update a comment")
		return
	}
	_ = r.ParseForm()
	body := r.FormValue("body")
	c, err := h.svc.Update(r.Context(), uint(id), body, authorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(w, "comment not found")
			return
		}
		if errors.Is(err, service.ErrCommentForbidden) {
			response.Forbidden(w, err.Error())
			return
		}
		response.Internal(w, "failed to update comment")
		return
	}
	if c == nil {
		response.NotFound(w, "comment not found")
		return
	}
	response.OK(w, c)
}

// Delete godoc
//
//	@Summary		Delete a comment
//	@Description	Deletes own comment. Requires Authorization: Bearer <token>. You can only delete your own comment.
//	@Tags			comments
//	@Security		Bearer
//	@Param			id	path	int	true	"Comment ID"
//	@Success		204	"No content"
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		403	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Failure		500	{object}	response.Body
//	@Router			/comments/{id} [delete]
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	authorID := middleware.GetAuthorID(r.Context())
	if authorID == 0 {
		response.Unauthorized(w, "authorization required to delete a comment")
		return
	}
	if err := h.svc.Delete(r.Context(), uint(id), authorID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(w, "comment not found")
			return
		}
		if errors.Is(err, service.ErrDeleteCommentForbidden) {
			response.Forbidden(w, err.Error())
			return
		}
		response.Internal(w, "failed to delete comment")
		return
	}
	response.NoContent(w)
}
