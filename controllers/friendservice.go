package controllers

import (
  "go-team-room/models/dao/interfaces"
  "go-team-room/models/dao/entity"
  "github.com/pkg/errors"
  "database/sql"
  "go-team-room/models/dto"
)

type FriendService struct {
  FriendshipDao interfaces.FriendshipDao
  UserDao       interfaces.UserDao
}

func (fs *FriendService) GetFriends(id int64) ([]dto.ShortUser, error) {
  friends, err := fs.FriendshipDao.FriendsByUserID(id)
  shortFriends := []dto.ShortUser{}
  if err != nil {
    return shortFriends, err
  }

  shortFriends = usersToShortUsers(friends)

  return shortFriends, nil
}

func (fs *FriendService) GetUsersWithRequests(id int64) ([]dto.ShortUser, error) {
  friends, err := fs.FriendshipDao.UsersWithRequestsTo(id)
  shortFriends := []dto.ShortUser{}
  if err != nil {
    return shortFriends, err
  }

  shortFriends = usersToShortUsers(friends)

  return shortFriends, nil
}

func usersToShortUsers(users []entity.User) []dto.ShortUser {
  shortUsers := []dto.ShortUser{}
  shortUser := dto.ShortUser{}
  for _, friend := range users {
    shortUser = dto.UserEntityToShortUser(&friend)
    shortUsers = append(shortUsers, shortUser)
  }

  return shortUsers
}


func (fs *FriendService) GetFriendIds(id int64) ([]int64, error) {
  friends, err := fs.FriendshipDao.FriendsByUserID(id)
  if err != nil {
    return []int64{}, err
  }

  ids := []int64{}
  for _, friend := range friends {
    ids = append(ids, friend.ID)
  }

  return ids, nil
}

func (fs *FriendService) NewFriendRequest(connection *entity.Connection) error {
  reversedConn := reversedConnection(connection)
  connectionExists, err := fs.connectionExists(&reversedConn)
  if err != nil {
    return err
  }
  if connectionExists {
    log.Error(errors.New("Connection already exists"))
    return errors.New("Connection already exists")
  }

  connection.Status = entity.Waiting
  _, err = fs.setConnection(connection)
  if err != nil {
    return err
  }

  return nil
}

func (fs *FriendService) setConnection(connection *entity.Connection) (entity.Connection, error) {

  if err := fs.validateConnectionUsers(connection); err != nil {
    return *connection, err
  }

  connectionExists, err := fs.connectionExists(connection)
  if err != nil {
    return *connection, err
  }

  if connectionExists {
    log.Error(errors.New("Connection already exists"))
    return *connection, errors.New("Connection already exists")
  }

  if err := fs.FriendshipDao.InsertConnection(connection); err != nil {
    return *connection, err
  }

  return *connection, nil
}

func (fs *FriendService) ApproveFriendRequest(connection *entity.Connection) error {
  if err := fs.validateConnectionUsers(connection); err != nil {
    return err
  }

  connectionExists, err := fs.connectionExists(connection)
  if err != nil {
    return err
  }

  if connectionExists == false {
    log.Error(errors.New("Unable to approve non existing connection"))
    return errors.New("Unable to approve non existing connection")
  }

  connection.Status = entity.Approved

  err = fs.FriendshipDao.UpdateStatus(connection)
  if err != nil {
    return err
  }

  reversed := reversedConnection(connection)
  _, err = fs.setConnection(&reversed)
  if err != nil {
    return err
  }

  return err
}

func (fs *FriendService) RejectFriendRequest(connection *entity.Connection) error {
  if err := fs.validateConnectionUsers(connection); err != nil {
    return err
  }

  connectionExists, err := fs.connectionExists(connection)
  if err != nil {
    return err
  }

  if connectionExists == false {
    log.Error(errors.New("Unable to reject non existing connection"))
    return errors.New("Unable to reject non existing connection")
  }

  return fs.FriendshipDao.Delete(connection)
}

func (fs *FriendService) DeleteFriendship(friendship *entity.Connection) error {
  if err := fs.validateConnectionUsers(friendship); err != nil {
    return err
  }

  friendshipExists, err := fs.connectionsExistBoth(friendship)
  if err != nil {
    return err
  }

  if friendshipExists == false {
    log.Error(errors.New("Unable to delete non existing friendship"))
    return errors.New("Unable to delete non existing friendship")
  }

  return fs.FriendshipDao.Delete(friendship)
}

func (fs *FriendService) validateConnectionUsers(connection *entity.Connection) error {
  if connection.FriendUserId == connection.UserId {
    log.Error(errors.New("Invalid connection"))
    return errors.New("Invalid connection")
  }

  user, err := fs.UserDao.FindUserById(connection.FriendUserId)
  if err != nil && user.AccStatus == entity.Deleted {
    log.Error(errors.New("Invalid connection"))
    return errors.New("Invalid connection")
  }

  user, err = fs.UserDao.FindUserById(connection.UserId)
  if err != nil && user.AccStatus == entity.Deleted{
    log.Error(errors.New("Invalid connection"))
    return errors.New("Invalid connection")
  }

  return nil
}

func (fs *FriendService) connectionExists(connection *entity.Connection) (bool, error) {
  _, err := fs.FriendshipDao.FindConnection(connection)
  switch err {
  case sql.ErrNoRows:
    return false, nil
  case nil:
    return true, nil
  default:
    return false, err
  }
}

func (fs *FriendService) connectionsExistBoth(connection *entity.Connection) (bool, error) {
  connectionReversed := reversedConnection(connection)

  exists1, err := fs.connectionExists(connection)
  if err != nil {
    return false, err
  }

  exists2, err := fs.connectionExists(&connectionReversed)
  if err != nil {
    return false, err
  }

  return exists1 && exists2, nil
}

func reversedConnection(connection *entity.Connection) entity.Connection {
  return entity.Connection{
    connection.UserId,
    connection.FriendUserId,
    connection.Status,
    }
}
