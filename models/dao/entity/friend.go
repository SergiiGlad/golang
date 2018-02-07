package entity

import (
  "database/sql/driver"
  "fmt"
)

//Friendship type represents user connectivity
type Friendship struct {
  FriendUserId int64
  UserId int64
  Status ConnectionStatus
}

type ConnectionStatus string

//Scan method for implementing Scanner interface. That allows custom types to be passed as Scanner type
//arguments.
func (r *ConnectionStatus) Scan(value interface{}) error {
  *r = ConnectionStatus(value.(string))
  return nil
}

//Value method for implementing Valuer interface. That allows custom types to be passed as Valuer type
//arguments
func (r ConnectionStatus) Value() (driver.Value, error) {
  return driver.Value(string(r)), nil
}

const (
  Approved ConnectionStatus = "approved"
  Rejected ConnectionStatus = "rejected"
  Waiting  ConnectionStatus = "waiting"
)

func (f Friendship) String() string {
  return fmt.Sprintf("Friendship:\n\tFriendUserId = %v\n\tUserId = %v\n\tStatus = %s",
    f.FriendUserId, f.UserId, f.Status)
}
