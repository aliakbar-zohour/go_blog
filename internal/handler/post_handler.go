// handler/post_handler: هندلرهای HTTP برای CRUD پست و آپلود فایل.
package handler

import (
	"net/http"
	"strconv"

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

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.BadRequest(w, "فرمت درخواست نامعتبر")
		return
	}
	title := r.FormValue("title")
	body := r.FormValue("body")
	files := r.MultipartForm.File["files"]
	post, err := h.svc.Create(r.Context(), title, body, files)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, post)
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "شناسه نامعتبر")
		return
	}
	post, err := h.svc.GetByID(r.Context(), uint(id))
	if err != nil {
		response.NotFound(w, "پست یافت نشد")
		return
	}
	if post == nil {
		response.NotFound(w, "پست یافت نشد")
		return
	}
	response.OK(w, post)
}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	posts, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		response.Internal(w, "خطا در دریافت لیست")
		return
	}
	response.OK(w, posts)
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "شناسه نامعتبر")
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.BadRequest(w, "فرمت درخواست نامعتبر")
		return
	}
	title := r.FormValue("title")
	body := r.FormValue("body")
	files := r.MultipartForm.File["files"]
	post, err := h.svc.Update(r.Context(), uint(id), title, body, files)
	if err != nil {
		response.Internal(w, err.Error())
		return
	}
	if post == nil {
		response.NotFound(w, "پست یافت نشد")
		return
	}
	response.OK(w, post)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "شناسه نامعتبر")
		return
	}
	if err := h.svc.Delete(r.Context(), uint(id)); err != nil {
		response.Internal(w, "خطا در حذف")
		return
	}
	response.NoContent(w)
}
