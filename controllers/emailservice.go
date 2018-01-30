package controllers

import (
  "go-team-room/models/dto"
  "errors"
  "fmt"
)

const (
  passwordSubject = "Please confrim password change!"
  registrationSubject = "Please confirm your email!"
  welcomeSubject   = "Welcome to GOHUM!"
  invalidEmail     = "Email %s is invalid."
  serverNoResponse = "No response from mail server"
)

// Default implementation of EmailServiceInterface uses body generator and email sender for sending emails to users.
type EmailService struct{
  BG EmailBodyGeneratorInterface
  ES EmailSendInterface
  TG TokenGeneratorInterface
}

var _ EmailServiceInterface = &EmailService{}

// Send all emails if at least one fails return error with explanation.
func (ess *EmailService) SendEmails(emails ...dto.Email) error {
  for _, e := range emails {
    if !ValidEmail(e.To) {
      return errors.New(fmt.Sprintf(invalidEmail, e.To))
    }
    err := ess.ES.SendEmail(e)
    if err != nil {
      return errors.New(serverNoResponse)
    }
  }
  return nil
}

// Send email with welcome text for user with 'CONFIRMED' email.
func (ess *EmailService) SendWelcomeEmail(user dto.RequestUserDto) error {
  body := ess.BG.GenerateWelcomeBody(user)
  email := dto.RequestUserDtoToEmail(user, welcomeSubject, body)
  return ess.SendEmails(email)
}

// Send email with request for email confirmation to User with unconfirmed email.
func (ess *EmailService) SendRegistrationConfirmationEmail(user dto.RequestUserDto) error {
  body := ess.BG.GenerateRegistrationConfirmationEmail(user, "")
  email := dto.RequestUserDtoToEmail(user, registrationSubject, body)
  return ess.SendEmails(email)
}

// Send email with confirmation for password change.
func (ess *EmailService) SendChangePasswordConfirmationEmail(user dto.RequestUserDto) error {
  body := ess.BG.GenerateChangePasswordConfirmationEmail(user, "")
  email := dto.RequestUserDtoToEmail(user, passwordSubject, body)
  return ess.SendEmails(email)
}
