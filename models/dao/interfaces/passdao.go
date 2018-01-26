package interfaces

import "go-team-room/models/dao"

type PasswordDao interface {
  InsertPass(pass *dao.Password) (int64, error)
}
