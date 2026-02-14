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

// Create godoc
//	@Summary		ساخت پست
//	@Description	پست جدید با عنوان و متن و اختیاری فایل (عکس/ویدئو) ایجاد می‌کند.
//	@Tags			posts
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			title	formData	string	true	"عنوان پست"
//	@Param			body	formData	string	false	"متن پست"
//	@Param			files	formData	file	false	"فایل‌های تصویر یا ویدئو"
//	@Success		201		{object}	response.Body{data=model.Post}
//	@Failure		400		{object}	response.Body
//	@Router			/posts [post]
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

// GetByID godoc
//	@Summary		دریافت یک پست
//	@Description	پست با شناسه مشخص را برمی‌گرداند.
//	@Tags			posts
//	@Produce		json
//	@Param			id	path		int	true	"شناسه پست"
//	@Success		200	{object}	response.Body{data=model.Post}
//	@Failure		400	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Router			/posts/{id} [get]
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

// List godoc
//	@Summary		لیست پست‌ها
//	@Description	لیست پست‌ها با صفحه‌بندی (limit و offset).
//	@Tags			posts
//	@Produce		json
//	@Param			limit	query		int	false	"تعداد در هر صفحه (پیش‌فرض 20)"
//	@Param			offset	query		int	false	"تعداد رد شده از ابتدا"
//	@Success		200		{object}	response.Body{data=[]model.Post}
//	@Failure		500		{object}	response.Body
//	@Router			/posts [get]
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

// Update godoc
//	@Summary		ویرایش پست
//	@Description	پست را به‌روزرسانی می‌کند؛ فیلدهای خالی تغییر نمی‌کنند. فایل‌های جدید اختیاری.
//	@Tags			posts
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		int		true	"شناسه پست"
//	@Param			title	formData	string	false	"عنوان جدید"
//	@Param			body	formData	string	false	"متن جدید"
//	@Param			files	formData	file	false	"فایل‌های جدید"
//	@Success		200		{object}	response.Body{data=model.Post}
//	@Failure		400		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Router			/posts/{id} [put]
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

// Delete godoc
//	@Summary		حذف پست
//	@Description	پست با شناسه مشخص حذف می‌شود (soft delete).
//	@Tags			posts
//	@Param			id	path	int	true	"شناسه پست"
//	@Success		204	"No content"
//	@Failure		400	{object}	response.Body
//	@Failure		500	{object}	response.Body
//	@Router			/posts/{id} [delete]
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
