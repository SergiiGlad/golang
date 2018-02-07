package interfaces

import (
  "go-team-room/models/dao/entity"
)

//PasswordDao interface is used in services to follow dependency inversion principle
type PasswordDao interface {
  InsertPass(pass *entity.Password) (int64, error)
  LastPassByUserId(id int64) (entity.Password, error)
  PasswdsByUserId(id int64) ([]entity.Password, error)
}
