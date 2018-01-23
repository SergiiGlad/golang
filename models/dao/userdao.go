package dao

import (
  "fmt"
  "database/sql/driver"
)

type User struct {
  ID          int
  Email       string
  FirstName   string
  SecondName  string
  Phone       string
  CurrentPass string
  Role        Role
  AccStatus   AccountStatus
  AvatarRef   string
}

type Role string

func (r *Role) Scan(value interface{}) error {
  *r = Role(value.(string))
  return nil
}

func (r Role) Value() (driver.Value, error) {
  return string(r), nil
}

const (
  AdminRole Role = "admin"
  UserRole  Role = "user"
)

type AccountStatus string

const (
  Active AccountStatus = "active"
  Deleted AccountStatus = "deleted"
)

func (a *AccountStatus) Scan(value interface{}) error {
  *a = AccountStatus(value.(string))
  return nil
}

func (a AccountStatus) Value() (driver.Value, error) {
  return string(a), nil
}

func (user User) String() string {
  return fmt.Sprintf("User object:\n\tID = %d\n\tEmail = %s\n\tFirstName = %s\n\tSecondName = %s\n\tPhone = %s\n\tCurrentPass = %s\n\tRole %s\n\tAccStatus = %s\n\tAvatarRef = %s\n",
    user.ID, user.FirstName, user.SecondName, user.Email, user.Phone, user.CurrentPass, user.Role, user.AccStatus, user.AvatarRef)
}
