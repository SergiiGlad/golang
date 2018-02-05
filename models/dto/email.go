package dto

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
