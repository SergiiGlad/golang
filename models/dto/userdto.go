package dto

import "fmt"

type UserDto struct {
  ID          int     `json:"id"`
  Email       string  `json:"email"`
  FirstName   string  `json:"firstName"`
  SecondName  string  `json:"secondName"`
  Phone       string  `json:"phone"`
  Friends     []int   `json:"friends"`
}

func (user UserDto) String() string {
  return fmt.Sprintf("User object:\n\tID = %d\n\tEmail = %s\n\tFirstName = %s\n\tSecondName = %s\n\tPhone = %s\n\tFriends = %v\n",
    user.ID, user.FirstName, user.SecondName, user.Email, user.Phone, user.Friends)
}
