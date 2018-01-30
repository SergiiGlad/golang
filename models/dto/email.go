package dto

type Email struct {
  To      string
  Subject string
  Body    []byte
}

func RequestUserDtoToEmail(user RequestUserDto, subject string, body []byte) Email {
  return Email{
    To:      user.Email,
    Subject: subject,
    Body:    body,
  }
}
