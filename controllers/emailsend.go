package controllers

import (
  "go-team-room/models/dto"
  "github.com/jordan-wright/email"
  "net/smtp"
)

// Default implementation of SendEmailInterface uses "github.com/jordan-wright/email" package.
type DefaultEmailSend struct{}

var _ EmailSendInterface = &DefaultEmailSend{}

// Send email
func (des *DefaultEmailSend) SendEmail(emailDto dto.Email) error {
  e := email.NewEmail()
  e.To = []string{emailDto.To}
  e.Subject = emailDto.Subject
  e.HTML = emailDto.Body
  return e.Send("addr", smtp.PlainAuth("", "from", "password", "host"))
}
