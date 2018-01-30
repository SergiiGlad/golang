package controllers

import "go-team-room/models/dto"

type UserServiceInterface interface {
  CreateUser(userDto *dto.RequestUserDto) (dto.ResponseUserDto, error)
  UpdateUser(id int64, userDto *dto.RequestUserDto) (dto.ResponseUserDto, error)
  DeleteUser(id int64) (dto.ResponseUserDto, error)
  GetUserFriends(id int64) ([]int64, error)
}

type EmailServiceInterface interface {
  // Send all emails if at least one fails return error with explanation.
  SendEmails(email ...dto.Email) error
  // Send email with welcome text for user with 'CONFIRMED' email.
  SendWelcomeEmail(user dto.RequestUserDto) error
  // Send email with request for email confirmation to User with unconfirmed email.
  SendRegistrationConfirmationEmail(user dto.RequestUserDto) error
  // Send email with confirmation for password change.
  SendChangePasswordConfirmationEmail(user dto.RequestUserDto) error
}

type EmailBodyGeneratorInterface interface {
  // Generate message body for welcome email.
  GenerateWelcomeBody(user dto.RequestUserDto) []byte
  // Generate message body for registration confirmation email where token is confirmation token.
  GenerateRegistrationConfirmationEmail(user dto.RequestUserDto, token string) []byte
  // Generate message body for password change confirmation where token is confirmation token.
  GenerateChangePasswordConfirmationEmail(user dto.RequestUserDto, token string) []byte
}

type TokenGeneratorInterface interface {
  GenerateToken(email string) (error, string) //TODO implement this interface
  }

type EmailSendInterface interface {
  // Send email
  SendEmail(email dto.Email) error
}
