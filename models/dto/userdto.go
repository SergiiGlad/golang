package dto

import (
  "fmt"
  "go-team-room/models/dao"
)

type LoginDto struct {
  PhoneOrEmail string `json:"phoneOrEmail"`
  Password     string `json:"password"`
}

func (login LoginDto) String() string {
  return fmt.Sprint("Login object:\n\tLogin = %s\n\tPassword = %s\n\t", login.PhoneOrEmail, login.Password)
}

type ResponseUserDto struct {
  ID        int64   `json:"id"`
  Email     string  `json:"email"`
  FirstName string  `json:"first_name"`
  LastName  string  `json:"last_name"`
  Phone     string  `json:"phone"`
  Friends   []int64 `json:"friends"`
}

func (user ResponseUserDto) String() string {
  return fmt.Sprintf("User object:\n\tID = %d\n\tEmail = %s\n\tFirstName = %s\n\tLastName = %s\n\tPhone = %s\n\tFriends = %v\n",
    user.ID, user.Email, user.FirstName, user.LastName, user.Phone, user.Friends)
}

type RequestUserDto struct {
  Email     string `json:"email"`
  FirstName string `json:"first_name"`
  LastName  string `json:"last_name"`
  Phone     string `json:"phone"`
  Password  string `json:"password"`
}

func (user RequestUserDto) String() string {
  return fmt.Sprintf("User object:\n\tEmail = %s\n\tFirstName = %s\n\tLastName = %s\n\tPhone = %s\n\tPassword = %s\n",
    user.Email, user.FirstName, user.LastName, user.Phone, user.Password)
}

func RequestUserDtoToEntity(userDto *RequestUserDto) dao.User {
  userDao := dao.User {
    0,
    userDto.Email,
    userDto.FirstName,
    userDto.LastName,
    userDto.Phone,
    dao.UserRole,
    dao.Active,
    "",
  }

  return userDao
}

//without friends
func UserEntityToResponseDto(userDao *dao.User) ResponseUserDto {
  userDto := ResponseUserDto{
    userDao.ID,
    userDao.Email,
    userDao.FirstName,
    userDao.LastName,
    userDao.Phone,
    []int64 {},
  }

  return userDto
}
