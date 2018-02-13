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
  log.Debugf("Generating new Welcome Email body for user: %s", user)
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
  log.Debugf("Generating new Registration Confirmation Email body for user: %s, token: %s", user, token)
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
            Link:  conf.Ip + "/confirm/email/" + token, // TODO Add property dns token
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
  password string) string {
  log.Debugf("Generating new Change Password body for user: %s", user)
  email := hermes.Email{
    Body: hermes.Body{
      Name: user.FirstName + " " + user.LastName,
      Intros: []string{
        "Attention! This is your new password! Change it after login!",
      },
      Dictionary: []hermes.Entry{
        {Key: "New Password", Value: password},
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
