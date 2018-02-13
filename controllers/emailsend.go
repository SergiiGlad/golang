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
  log.Debugf("Sending new email: %s", emailDto)
  m := gomail.NewMessage()
  m.SetHeader("From", conf.GohumEmail)
  m.SetHeader("To", emailDto.To)
  m.SetHeader("Subject", emailDto.Subject)
  m.SetBody("text/html", emailDto.Body)
  log.Debugf("Email to send: %s", m)
  err := d.DialAndSend(m)
  d.Dial()
  if err != nil {
    log.Errorf("Faild to send email: %s, formed email was: %s, error: %s", emailDto, m, err)
    return err
  }
  log.Debugf("Successfully send email: %s", emailDto)
  return nil
}
