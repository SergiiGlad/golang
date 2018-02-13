package controllers

import (
  "go-team-room/models/dto"
  "github.com/pkg/errors"
  "testing"
  "github.com/stretchr/testify/assert"
)

type EmailSendMock struct {
}

func (es EmailSendMock) SendEmail(email dto.Email) error {
  if email.Subject == "Bad Email" {
    return errors.New("Error")
  }
  return nil
}

type TokenGeneratorMock struct {
}

func (tg TokenGeneratorMock) ApproveUser(token string) (bool, error) {
  if token == "badToken" {
    return false, errors.New("Bad token")
  }
  if token == "usedToken" {
    return false, nil
  }
  return true, nil
}

func (tg TokenGeneratorMock) GenerateTokenForEmail(email string) (string, error) {
  if email == "email@email.com" {
    return "newToken", nil
  }
  return "", errors.New("Error")
}

var es = UserEmailService{
  HermesEmailBodyGenerator{},
  EmailSendMock{},
  TokenGeneratorMock{},
}

func TestUserEmailService_SendEmails(t *testing.T) {
  email := dto.Email{Subject: "Whatever", To: "email@email.com"}
  assert.Nil(t, es.SendEmails(email))
}

func TestUserEmailService_SendChangePasswordConfirmationEmail(t *testing.T) {
  user := dto.RequestUserDto{Email: "email@email.com"}
  assert.Nil(t, es.SendChangePasswordConfirmationEmail(user, "password123"))
}

func TestUserEmailService_SendRegistrationConfirmationEmail(t *testing.T) {
  user := dto.RequestUserDto{Email: "email@email.com"}
  assert.Nil(t, es.SendRegistrationConfirmationEmail(user))
}

func TestUserEmailService_SendWelcomeEmail(t *testing.T) {
  user := dto.RequestUserDto{Email: "email@email.com"}
  assert.Nil(t, es.SendWelcomeEmail(user))
}
