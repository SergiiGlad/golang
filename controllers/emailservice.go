package controllers

import (
  "go-team-room/models/dto"
  "go-team-room/conf"
)

const (
  passwordSubject     = "Please confrim password change!"
  registrationSubject = "Please confirm your email!"
  welcomeSubject      = "Welcome to GOHUM!"
  invalidEmail        = "Email %s is invalid."
  serverNoResponse    = "No response from mail server"
)

// Default implementation of UserEmailServiceInterface uses body generator and email sender for sending emails to users.
type UserEmailService struct {
  BG EmailBodyGeneratorInterface
  ES EmailSendInterface
  TG TokenGeneratorInterface
}

var _ UserEmailServiceInterface = &UserEmailService{}

func (ess *UserEmailService) SendEmails(emails ...dto.Email) {
  log.Debugf("Start sending emails: %s", emails)
  for _, e := range emails {
    if !ValidEmail(e.To) {
      log.Errorf("Email address of email is invalid, email: %s, address: %s", e, e.To)
    }
    if conf.EnableSendMails {
      err := ess.ES.SendEmail(e)
      if err != nil {
        log.Errorf("Fail to send email: %s, error: %s", e, err)
      }
    }
  }
  log.Infof("Successfully send emails: %s", emails)
}

// Send email with welcome text for user with 'CONFIRMED' email.
func (ess *UserEmailService) SendWelcomeEmail(user dto.RequestUserDto) error {
  log.Debugf("Sending new Welcome email for user: %s, subject: %s", user, welcomeSubject)
  body := ess.BG.GenerateWelcomeBody(user)
  email := dto.RequestUserDtoToEmail(user, welcomeSubject, body)
  go ess.SendEmails(email)
  return nil
}

// Send email with request for email confirmation to User with unconfirmed email.
func (ess *UserEmailService) SendRegistrationConfirmationEmail(user dto.RequestUserDto) error {
  log.Debugf("Sending new Registration Confirmation email for user: %s, subject: %s", user, registrationSubject)
  token, err := ess.TG.GenerateTokenForEmail(user.Email)
  if err != nil {
    log.Errorf("Fail to send registration confirmation email for user: %s, err: %s", user, err)
    return err
  }
  body := ess.BG.GenerateRegistrationConfirmationEmail(user, token)
  email := dto.RequestUserDtoToEmail(user, registrationSubject, body)
  go ess.SendEmails(email)
  return nil
}

// Send email with confirmation for password change.
func (ess *UserEmailService) SendChangePasswordConfirmationEmail(user dto.RequestUserDto, newPassword string) error {
  log.Debugf("Sending new Password Confirmation email for user: %s, subject: %s", user, passwordSubject)
  body := ess.BG.GenerateChangePasswordConfirmationEmail(user, newPassword)
  email := dto.RequestUserDtoToEmail(user, passwordSubject, body)
  go ess.SendEmails(email)
  return nil
}
