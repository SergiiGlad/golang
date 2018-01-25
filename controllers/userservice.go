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

  userEntity := dto.RequestUserDtoToEntity(userDto)
  var respUserDto dto.ResponseUserDto

  err := checkUniqueEmail(userEntity.Email)

  if err != nil && err != sql.ErrNoRows {
    return respUserDto, err
  }

  err = checkUniquePhone(userEntity.Phone)

  if err != nil && err != sql.ErrNoRows {
    return respUserDto, err
  }

  if validPasswordLength(userEntity.CurrentPass) == false {
    return respUserDto, errors.New("Password too short.")
  }

  hashPass, err := bcrypt.Crypter.Hash(userEntity.CurrentPass)

  if err != nil {
    return respUserDto, err
  }

  userEntity.CurrentPass = hashPass
  nameLetterToUppep(&userEntity)

  _, err = mysql.DB.AddUser(&userEntity)

  if err != nil {
    log.Println(err)
    return respUserDto, err
  }

  respUserDto = dto.UserEntityToResponseDto(&userEntity)

  return respUserDto, nil
}

func UpdateUser(id int64, userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {

  newUserData := dto.RequestUserDtoToEntity(userDto)
  oldUserData, err := mysql.DB.FindUserById(id)
  var responseUser dto.ResponseUserDto

  if err != nil {
    return responseUser, err
  }

  if len(newUserData.FirstName) == 0 {
    newUserData.FirstName = oldUserData.FirstName
  }

  if len(newUserData.SecondName) == 0 {
    newUserData.SecondName = oldUserData.SecondName
  }

  if len(newUserData.Email) != 0 {
    err = checkUniqueEmail(newUserData.Email)

    if err != nil && err != sql.ErrNoRows {
      return responseUser, err
    }
  } else {
    newUserData.Email = oldUserData.Email
  }

  if len(newUserData.Phone) != 0 {
    err = checkUniquePhone(newUserData.Phone)

    if err != nil && err != sql.ErrNoRows {
      return responseUser, err
    }
  } else {
    newUserData.Phone = oldUserData.Phone
  }

  if len(newUserData.CurrentPass) != 0 {
    if validPasswordLength(newUserData.CurrentPass) == false {
      return responseUser, errors.New("Password too short.")
    }

    hashPass, err := bcrypt.Crypter.Hash(newUserData.CurrentPass)

    if err != nil {
      log.Println(err)
      return responseUser, err
    }

    newUserData.CurrentPass = hashPass
  } else {
    newUserData.CurrentPass = oldUserData.CurrentPass
  }

  nameLetterToUppep(&newUserData)

  err = mysql.DB.UpdateUser(id, &newUserData)
  if err != nil {
    log.Println(err)
    return responseUser, err
  }

  newUserData.ID = id
  responseUser = dto.UserEntityToResponseDto(&newUserData)

  return responseUser, nil
}

func DeleteUser(id int64) (dto.ResponseUserDto, error) {

  var responseUserDto dto.ResponseUserDto

  userEntity, err := mysql.DB.FindUserById(id)

  if err != nil {
    return responseUserDto, err
  }

  responseUserDto = dto.UserEntityToResponseDto(userEntity)
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

