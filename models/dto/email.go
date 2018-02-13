package dto

import "fmt"

type Email struct {
  To      string
  Subject string
  Body    string
}

func RequestUserDtoToEmail(user RequestUserDto, subject string, body string) Email {
  return Email{
    To:      user.Email,
    Subject: subject,
    Body:    body,
  }
}

func (e Email) String() string {
  return fmt.Sprintf("Email objetct: to=%s, subject=%s", e.To, e.Subject)
}
