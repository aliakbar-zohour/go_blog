# Go Blog API

بلاگ REST API با Go، دیتابیس PostgreSQL و پشتیبانی از آپلود تصویر و ویدئو.

## ساختار پروژه

- `cmd/api` — نقطه ورود برنامه
- `internal/config` — تنظیمات از env
- `internal/database` — اتصال و مایگریشن
- `internal/model` — مدل‌های پست و مدیا
- `internal/repository` — لایه دسترسی به داده
- `internal/service` — منطق کسب‌وکار
- `internal/handler` — هندلرهای HTTP
- `internal/router` — مسیرها و میدلور
- `internal/middleware` — بازیابی پنیک، هدر امن، لاگ
- `internal/upload` — اعتبارسنجی و ذخیره فایل
- `pkg/response` — پاسخ JSON یکسان

## اجرا با Docker

```bash
docker compose up --build
```

API روی `http://localhost:8080` در دسترس است.

## اجرای محلی

1. PostgreSQL را اجرا کنید.
2. کپی فایل env:
   ```bash
   cp .env.example .env
   ```
3. وابستگی و اجرا:
   ```bash
   go mod tidy
   go run ./cmd/api
   ```

## endpoints

| متد   | مسیر           | توضیح                    |
|-------|----------------|---------------------------|
| GET   | /api/posts     | لیست پست‌ها (limit, offset) |
| POST  | /api/posts     | ساخت پست (form: title, body, files[]) |
| GET   | /api/posts/:id | دریافت یک پست            |
| PUT   | /api/posts/:id | ویرایش پست (form: title, body, files[]) |
| DELETE| /api/posts/:id | حذف پست                  |

فایل‌های آپلودشده از طریق `/uploads/<path>` سرو می‌شوند.  
فرمت‌های مجاز تصویر: jpg, jpeg, png, gif, webp. ویدئو: mp4, webm, mov.

## لایسنس

MIT
