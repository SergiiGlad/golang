package controllers

import (
  "errors"
  "gopkg.in/hlandau/passlib.v1/hash/bcrypt"
  "database/sql"
  "strings"
  "go-team-room/models/dto"
  "go-team-room/models/dao/interfaces"
  "time"
  "go-team-room/models/dao/entity"
  "go-team-room/conf"
)

// Get instance of logger (Formatter, Hook，Level, Output ).
// If you want to use only your log message  It will need use own call logs example
var log = conf.GetLog()

//UserService type implements UserServiceInterface and holds one field DB to access to database
type UserService struct {
  DB interfaces.MySqlDal
}

var _ UserServiceInterface = &UserService{}

func (us *UserService) CreateUser(userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {

  var responseUserDto dto.ResponseUserDto
  err := CheckUniqueEmail(userDto.Email, us.DB)
  if err != nil && err != sql.ErrNoRows {
    return responseUserDto, err
  }

  err = CheckUniquePhone(userDto.Phone, us.DB)
  if err != nil && err != sql.ErrNoRows {
    return responseUserDto, err
  }

  if ValidPasswordLength(userDto.Password) == false {
    return responseUserDto, errors.New("Password too short.")

  }

  hashPass, err := bcrypt.Crypter.Hash(userDto.Password)
  if err != nil {
    return responseUserDto, err
  }

  userEntity := dto.RequestUserDtoToEntity(userDto)
  NameLetterToUppep(&userEntity)
  user, err := us.DB.AddUser(&userEntity)
  if err != nil {
    return responseUserDto, err
  }

  newPass := entity.Password{
    0,
    hashPass,
    time.Now().Format("2006-01-02 15:04:05"),
    user.ID,
  }

  _, err = us.DB.InsertPass(&newPass)
  if err != nil {
    us.DB.ForceDeleteUser(user.ID)
    return responseUserDto, err
  }

  responseUserDto = dto.UserEntityToResponseDto(&userEntity)

  return responseUserDto, nil
}

func (us *UserService) UpdateUser(id int64, userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {

  oldUserData, err := us.DB.FindUserById(id)
  var responseUserDto dto.ResponseUserDto
  if err != nil {
    return responseUserDto, err
  }

  if len(userDto.FirstName) == 0 {
    userDto.FirstName = oldUserData.FirstName
  }

  if len(userDto.LastName) == 0 {
    userDto.LastName = oldUserData.LastName
  }

  if len(userDto.Email) != 0 {
    err = CheckUniqueEmail(userDto.Email, us.DB)
    if err != nil && err != sql.ErrNoRows {
      return responseUserDto, err
    }
  } else {
    userDto.Email = oldUserData.Email
  }

  if len(userDto.Phone) != 0 {
    err = CheckUniquePhone(userDto.Phone, us.DB)
    if err != nil && err != sql.ErrNoRows {
      return responseUserDto, err
    }
  } else {
    userDto.Phone = oldUserData.Phone
  }

  err = us.newPassIfValid(id, userDto.Password)
  if err != nil {
    return responseUserDto, err
  }

  newUserData := dto.RequestUserDtoToEntity(userDto)
  NameLetterToUppep(&newUserData)

  _, err = us.DB.UpdateUser(id, &newUserData)
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
  userEntity, err := us.DB.FindUserById(id)
  if userEntity.Role == entity.AdminRole {
    admins, err := us.DB.CountByRole(entity.AdminRole)
    if err != nil {
      return responseUserDto, err
    }

    if admins == 1 {
      return responseUserDto, errors.New("could not delete user with admin status")
    }
  }

  if err != nil {
    return responseUserDto, err
  }

  responseUserDto = dto.UserEntityToResponseDto(&userEntity)
  responseUserDto.Friends, _ = us.GetUserFriends(id)

  return responseUserDto, us.DB.DeleteUser(id)
}

func (us *UserService) GetUserFriends(id int64) ([]int64, error) {

  _, err := us.DB.FindUserById(id)
  if err != nil {
    return nil, err
  }

  return us.DB.FriendsByUserID(id)
}

//CheckUniqueEmail validates email string and queries to database to make sure that email is unique
func CheckUniqueEmail(email string, dao interfaces.UserDao) error {

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

//CheckUniquePhone validates phone string and queries to database to make sure that input phone is unique
func CheckUniquePhone(phone string, dao interfaces.UserDao) error {
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



//newPassIfValid method validate and create new password. It just checks password length and
//if length ok then password should be hashed and written into database
func (us *UserService) newPassIfValid(userId int64, password string) error {

  if ValidPasswordLength(password) == false {
    return errors.New("Password too short.")
  }

  hashPass, err := bcrypt.Crypter.Hash(password)
  if err != nil {
    return err
  }

  newPass := entity.Password{
    0,
    hashPass,
    time.Now().Format("2006-01-02 15:04:05"),
    userId,
  }

  _, err = us.DB.InsertPass(&newPass)
  if err != nil {
    return err
  }

  return nil
}

//NameLetterToUppep makes sure first letters of user name and surname will be upper case
func NameLetterToUppep(user *entity.User) {
  user.FirstName = strings.ToUpper(string([]rune(user.FirstName)[0])) + string([]rune(user.FirstName)[1:])
  user.LastName = strings.ToUpper(string([]rune(user.LastName)[0])) + string([]rune(user.LastName)[1:])
}
