package dto

import (
  "fmt"
  "go-team-room/models/dao"
)

type ResponseUserDto struct {
  ID        int64   `json:"id"`
  Email     string  `json:"email"`
  FirstName string  `json:"firstName"`
  LastName  string  `json:"lastName"`
  Phone     string  `json:"phone"`
  Friends   []int64 `json:"friends"`
}

func (user ResponseUserDto) String() string {
  return fmt.Sprintf("User object:\n\tID = %d\n\tEmail = %s\n\tFirstName = %s\n\tSecondName = %s\n\tPhone = %s\n\tFriends = %v\n",
    user.ID, user.Email, user.FirstName, user.LastName, user.Phone, user.Friends)
}

type RequestUserDto struct {
  Email       string `json:"email"`
  FirstName   string `json:"firstName"`
  LastName    string `json:"lastName"`
  Phone       string `json:"phone"`
  CurrentPass string `json:"password"`
}

func (user RequestUserDto) String() string {
  return fmt.Sprintf("User object:\n\tEmail = %s\n\tFirstName = %s\n\tSecondName = %s\n\tPhone = %s\n\tPassword = %s\n",
    user.Email, user.FirstName, user.LastName, user.Phone, user.CurrentPass)
}

func RequestUserDtoToDao(userDto *RequestUserDto) dao.User {
  userDao := dao.User {
    0,
    userDto.Email,
    userDto.FirstName,
    userDto.LastName,
    userDto.Phone,
    userDto.CurrentPass,
    dao.UserRole,
    dao.Active,
    "",
  }

  return userDao
}

//without friends
func UserDaoToResponseDto(userDao *dao.User) ResponseUserDto {
  userDto := ResponseUserDto{
    userDao.ID,
    userDao.Email,
    userDao.FirstName,
    userDao.SecondName,
    userDao.Phone,
    []int64 {},
  }

  return userDto
}
