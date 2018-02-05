package controllers

import (
  "go-team-room/models/dto"
  "github.com/matcornic/hermes"
  "go-team-room/conf"
)

//Generate emails bodys using Hermes package
type HermesEmailBodyGenerator struct {
}

var h = hermes.Hermes{
  Theme: new(hermes.Flat),
  Product: hermes.Product{
    Name: "GoHum",
    Link: "https://gohum.go",
  },
}

var _ EmailBodyGeneratorInterface = &HermesEmailBodyGenerator{}

// Generate message body for welcome email.
func (eg HermesEmailBodyGenerator) GenerateWelcomeBody(user dto.RequestUserDto) string {
  email := hermes.Email{
    Body: hermes.Body{
      Name: user.FirstName + " " + user.LastName,
      Intros: []string{
        "Welcome to GoHum! Thank you for your registration! Your account details: ",
      },
      Dictionary: []hermes.Entry{
        {Key: "Firstname", Value: user.FirstName},
        {Key: "Lastname", Value: user.LastName},
        {Key: "Email", Value: user.Email},
        {Key: "Phone", Value: user.Phone},
      },
      Actions: []hermes.Action{
        {
          Instructions: "Login to your account",
          Button: hermes.Button{
            Color: "#32CD32",
            Text:  "Login",
            Link:  conf.Ip + "/dist/login", // TODO Add property dns login
          },
        },
      },
      Outros: []string{
        "Need help, or have questions? Just reply to this email, we'd love to help.",
      },
    },
  }
  body, _ := h.GenerateHTML(email)
  return body
}

// Generate message body for registration confirmation email.
func (eg HermesEmailBodyGenerator) GenerateRegistrationConfirmationEmail(user dto.RequestUserDto, token string) string {
  email := hermes.Email{
    Body: hermes.Body{
      Name: user.FirstName + " " + user.LastName,
      Intros: []string{
        "Only one step remains!",
      },
      Actions: []hermes.Action{
        {
          Instructions: "To get started with GoHum, please confirm your email:",
          Button: hermes.Button{
            Color: "#1AACF5", // Optional action button color
            Text:  "Confirm your email",
            Link:  conf.Ip + "/confirm/email?token=" + token, // TODO Add property dns token
          },
        },
      },
      Outros: []string{
        "Need help, or have questions? Just reply to this email, we'd love to help.",
      },
    },
  }
  body, _ := h.GenerateHTML(email)
  return body
}

// Generate message body for password change confirmation.
func (eg HermesEmailBodyGenerator) GenerateChangePasswordConfirmationEmail(user dto.RequestUserDto,
  token string) string {
  email := hermes.Email{
    Body: hermes.Body{
      Name: user.FirstName + " " + user.LastName,
      Intros: []string{
        "Attention! Your About to change your password!",
      },
      Actions: []hermes.Action{
        {
          Instructions: "To confirm your new password press button. If your not changing password contact as immediately!",
          Button: hermes.Button{
            Color: "#FF0000", // Optional action button color
            Text:  "I confirm password change!",
            Link:  conf.Ip + "/confirm/password?token=" + token,
          },
        },
      },
      Outros: []string{
        "Need help, or have questions? Just reply to this email, we'd love to help.",
      },
    },
  }
  body, _ := h.GenerateHTML(email)
  return body
}
