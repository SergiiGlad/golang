package controllers

import "go-team-room/models/dto"

type UserServiceInterface interface {
  CreateUser(userDto *dto.RequestUserDto) (dto.ResponseUserDto, error)
  UpdateUser(id int64, userDto *dto.RequestUserDto) (dto.ResponseUserDto, error)
  DeleteUser(id int64) (dto.ResponseUserDto, error)
  GetUserFriends(id int64) ([]int64, error)
}
