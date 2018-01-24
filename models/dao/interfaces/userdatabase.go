package interfaces

import "go-team-room/models/dao"

type UserDatabase interface {
  AddUser(user *dao.User) (int64, error)
  DeleteUser(id int64) error
  UpdateUser(id int64, user *dao.User) error

  FindUserById(id int64) (*dao.User, error)
  FindUserByEmail(email string) (*dao.User, error)
  FindUserByPhone(phone string) (*dao.User, error)
}
