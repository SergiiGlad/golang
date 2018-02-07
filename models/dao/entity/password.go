package entity

import "fmt"

type Password struct {
  ID        int64
  Password  string
  CreatedAt string
  UserId    int64
}

func (p *Password) String() string {
  return fmt.Sprintf("User Password:\n\tID = %s\n\tPassword = %s\n\tCreatedAt = %s\n\tUserId = %v\n",
    p.ID, p.Password, p.CreatedAt, p.UserId)
}
