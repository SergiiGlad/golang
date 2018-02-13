package controllers

import (
  "go-team-room/models/dto"
  "gopkg.in/gomail.v2"
  "go-team-room/conf"
)

var d = gomail.NewPlainDialer(conf.SmtpServer, conf.SmtpPort, conf.GohumEmail, conf.GohumEmailPass)

// Default implementation of SendEmailInterface uses "github.com/jordan-wright/email" package.
type DefaultEmailSend struct{}

var _ EmailSendInterface = &DefaultEmailSend{}

// Send email using conf.GohumEmail.
func (des *DefaultEmailSend) SendEmail(emailDto dto.Email) error {
  log.Debug("Sending new email: {}", emailDto)
  m := gomail.NewMessage()
  m.SetHeader("From", conf.GohumEmail)
  m.SetHeader("To", emailDto.To)
  m.SetHeader("Subject", emailDto.Subject)
  m.SetBody("text/html", emailDto.Body)
  log.Debug("Email to send: {}", m)
  err := d.DialAndSend(m)
  d.Dial()
  if err != nil {
    log.Error("Faild to send email: {}, formed email was: {}, error: ", emailDto, m, err)
    return err
  }
  log.Debug("Successfully send email: {}", emailDto)
  return nil
}
