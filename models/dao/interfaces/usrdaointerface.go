package interfaces

import "go-team-room/models/dao"

type UserDao interface {
  Create(user *dao.User) error
  Delete(id int) error
  Update(id int, user *dao.User) error
  InitUsersTable() error

  FindById(id int) (dao.User, error)
  FindByEmail(email string) (dao.User, error)
  FindByPhone(phone string) (dao.User, error)
}
