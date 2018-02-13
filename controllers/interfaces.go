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
  GetFriends(id int64) ([]dto.ShortUser, error)
  GetUsersWithRequests(id int64) ([]dto.ShortUser, error)
  GetFriendIds(id int64) ([]int64, error)
  NewFriendRequest(connection *entity.Connection) error
  ApproveFriendRequest(connection *entity.Connection) error
  RejectFriendRequest(connection *entity.Connection) error
  DeleteFriendship(friendship *entity.Connection) error
}
