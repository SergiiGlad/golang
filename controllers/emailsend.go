package controllers

import (
  "go-team-room/models/dto"
  "gopkg.in/gomail.v2"
)

var d = gomail.NewPlainDialer("smtp.gmail.com", 587, "gohum.assistant@gmail.com", "gohum123")


// Default implementation of SendEmailInterface uses "github.com/jordan-wright/email" package.
type DefaultEmailSend struct{}

var _ EmailSendInterface = &DefaultEmailSend{}

// Send email
func (des *DefaultEmailSend) SendEmail(emailDto dto.Email) error {
  m := gomail.NewMessage()
  m.SetHeader("From", "gohum.assistant@gmail.com")
  m.SetHeader("To", emailDto.To)
  m.SetHeader("Subject", emailDto.Subject)
  m.SetBody("text/html", emailDto.Body)
  err := d.DialAndSend(m)
  d.Dial()
  return err
}
