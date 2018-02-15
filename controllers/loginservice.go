package controllers

import (
  "go-team-room/models/dao/entity"
  "go-team-room/models/dao/mysql"
  "gopkg.in/hlandau/passlib.v1/hash/bcrypt"
  "errors"
)

func Login(phoneOrEmail string, password string) (*entity.User, error) {
  var user entity.User
  var err error

  if ValidPhone(phoneOrEmail) {
    user, err = mysql.UserDao.FindUserByPhone(phoneOrEmail)
  } else {
    user, err = mysql.UserDao.FindUserByEmail(phoneOrEmail)
  }

  if err != nil {
    return nil, errors.New("wrong credentials")
  }

  pass, err := mysql.PasswordDao.LastPassByUserId(user.ID)

  if err != nil {
    return nil, err
  }

  if bcrypt.Crypter.Verify(password, pass.Password) != nil {
    return nil, errors.New("wrong credentials")
  }

  return &user, nil
}
