package dto

import "fmt"

//ResponseError represents dto for error. It contains just one field Reason where error message located.
type ResponseError struct {
  Reason string `json:"reason"`
}

func (rerror ResponseError) String() string {
  return fmt.Sprintf("ResponseError:\n\tReason: %s", rerror.Reason)
}
