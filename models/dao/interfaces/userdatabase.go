package interfaces

import "go-team-room/models/dao"

type UserDatabase interface {
  AddUser(user *dao.User) (int64, error)
  DeleteUser(id int) error
  UpdateUser(id int, user *dao.User) error

  FindUserById(id int) (*dao.User, error)
  FindUserByEmail(email string) (*dao.User, error)
  FindUserByPhone(phone string) (*dao.User, error)
}
