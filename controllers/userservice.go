package controllers

import (
  "go-team-room/models/dao"
  "regexp"
  "errors"
  "gopkg.in/hlandau/passlib.v1/hash/bcrypt"
  "log"
  "database/sql"
  "go-team-room/models/dao/mysql"
)

func CreateUser(user *dao.User) error {
  err := checkEmail(user.Email)

  if err != nil {
    log.Fatal(err)
    return err
  }

  err = checkPhone(user.Phone)

  if err != nil {
    log.Fatal(err)
    return err
  }

  if validPasswordLength(user.CurrentPass) == false {
    log.Fatal("Invalid password size")
    return errors.New("Password too short.")
  }

  hashPass, err := bcrypt.Crypter.Hash(user.CurrentPass)

  if err != nil {
    log.Fatal(err)
    return err
  }

  user.CurrentPass = hashPass

  id, err := mysql.DB.AddUser(user)

  if err != nil {
    log.Fatal(err)
    return err
  }

  user.ID = id

  return nil
}

func UpdateUser(id int, user *dao.User) error {
  return mysql.DB.UpdateUser(id, user)
}

func DeleteUser(id int, user *dao.User) error {
  return mysql.DB.DeleteUser(id)
}

func checkEmail(email string) error {
  _, err := mysql.DB.FindUserByEmail(email)

  switch err {
  case sql.ErrNoRows:
    if validEmail(email) == false {
      log.Fatal("Invalid email format.")
      return errors.New("Invalid email format.")
    }

  case nil:
    log.Fatal("There is user with such email.")
    return errors.New("There is user with such email.")

  default:
    log.Fatal(err)
    return err
  }

  return nil
}

func checkPhone(phone string) error {
  if len(phone) > 0 {
    _, err := mysql.DB.FindUserByPhone(phone)

    switch err {
    case sql.ErrNoRows:
      if validPhone(phone) == false {
        log.Fatal("Invalid phone number format.")
        return errors.New("Invalid phone number format.")
      }

    case nil:
      log.Fatal("There is user with such phone.")
      return errors.New("There is user with such phone.")

    default:
      log.Fatal(err)
      return err
    }
  }

  return nil
}

func validRegexItem(item string, pattern string) bool {
  itemRegex := regexp.MustCompile(pattern)

  if isItemOk := itemRegex.MatchString(item); isItemOk == false {
    return false
  }

  return true
}

func validEmail(email string) bool {
  return validRegexItem(email, "^[a-z0-9]+@[a-z]+[.][a-z]+$")
}

func validPhone(phone string) bool {
  return validRegexItem(phone, "^[+][0-9]{12}")
}

func validCyrillicName(name string) bool {
  return validRegexItem(name, "^[А-Я][а-я]{1,49}")
}

func validLatinName(name string) bool {
  return validRegexItem(name, "^[A-Z][a-z]{1,49}")
}

func validPasswordLength(password string) bool {
  if len(password) < 6 {
    return false
  }

  return true
}
