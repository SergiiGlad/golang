package server

import (
  "net/http"
  "go-team-room/controllers"
  "go-team-room/models/dao/mysql"
  "go-team-room/models/dao/entity"
  "io"
  "crypto/rand"
  "gopkg.in/hlandau/passlib.v1/hash/bcrypt"
  "time"
  "go-team-room/models/dto"
)

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func NewPassword(length int) string {
  return rand_char(length, StdChars)
}

func rand_char(length int, chars []byte) string {
  new_pword := make([]byte, length)
  random_data := make([]byte, length+(length/4)) // storage for random bytes.
  clen := byte(len(chars))

  maxrb := byte(256 - (256 % len(chars)))
  i := 0
  for {
    if _, err := io.ReadFull(rand.Reader, random_data); err != nil {
      panic(err)
    }
    for _, c := range random_data {
      if c >= maxrb {
        continue
      }
      new_pword[i] = chars[c%clen]
      i++
      if i == length {
        return string(new_pword)
      }
    }
  }
  panic("unreachable")
}

func recoveryPass(service controllers.UserServiceInterface, emailService controllers.UserEmailServiceInterface) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    email := r.Form["email"][0]
    user, err := mysql.UserDao.FindUserByEmail(email)

    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }
    userId := user.ID

    newPass := NewPassword(6)
    hashPass, err := bcrypt.Crypter.Hash(newPass)
    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }

    newPassStruct := entity.Password{
      0,
      hashPass,
      time.Now().Format("2006-01-02 15:04:05"),
      userId,
    }
    _, err = mysql.PasswordDao.InsertPass(&newPassStruct)
    if err != nil {
      responseError(w, err, http.StatusForbidden)
      return
    }
    emailService.SendChangePasswordConfirmationEmail(dto.RequestUserDto{
      Email:     user.Email,
      FirstName: user.FirstName,
      LastName:  user.LastName,
      Phone:     user.Phone,
      Role:      user.Role,
      Password:  hashPass,
    }, newPass)
    log.Println(newPassStruct)
  }
}
