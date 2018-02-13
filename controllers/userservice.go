package controllers

import (
  "regexp"
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

//UserService type implements UserServiceInterface and holds one field UserDao to access to database
type UserService struct {
  FriendService FriendServiceInterface
  PassDao       interfaces.PasswordDao
  UserDao       interfaces.UserDao
}

var _ UserServiceInterface = &UserService{}

func (us *UserService) CreateUser(userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {

  var responseUserDto dto.ResponseUserDto
  err := CheckUniqueEmail(userDto.Email, us.UserDao)
  if err != nil && err != sql.ErrNoRows {
    return responseUserDto, err
  }

  err = CheckUniquePhone(userDto.Phone, us.UserDao)
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
  user, err := us.UserDao.AddUser(&userEntity)
  if err != nil {
    return responseUserDto, err
  }

  newPass := entity.Password{
    0,
    hashPass,
    time.Now().Format("2006-01-02 15:04:05"),
    user.ID,
  }

  _, err = us.PassDao.InsertPass(&newPass)
  if err != nil {
    us.UserDao.ForceDeleteUser(user.ID)
    return responseUserDto, err
  }

  responseUserDto = dto.UserEntityToResponseDto(&userEntity)

  return responseUserDto, nil
}

func (us *UserService) UpdateUser(id int64, userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {

  oldUserData, err := us.UserDao.FindUserById(id)
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
    err = CheckUniqueEmail(userDto.Email, us.UserDao)
    if err != nil && err != sql.ErrNoRows {
      return responseUserDto, err
    }
  } else {
    userDto.Email = oldUserData.Email
  }

  if len(userDto.Phone) != 0 {
    err = CheckUniquePhone(userDto.Phone, us.UserDao)
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

  _, err = us.UserDao.UpdateUser(id, &newUserData)
  if err != nil {
    log.Println(err)
    return responseUserDto, err
  }

  newUserData.ID = id
  responseUserDto = dto.UserEntityToResponseDto(&newUserData)
  friends, _ := us.FriendService.GetFriendIds(id)
  responseUserDto.Friends = int64(len(friends))

  return responseUserDto, nil
}

func (us *UserService) DeleteUser(id int64) (dto.ResponseUserDto, error) {

  var responseUserDto dto.ResponseUserDto
  userEntity, err := us.UserDao.FindUserById(id)
  if userEntity.Role == entity.AdminRole {
    admins, err := us.UserDao.CountByRole(entity.AdminRole)
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
  friends, _ := us.FriendService.GetFriendIds(id)
  responseUserDto.Friends = int64(len(friends))

  return responseUserDto, us.UserDao.DeleteUser(id)
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

//ValidRegexItem is helper function that is used by ValidEmail and ValidPhone as common logic for both.
func ValidRegexItem(item string, pattern string) bool {

  itemRegex := regexp.MustCompile(pattern)
  if isItemOk := itemRegex.MatchString(item); isItemOk == false {
    return false
  }

  return true
}

//ValidEmail return ValidRegexItem function result that checks email string argument
// matches the email regex pattern
func ValidEmail(email string) bool {
  return ValidRegexItem(email, "^[a-z0-9]+@[a-z]+[.][a-z]+$")
}

//ValidPhone return ValidRegexItem function result that checks phone string argument
// matches the phone regex pattern
func ValidPhone(phone string) bool {
  return ValidRegexItem(phone, "^[+][0-9]{12}$")
}

//ValidCyrillicName return ValidRegexItem result that checks if name string is valid cyrillic name pattern
func ValidCyrillicName(name string) bool {
  return ValidRegexItem(name, "^[А-Я][а-я]{1,49}$")
}

//ValidLatinName return ValidRegexItem result that checks if name string is valid latin name pattern
func ValidLatinName(name string) bool {
  return ValidRegexItem(name, "^[A-Z][a-z]{1,49}$")
}

//ValidPasswordLength checks password length is bigger than 6 symbols
func ValidPasswordLength(password string) bool {
  if len(password) < 6 {
    return false
  }

  return true
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

  _, err = us.PassDao.InsertPass(&newPass)
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

