package controllers

import (
  "go-team-room/models/dto"
  "go-team-room/models/dao/entity"
  "github.com/pkg/errors"
)

type friendServiceMock struct {}

func (fs friendServiceMock) GetFriends(id int64) ([]dto.ShortUser, error) {
  if id < 1 {
    return nil, errors.New("Invalid id")
  }

  return []dto.ShortUser{}, nil
}

func (fs friendServiceMock) GetUsersWithRequests(id int64) ([]dto.ShortUser, error) {
  if id < 1 {
    return nil, errors.New("Invalid id")
  }

  return []dto.ShortUser{}, nil
}

func (fs friendServiceMock) GetFriendIds(id int64) ([]int64, error) {
  if id < 1 {
    return nil, errors.New("Invalid id")
  }

  return []int64{}, nil
}

func (fs friendServiceMock) NewFriendRequest(connection *entity.Connection) error {
  return nil
}

func (fs friendServiceMock) ApproveFriendRequest(connection *entity.Connection) error {
  return nil
}

func (fs friendServiceMock) RejectFriendRequest(connection *entity.Connection) error {
  return nil
}

func (fs friendServiceMock) DeleteFriendship(friendship *entity.Connection) error {
  return nil
}

var friendServiceMocked = friendServiceMock{}
