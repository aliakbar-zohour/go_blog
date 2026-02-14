// mail: Sends verification code email with an HTML template.
package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
)

// SendVerificationCode sends a beautiful HTML email with the code. If SMTP is not configured, returns nil (no error) and no email is sent.
func SendVerificationCode(toEmail, code, smtpHost, smtpPort, smtpUser, smtpPass, from string) error {
	if smtpHost == "" {
		return nil
	}
	subject := "Your verification code – Go Blog"
	htmlBody := buildVerificationHTML(code)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", from, toEmail, subject, htmlBody)
	addr := smtpHost + ":" + smtpPort
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	return smtp.SendMail(addr, auth, from, []string{toEmail}, []byte(msg))
}

func buildVerificationHTML(code string) string {
	var b bytes.Buffer
	b.WriteString(`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"></head><body style="margin:0;font-family:'Segoe UI',system-ui,sans-serif;background:linear-gradient(135deg,#1a1a2e 0%,#16213e 50%,#0f3460 100%);min-height:100vh;display:flex;align-items:center;justify-content:center;padding:20px;box-sizing:border-box">`)
	b.WriteString(`<div style="background:rgba(255,255,255,0.08);backdrop-filter:blur(12px);border:1px solid rgba(255,255,255,0.12);border-radius:20px;padding:48px 40px;max-width:420px;width:100%;text-align:center;box-shadow:0 25px 50px -12px rgba(0,0,0,0.4)">`)
	b.WriteString(`<div style="font-size:28px;font-weight:700;color:#e94560;margin-bottom:8px;letter-spacing:-0.5px">Go Blog</div>`)
	b.WriteString(`<div style="color:rgba(255,255,255,0.7);font-size:14px;margin-bottom:32px">Writer verification</div>`)
	b.WriteString(`<p style="color:rgba(255,255,255,0.9);font-size:15px;line-height:1.6;margin:0 0 24px">Use this code to complete your registration:</p>`)
	b.WriteString(`<div style="background:rgba(233,69,96,0.2);border:2px solid #e94560;border-radius:12px;padding:20px 28px;margin:0 0 32px">`)
	b.WriteString(`<span style="font-size:32px;font-weight:700;letter-spacing:8px;color:#fff">`)
	b.WriteString(escapeHTML(code))
	b.WriteString(`</span></div>`)
	b.WriteString(`<p style="color:rgba(255,255,255,0.5);font-size:12px;margin:0">This code expires in 15 minutes. If you didn't request it, ignore this email.</p>`)
	b.WriteString(`</div></body></html>`)
	return b.String()
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
