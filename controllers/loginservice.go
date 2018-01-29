package controllers

import (
  "go-team-room/models/dao"
  "go-team-room/models/dao/mysql"
  "golang.org/x/crypto/bcrypt"
  "errors"
)

func Login(email string, password string) (*dao.User, error) {
  user, err := mysql.DB.FindUserByEmail(email)

  if err != nil {
    return nil, err
  }

  pass, err := mysql.DB.LastPassByUserId(user.ID)

  if err != nil {
    return nil, err
  }

  if bcrypt.CompareHashAndPassword([]byte(pass.Password), []byte(password)) != nil {
    return nil, errors.New("wrong password")
  }

  return &user, nil
}
