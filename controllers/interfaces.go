package controllers

import (
  "go-team-room/models/dto"
  "go-team-room/models/dao/entity"
)

//UserServiceInterface interface is used as HandlerFunc wrappers to follow dependency inversion principle
type UserServiceInterface interface {
  CreateUser(userDto *dto.RequestUserDto) (dto.ResponseUserDto, error)
  UpdateUser(id int64, userDto *dto.RequestUserDto) (dto.ResponseUserDto, error)
  DeleteUser(id int64) (dto.ResponseUserDto, error)
}

type FriendServiceInterface interface {
  GetFriends(id int64) ([]dto.ShortUser, error)
  GetUsersWithRequests(id int64) ([]dto.ShortUser, error)
  GetFriendIds(id int64) ([]int64, error)
  NewFriendRequest(connection *entity.Connection) error
  ApproveFriendRequest(connection *entity.Connection) error
  RejectFriendRequest(connection *entity.Connection) error
  DeleteFriendship(friendship *entity.Connection) error
}

type EmailServiceInterface interface {
  // Send all emails if at least one fails return error with explanation.
  SendEmails(email ...dto.Email) error
  // Send email with welcome text for user with 'CONFIRMED' email.
  SendWelcomeEmail(user dto.RequestUserDto) error
  // Send email with request for email confirmation to User with unconfirmed email.
  SendRegistrationConfirmationEmail(user dto.RequestUserDto) error
  // Send email with confirmation for password change.
  SendChangePasswordConfirmationEmail(user dto.RequestUserDto, newPassword string) error
}

type EmailBodyGeneratorInterface interface {
  // Generate message body for welcome email.
  GenerateWelcomeBody(user dto.RequestUserDto) string
  // Generate message body for registration confirmation email where token is confirmation token.
  GenerateRegistrationConfirmationEmail(user dto.RequestUserDto, token string) string
  // Generate message body for password change confirmation where token is confirmation token.
  GenerateChangePasswordConfirmationEmail(user dto.RequestUserDto, token string) string
}

type TokenGeneratorInterface interface {
  //Generate Token For Email And Save it To Database change user status to inactive
  GenerateTokenForEmail(email string) (string, error)
  // Change status for token to inActive and update user account status to active
  // If token wasn't found return false
  // If successfully update user status to Active return true
  ApproveUser(token string) (bool, error)
}

type EmailSendInterface interface {
  // Send email
  SendEmail(email dto.Email) error
}
