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
  "go-team-room/models/dto"
)

func CreateUser(userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {

  user := dto.RequestUserDtoToDao(userDto)
  var responseUser dto.ResponseUserDto

  err := checkUniqueEmail(user.Email)

  if err != nil && err != sql.ErrNoRows {
    log.Println(err)
    return responseUser, err
  }

  err = checkUniquePhone(user.Phone)

  if err != nil && err != sql.ErrNoRows {
    return responseUser, err
  }

  if validPasswordLength(user.CurrentPass) == false {
    return responseUser, errors.New("Password too short.")
  }

  hashPass, err := bcrypt.Crypter.Hash(user.CurrentPass)

  if err != nil {
    log.Println(err)
    return responseUser, err
  }

  user.CurrentPass = hashPass
  nameLetterToUppep(&user)

  id, err := mysql.DB.AddUser(&user)

  if err != nil {
    log.Println(err)
    return responseUser, err
  }

  user.ID = id
  responseUser = dto.UserDaoToResponseDto(&user)

  return responseUser, nil
}

func UpdateUser(id int64, userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {

  userNew := dto.RequestUserDtoToDao(userDto)
  userOld, err := mysql.DB.FindUserById(id)
  var responseUser dto.ResponseUserDto

  if err != nil {
    return responseUser, err
  }

  if len(userNew.FirstName) == 0 {
    userNew.FirstName = userOld.FirstName
  }

  if len(userNew.SecondName) == 0 {
    userNew.SecondName = userOld.SecondName
  }

  if len(userNew.Email) != 0 {
    err = checkUniqueEmail(userNew.Email)

    if err != nil && err != sql.ErrNoRows {
      return responseUser, err
    }
  } else {
    userNew.Email = userOld.Email
  }

  if len(userNew.Phone) != 0 {
    err = checkUniquePhone(userNew.Phone)

    if err != nil && err != sql.ErrNoRows {
      return responseUser, err
    }
  } else {
    userNew.Phone = userOld.Phone
  }

  if len(userNew.CurrentPass) != 0 {
    if validPasswordLength(userNew.CurrentPass) == false {
      return responseUser, errors.New("Password too short.")
    }

    hashPass, err := bcrypt.Crypter.Hash(userNew.CurrentPass)

    if err != nil {
      log.Println(err)
      return responseUser, err
    }

    userNew.CurrentPass = hashPass
  } else {
    userNew.CurrentPass = userOld.CurrentPass
  }

  nameLetterToUppep(&userNew)

  err = mysql.DB.UpdateUser(id, &userNew)
  if err != nil {
    log.Println(err)
    return responseUser, err
  }

  userNew.ID = id
  responseUser = dto.UserDaoToResponseDto(&userNew)

  return responseUser, nil
}

func DeleteUser(id int64) (dto.ResponseUserDto, error) {

  var responseUserDto dto.ResponseUserDto

  user, err := mysql.DB.FindUserById(id)

  if err != nil {
    return responseUserDto, err
  }

  responseUserDto = dto.UserDaoToResponseDto(user)
  responseUserDto.Friends, _ = getUserFriends(id)

  return responseUserDto, mysql.DB.DeleteUser(id)
}

func getUserFriends(id int64) ([]int64, error) {
  _, err := mysql.DB.FindUserById(id)

  if err != nil {
    return nil, err
  }

  return mysql.DB.FriendsByUserID(id)
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

