// handler/auth_handler: Registration (request code, verify & register) and login; returns JWT.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aliakbar-zohour/go_blog/internal/service"
	"github.com/aliakbar-zohour/go_blog/pkg/response"
)

// AuthRegisterRequest body for POST /auth/register/request
type AuthRegisterRequest struct {
	Email string `json:"email" example:"writer@example.com"`
}

// AuthRegisterVerifyRequest body for POST /auth/register/verify
type AuthRegisterVerifyRequest struct {
	Email    string `json:"email" example:"writer@example.com"`
	Code     string `json:"code" example:"123456"`
	Name     string `json:"name" example:"Jane Doe"`
	Password string `json:"password" example:"secret123"`
}

// AuthLoginRequest body for POST /auth/login
type AuthLoginRequest struct {
	Email    string `json:"email" example:"writer@example.com"`
	Password string `json:"password" example:"secret123"`
}

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// RequestVerification godoc
//
//	@Summary		Request verification code
//	@Description	Sends a 6-digit code to the given email (if SMTP configured). For testing without SMTP, the code is logged on the server.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		AuthRegisterRequest	true	"Email to send the code to"
//	@Success		200		{object}	response.Body{data=object}
//	@Failure		400		{object}	response.Body
//	@Router			/auth/register/request [post]
func (h *AuthHandler) RequestVerification(w http.ResponseWriter, r *http.Request) {
	var body AuthRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid body")
		return
	}
	if len(body.Email) > 255 {
		response.BadRequest(w, "email too long")
		return
	}
	devCode, err := h.svc.RequestVerification(r.Context(), body.Email)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	res := map[string]interface{}{"sent": true}
	if devCode != "" {
		res["dev_code"] = devCode
		res["message"] = "SMTP not configured; use dev_code in /auth/register/verify to complete registration."
	}
	response.OK(w, res)
}

// VerifyAndRegister godoc
//
//	@Summary		Verify code and complete registration
//	@Description	Verifies the code sent to email, creates the author account with name and password, returns author and JWT. Password must be at least 8 characters.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		AuthRegisterVerifyRequest	true	"Email, code, name and password"
//	@Success		201		{object}	response.Body	"data contains author and token"
//	@Failure		400		{object}	response.Body
//	@Router			/auth/register/verify [post]
func (h *AuthHandler) VerifyAndRegister(w http.ResponseWriter, r *http.Request) {
	var body AuthRegisterVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid body")
		return
	}
	if len(body.Email) > 255 || len(body.Name) > 255 || len(body.Code) > 10 || len(body.Password) > 128 {
		response.BadRequest(w, "field too long")
		return
	}
	a, token, err := h.svc.VerifyAndRegister(r.Context(), body.Email, body.Code, body.Name, body.Password)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	response.Created(w, map[string]interface{}{"author": a, "token": token})
}

// Login godoc
//
//	@Summary		Login
//	@Description	Returns author and JWT for valid email/password. Use the token in Authorization: Bearer <token> for protected routes.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		AuthLoginRequest	true	"Email and password"
//	@Success		200		{object}	response.Body	"data contains author and token"
//	@Failure		400		{object}	response.Body
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body AuthLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid body")
		return
	}
	if len(body.Email) > 255 || len(body.Password) > 128 {
		response.BadRequest(w, "field too long")
		return
	}
	a, token, err := h.svc.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	response.OK(w, map[string]interface{}{"author": a, "token": token})
}
