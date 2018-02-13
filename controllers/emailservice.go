package controllers

import (
  "go-team-room/models/dto"
  "errors"
  "fmt"
)

const (
  passwordSubject     = "Please confrim password change!"
  registrationSubject = "Please confirm your email!"
  welcomeSubject      = "Welcome to GOHUM!"
  invalidEmail        = "Email %s is invalid."
  serverNoResponse    = "No response from mail server"
)

// Default implementation of EmailServiceInterface uses body generator and email sender for sending emails to users.
type EmailService struct {
  BG EmailBodyGeneratorInterface
  ES EmailSendInterface
  TG TokenGeneratorInterface
}

var _ EmailServiceInterface = &EmailService{}

// Send all emails if at least one fails return error with explanation.
func (ess *EmailService) SendEmails(emails ...dto.Email) error {
  log.Debug("Start sending emails: ", emails)
  for _, e := range emails {
    if !ValidEmail(e.To) {
      log.Error("Email address of email is invalid, email: {}, address: {}", e, e.To)
      return errors.New(fmt.Sprintf(invalidEmail, e.To))
    }
    err := ess.ES.SendEmail(e)
    if err != nil {
      log.Error("Fail to send email: {}, error: {}", e, err)
      return errors.New(serverNoResponse)
    }
  }
  log.Info("Successfully send emails: {}", emails)
  return nil
}

// Send email with welcome text for user with 'CONFIRMED' email.
func (ess *EmailService) SendWelcomeEmail(user dto.RequestUserDto) error {
  log.Debug("Sending new Welcome email for user: {}, subject: {}", user, welcomeSubject)
  body := ess.BG.GenerateWelcomeBody(user)
  email := dto.RequestUserDtoToEmail(user, welcomeSubject, body)
  return ess.SendEmails(email)
}

// Send email with request for email confirmation to User with unconfirmed email.
func (ess *EmailService) SendRegistrationConfirmationEmail(user dto.RequestUserDto) error {
  //log.Debug("Sending new Registration Confirmation email for user: {}, subject: {}", user, registrationSubject)
  token, err := ess.TG.GenerateTokenForEmail(user.Email)
  if err != nil {log.Error("Fail to send registration confirmation email for user: {}, err: {}",user, err)
    return err
  }
  body := ess.BG.GenerateRegistrationConfirmationEmail(user, token)
  email := dto.RequestUserDtoToEmail(user, registrationSubject, body)
  return ess.SendEmails(email)
}

// Send email with confirmation for password change.
func (ess *EmailService) SendChangePasswordConfirmationEmail(user dto.RequestUserDto, newPassword string) error {
  log.Debug("Sending new Password Confirmation email for user: {}, subject: {}", user, passwordSubject)
  body := ess.BG.GenerateChangePasswordConfirmationEmail(user, newPassword)
  email := dto.RequestUserDtoToEmail(user, passwordSubject, body)
  return ess.SendEmails(email)
}
