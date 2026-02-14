// handler/comment_handler: HTTP handlers for comments (under a post).
package handler

import (
	"net/http"
	"strconv"

	"github.com/aliakbar-zohour/go_blog/internal/service"
	"github.com/aliakbar-zohour/go_blog/pkg/response"
	"github.com/go-chi/chi/v5"
)

type CommentHandler struct {
	svc *service.CommentService
}

func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

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

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseUint(chi.URLParam(r, "postId"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid post id")
		return
	}
	_ = r.ParseForm()
	body := r.FormValue("body")
	authorName := r.FormValue("author_name")
	c, err := h.svc.Create(r.Context(), uint(postID), body, authorName)
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

func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	_ = r.ParseForm()
	body := r.FormValue("body")
	c, err := h.svc.Update(r.Context(), uint(id), body)
	if err != nil {
		response.Internal(w, err.Error())
		return
	}
	if c == nil {
		response.NotFound(w, "comment not found")
		return
	}
	response.OK(w, c)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), uint(id)); err != nil {
		response.Internal(w, "failed to delete comment")
		return
	}
	response.NoContent(w)
}
