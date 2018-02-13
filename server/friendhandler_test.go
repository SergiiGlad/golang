package server

import (
  "go-team-room/models/dto"
  "go-team-room/models/dao/entity"
  "github.com/pkg/errors"
  "testing"
  "github.com/gorilla/mux"
  "net/http"
  "net/http/httptest"
  "fmt"
  "strings"
  "go-team-room/controllers"
)

type friendServiceMock struct {
}

func (fs friendServiceMock) GetFriends(id int64) ([]dto.ShortUser, error) {
  if id < 1 {
    return nil, errors.New("Invalid id")
  }

  return []dto.ShortUser{}, nil
}

func (fs friendServiceMock) GetUsersWithRequests(id int64) ([]dto.ShortUser, error) {
  if id < 1 {
    return nil, errors.New("Invalid id")
  }

  return []dto.ShortUser{}, nil
}

func (fs friendServiceMock) GetFriendIds(id int64) ([]int64, error) {
  if id < 1 {
    return nil, errors.New("Invalid id")
  }

  return []int64{}, nil
}

func (fs friendServiceMock) NewFriendRequest(connection *entity.Connection) error {
  return nil
}

func (fs friendServiceMock) ApproveFriendRequest(connection *entity.Connection) error {
  return nil
}

func (fs friendServiceMock) RejectFriendRequest(connection *entity.Connection) error {
  return nil
}

func (fs friendServiceMock) DeleteFriendship(friendship *entity.Connection) error {
  return nil
}

func newGorilaFriendServerMock(hf http.HandlerFunc) http.Handler {
  r := mux.NewRouter()
  r.HandleFunc("/profile/{user_id:[0-9]+}/friends", hf).Methods("GET")
  r.HandleFunc("/profile/{user_id:[0-9]+}/friends/requests", hf).Methods("GET")
  r.HandleFunc("/friend", hf).Methods("POST", "PUT", "DELETE")
  return r
}

var fservice controllers.FriendServiceInterface = friendServiceMock{}

func TestGetFriends(t *testing.T) {
  tests := []struct {
    description string
    userId int
    expectedCode int
  } {
    {
      description: "should perform successfully and return 200 OK code",
      userId: 1,
      expectedCode: http.StatusOK,
    },
  }

  for testNum, tc := range tests {
    req := httptest.NewRequest("GET", fmt.Sprintf("/profile/%s/friends", tc.userId), strings.NewReader(""))

    rr := httptest.NewRecorder()
    handler := newGorilaFriendServerMock(getFriends(fservice))
    handler.ServeHTTP(rr, req)

    if rr.Code
  }
}
