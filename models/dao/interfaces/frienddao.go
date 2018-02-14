package interfaces

import "go-team-room/models/dao/entity"

type FriendDao interface {
  InsertConnection(friendship *entity.Connection) error
  UpdateStatus(friendship *entity.Connection) error
  Delete(friendship *entity.Connection) error
  FriendsByUserID(id int64) ([]entity.User, error)
  UsersWithRequestsTo(id int64) ([]entity.User, error)
  FindConnection(connection *entity.Connection) (entity.Connection, error)
}
