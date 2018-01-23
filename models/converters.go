package models

import (
  "go-team-room/models/dao"
  "go-team-room/models/dto"
)

func RequestUserDtoToDao(userDto dto.RequestUserDto) dao.User {
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
func UserDaoToResponseDto(userDao dao.User) dto.ResponseUserDto {
  userDto := dto.ResponseUserDto{
    userDao.ID,
    userDao.Email,
    userDao.FirstName,
    userDao.SecondName,
    userDao.Phone,
    []int64 {},
  }

  return userDto
}

