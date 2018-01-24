package controllers

import (
  "go-team-room/models/dao"
  "regexp"
  "errors"
  "gopkg.in/hlandau/passlib.v1/hash/bcrypt"
  "log"
  "database/sql"
  "go-team-room/models/dao/mysql"
  "strings"
)

func CreateUser(user *dao.User) error {
  err := checkUniqueEmail(user.Email)

  if err != nil && err != sql.ErrNoRows {
    log.Println(err)
    return err
  }

  err = checkUniquePhone(user.Phone)

  if err != nil && err != sql.ErrNoRows {
    return err
  }

  if validPasswordLength(user.CurrentPass) == false {
    return errors.New("Password too short.")
  }

  hashPass, err := bcrypt.Crypter.Hash(user.CurrentPass)

  if err != nil {
    log.Println(err)
    return err
  }

  user.CurrentPass = hashPass
  nameLetterToUppep(user)

  id, err := mysql.DB.AddUser(user)

  if err != nil {
    log.Println(err)
    return err
  }

  user.ID = id

  return nil
}

func UpdateUser(id int64, user *dao.User) error {

  _, err := mysql.DB.FindUserById(id)

  if err != nil {
    return err
  }

  err = checkUniqueEmail(user.Email)

  if err != nil && err != sql.ErrNoRows {
    return err
  }

  err = checkUniquePhone(user.Phone)

  if err != nil && err != sql.ErrNoRows {
    return err
  }

  if validPasswordLength(user.CurrentPass) == false {
    return errors.New("Password too short.")
  }

  hashPass, err := bcrypt.Crypter.Hash(user.CurrentPass)

  if err != nil {
    log.Println(err)
    return err
  }

  user.CurrentPass = hashPass
  nameLetterToUppep(user)

  err = mysql.DB.UpdateUser(id, user)

  if err != nil {
    log.Println(err)
    return err
  }

  return nil
}

func DeleteUser(id int64) (*dao.User, error) {
  user, err := mysql.DB.FindUserById(id)

  if err != nil {
    return user, err
  }

  return user, mysql.DB.DeleteUser(id)
}

func checkUniqueEmail(email string) error {


  if validEmail(email) == false {
    return errors.New("Invalid email format.")
  } else {

    _, err := mysql.DB.FindUserByEmail(email)

    switch err {
    case sql.ErrNoRows:
      return err

    case nil:
      return errors.New("There is user with such email.")

    default:
      log.Println(err)
      return err
    }
  }

  return nil
}

func checkUniquePhone(phone string) error {
  if len(phone) > 0 {

    if validPhone(phone) == false {
      return errors.New("Invalid phone number format.")
    } else {

      _, err := mysql.DB.FindUserByPhone(phone)

      switch err {
      case sql.ErrNoRows:
        return err

      case nil:
        return errors.New("There is user with such phone.")

      default:
        log.Println(err)
        return err
      }
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

func nameLetterToUppep(user *dao.User) {
  user.FirstName = strings.ToUpper(string([]rune(user.FirstName)[0])) + string([]rune(user.FirstName)[1:])
  user.SecondName = strings.ToUpper(string([]rune(user.SecondName)[0])) + string([]rune(user.SecondName)[1:])
}

