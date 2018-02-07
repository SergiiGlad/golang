package interfaces

import (
  "go-team-room/models/dao/entity"
)

//UserDao interface is used in services to follow dependency inversion principle
type UserDao interface {
  AddUser(user *entity.User) (entity.User, error)
  DeleteUser(id int64) error
  ForceDeleteUser(id int64) error
  UpdateUser(id int64, user *entity.User) (entity.User, error)
  CountByRole(role entity.Role) (int64, error)

  FindUserById(id int64) (entity.User, error)
  FindUserByEmail(email string) (entity.User, error)
  FindUserByPhone(phone string) (entity.User, error)
}
