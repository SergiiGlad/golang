package controllers

import (
  "go-team-room/models/dto"
  "go-team-room/models/dao/entity"
)

//UserServiceInterface interface is used as HandlerFunc wrappers to follow dependency inversion principle
type UserServiceInterface interface {
  CreateUser(userDto *dto.RequestUserDto) (dto.ResponseUserDto, error)
  UpdateUser(id int64, userDto *dto.RequestUserDto) (dto.ResponseUserDto, error)
  DeleteUser(id int64) (dto.ResponseUserDto, error)
}

type FriendServiceInterface interface {
  GetUserFriends(id int64) ([]entity.User, error)
  GetUserFriedsIds(id int64)([]int64, error)
}
