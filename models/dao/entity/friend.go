package entity

import (
  "database/sql/driver"
  "fmt"
)

//Connection type represents user connectivity
type Connection struct {
  FriendUserId int64            `json:"friend_user_id"`
  UserId       int64            `json:"user_id"`
  Status       ConnectionStatus `json:"connection_status"`
}

type ConnectionStatus string

const (
  Approved ConnectionStatus = "approved"
  Rejected ConnectionStatus = "rejected"
  Waiting  ConnectionStatus = "waiting"
)

//Scan method for implementing Scanner interface. That allows custom types to be passed as Scanner type
//arguments.
func (cs *ConnectionStatus) Scan(value interface{}) error {
  *cs = ConnectionStatus(value.(string))
  return nil
}

//Value method for implementing Valuer interface. That allows custom types to be passed as Valuer type
//arguments
func (cs ConnectionStatus) Value() (driver.Value, error) {
  return string(cs), nil
}

func (f Connection) String() string {
  return fmt.Sprintf("Connection:\n\tFriendUserId = %v\n\tUserId = %v\n\tStatus = %s",
    f.FriendUserId, f.UserId, f.Status)
}
