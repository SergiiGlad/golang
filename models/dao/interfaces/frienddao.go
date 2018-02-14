package interfaces

import "go-team-room/models/dao/entity"

type FriendDao interface {
  InsertConnection(connection *entity.Connection) error
  UpdateStatus(connection *entity.Connection) error
  Delete(connection *entity.Connection) error
  FriendsByUserID(id int64) ([]entity.User, error)
  UsersWithRequestsTo(id int64) ([]entity.User, error)
  FindConnection(connection *entity.Connection) (entity.Connection, error)
}
