//Package dto provides data transfer objects used in this project.
package dto

import (
  "fmt"
  "go-team-room/models/dao/entity"
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
  Avatar    string  `json:"avatar_ref"`
  Friends   int64   `json:"friends_num"`
}

func (user ResponseUserDto) String() string {
  return fmt.Sprintf("User object:\n\tID = %d\n\tEmail = %s\n\tFirstName = %s\n\tLastName = %s\n\tPhone = %s\n\tFriends = %v\n",
    user.ID, user.Email, user.FirstName, user.LastName, user.Phone, user.Friends)
}

//RequestUserDto is used in converting from json structure pulled from request body.
type RequestUserDto struct {
  Email     string      `json:"email"`
  FirstName string      `json:"first_name"`
  LastName  string      `json:"last_name"`
  Phone     string      `json:"phone"`
  Role      entity.Role `json:"role"`
  Password  string      `json:"password"`
  Avatar    string      `json:"avatar_ref"`
}

func (user RequestUserDto) String() string {
  return fmt.Sprintf("User object:\n\tEmail = %s\n\tFirstName = %s\n\tLastName = %s\n\tPhone = %s\n\tPassword = %s\n",
    user.Email, user.FirstName, user.LastName, user.Phone, "XXXXXX")
}

type ShortUser struct {
  ID        int64   `json:"id"`
  FirstName string  `json:"first_name"`
  LastName  string  `json:"last_name"`
  Avatar    string  `json:"avatar_ref"`
}

func (user ShortUser) String() string {
  return fmt.Sprintf("User object:\n\tID = %d\n\tFirstName = %s\n\tLastName = %s\n\tAvaterRef = %s\n",
    user.ID, user.FirstName, user.FirstName, user.Avatar)
}

func UserEntityToShortUser(userDao *entity.User) ShortUser {
  return ShortUser {userDao.ID, userDao.FirstName, userDao.LastName, userDao.AvatarRef}
}

//RequestUserDtoToEntity converts RequestUserDto to dao.User. user role is set by default if such field is empty.
func RequestUserDtoToEntity(userDto *RequestUserDto) entity.User {
  if userDto.Role == "" {
    userDto.Role = entity.UserRole
  }

  userDao := entity.User {
    0,
    userDto.Email,
    userDto.FirstName,
    userDto.LastName,
    userDto.Phone,
    userDto.Role,
    entity.InActive,
    userDto.Avatar,
  }

  return userDao
}

//UserEntityToResponseDto converts *dao.User to ResponseUserDto. For now Friends field is empty slice.
func UserEntityToResponseDto(userDao *entity.User) ResponseUserDto {
  userDto := ResponseUserDto{
    userDao.ID,
    userDao.Email,
    userDao.FirstName,
    userDao.LastName,
    userDao.Phone,
    userDao.AvatarRef,
    0,
  }

  return userDto
}
