package controllers

import (
  "go-team-room/models/dao"
  "go-team-room/models/dao/mysql"
  "gopkg.in/hlandau/passlib.v1/hash/bcrypt"
  "errors"
)

func Login(phoneOrEmail string, password string) (*dao.User, error) {
  var user dao.User
  var err error

  if validPhone(phoneOrEmail) {
    user, err = mysql.DB.FindUserByPhone(phoneOrEmail)
  } else {
    user, err = mysql.DB.FindUserByEmail(phoneOrEmail)
  }

  if err != nil {
    return nil, errors.New("wrong credentials")
  }

  pass, err := mysql.DB.LastPassByUserId(user.ID)

  if err != nil {
    return nil, err
  }

  if bcrypt.Crypter.Verify(password, pass.Password) != nil {
    return nil, errors.New("wrong credentials")
  }

  return &user, nil
}
