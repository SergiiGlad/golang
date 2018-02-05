//Package entities provides types for mapping with database tables.
package entity

import (
  "fmt"
  "database/sql/driver"
)

//User is entity that is used in data access layer operations for writing to and reading data from database.
type User struct {
  ID          int64
  Email       string
  FirstName   string
  LastName    string
  Phone       string
  Role        Role
  AccStatus   AccountStatus
  AvatarRef   string
}

//Role type is like enum for Role field in User type
type Role string

//Scan method for implementing Scanner interface. That allows custom types to be passed as Scanner type
//arguments.
func (r *Role) Scan(value interface{}) error {
  *r = Role(value.(string))
  return nil
}

//Value method for implementing Valuer interface. That allows custom types to be passed as Valuer type
//arguments
func (r Role) Value() (driver.Value, error) {
  return driver.Value(string(r)), nil
}

//Role type contains 2 types: Admin and User
const (
  AdminRole Role = "admin"
  UserRole  Role = "user"
)

//AccountStatus type is like enum for AccStatus field in User type
type AccountStatus string

//AccountStatus type contains 2 types: Active and Deleted
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
  return fmt.Sprintf("User object:\n\tID = %d\n\tEmail = %s\n\tFirstName = %s\n\tLastName = %s\n\tPhone = %s\n\tRole %s\n\tAccStatus = %s\n\tAvatarRef = %s\n",
    user.ID, user.FirstName, user.LastName, user.Email, user.Phone, user.Role, user.AccStatus, user.AvatarRef)
}
