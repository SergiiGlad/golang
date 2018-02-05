package dao

import (
  "fmt"
)

type UserToken struct {
  ID       int64
  Email    string
  Token    string
  IsActive bool
  UserId   int64
}

func (ut UserToken) String() string {
  return fmt.Sprintf("UserToken objetct: ID=%s, Email=%s, Token=%s, IsActive=%s, UserId=%s",
    ut.ID, ut.Email, ut.Token, ut.IsActive, ut.UserId)
}
