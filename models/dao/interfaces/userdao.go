package interfaces

import "go-team-room/models/dao"

type UserDao interface {
  AddUser(user *dao.User) (dao.User, error)
  DeleteUser(id int64) error
  UpdateUser(id int64, user *dao.User) (dao.User, error)

  FindUserById(id int64) (dao.User, error)
  FindUserByEmail(email string) (dao.User, error)
  FindUserByPhone(phone string) (dao.User, error)
  FriendsByUserID(id int64) ([]int64, error)
  ForceDeleteUser(id int64) error
}
