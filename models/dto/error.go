package dto

import "fmt"

type ResponseError struct {
  Reason string `json:"reason"`
}

func (rerror ResponseError) String() string {
  return fmt.Sprintf("ResponseError:\n\tReason: %s", rerror.Reason)
}
