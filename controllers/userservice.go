package controllers

import (
  "go-team-room/models/dao"
  "errors"
  "gopkg.in/hlandau/passlib.v1/hash/bcrypt"
  "log"
  "database/sql"
  "strings"
  "go-team-room/models/dto"
  "go-team-room/models/dao/interfaces"
)

type UserService struct {
  Dao interfaces.UserDao
}

var _ UserServiceInterface = &UserService{}

func (us *UserService) CreateUser(userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {

  userEntity := dto.RequestUserDtoToEntity(userDto)
  var respUserDto dto.ResponseUserDto

  err := checkUniqueEmail(userEntity.Email, us.Dao)

  if err != nil && err != sql.ErrNoRows {
    return respUserDto, err
  }

  err = checkUniquePhone(userEntity.Phone, us.Dao)

  if err != nil && err != sql.ErrNoRows {
    return respUserDto, err
  }

  if ValidPasswordLength(userEntity.CurrentPass) == false {
    return respUserDto, errors.New("Password too short.")
  }

  hashPass, err := bcrypt.Crypter.Hash(userEntity.CurrentPass)

  if err != nil {
    return respUserDto, err
  }

  userEntity.CurrentPass = hashPass
  nameLetterToUppep(&userEntity)

  _, err = us.Dao.AddUser(&userEntity)

  if err != nil {
    log.Println(err)
    return respUserDto, err
  }

  respUserDto = dto.UserEntityToResponseDto(&userEntity)

  return respUserDto, nil
}

func (us *UserService) UpdateUser(id int64, userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {

  newUserData := dto.RequestUserDtoToEntity(userDto)
  oldUserData, err := us.Dao.FindUserById(id)
  var responseUserDto dto.ResponseUserDto

  if err != nil {
    return responseUserDto, err
  }

  if len(newUserData.FirstName) == 0 {
    newUserData.FirstName = oldUserData.FirstName
  }

  if len(newUserData.SecondName) == 0 {
    newUserData.SecondName = oldUserData.SecondName
  }

  if len(newUserData.Email) != 0 {
    err = checkUniqueEmail(newUserData.Email, us.Dao)

    if err != nil && err != sql.ErrNoRows {
      return responseUserDto, err
    }
  } else {
    newUserData.Email = oldUserData.Email
  }

  if len(newUserData.Phone) != 0 {
    err = checkUniquePhone(newUserData.Phone, us.Dao)

    if err != nil && err != sql.ErrNoRows {
      return responseUserDto, err
    }
  } else {
    newUserData.Phone = oldUserData.Phone
  }

  if len(newUserData.CurrentPass) != 0 {
    if ValidPasswordLength(newUserData.CurrentPass) == false {
      return responseUserDto, errors.New("Password too short.")
    }

    hashPass, err := bcrypt.Crypter.Hash(newUserData.CurrentPass)

    if err != nil {
      log.Println(err)
      return responseUserDto, err
    }

    newUserData.CurrentPass = hashPass
  } else {
    newUserData.CurrentPass = oldUserData.CurrentPass
  }

  nameLetterToUppep(&newUserData)

  err = us.Dao.UpdateUser(id, &newUserData)
  if err != nil {
    log.Println(err)
    return responseUserDto, err
  }

  newUserData.ID = id
  responseUserDto = dto.UserEntityToResponseDto(&newUserData)
  responseUserDto.Friends, _ = us.GetUserFriends(id)

  return responseUserDto, nil
}

func (us *UserService) DeleteUser(id int64) (dto.ResponseUserDto, error) {

  var responseUserDto dto.ResponseUserDto

  userEntity, err := us.Dao.FindUserById(id)

  if err != nil {
    return responseUserDto, err
  }

  responseUserDto = dto.UserEntityToResponseDto(userEntity)
  responseUserDto.Friends, _ = us.GetUserFriends(id)

  return responseUserDto, us.Dao.DeleteUser(id)
}

func (us *UserService) GetUserFriends(id int64) ([]int64, error) {
  _, err := us.Dao.FindUserById(id)

  if err != nil {
    return nil, err
  }

  return us.Dao.FriendsByUserID(id)
}

func checkUniqueEmail(email string, dao interfaces.UserDao) error {

  if ValidEmail(email) == false {
    return errors.New("Invalid email format.")
  } else {

    _, err := dao.FindUserByEmail(email)

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

func checkUniquePhone(phone string, dao interfaces.UserDao) error {
  if len(phone) > 0 {

    if ValidPhone(phone) == false {
      return errors.New("Invalid phone number format.")
    } else {

      _, err := dao.FindUserByPhone(phone)

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

func nameLetterToUppep(user *dao.User) {
  user.FirstName = strings.ToUpper(string([]rune(user.FirstName)[0])) + string([]rune(user.FirstName)[1:])
  user.SecondName = strings.ToUpper(string([]rune(user.SecondName)[0])) + string([]rune(user.SecondName)[1:])
}
