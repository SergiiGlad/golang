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
    description      string
    userId           string
    expectedRespCode int
  } {
    {
      description:      "should perform successfully and return 200 OK code",
      userId:           "1",
      expectedRespCode: http.StatusOK,
    },
    {
      description:      "should perform unsuccessfully and return 400 BadRequest code because of illegal id",
      userId:           "0",
      expectedRespCode: http.StatusBadRequest,
    },
    {
      description:      "should perform unsuccessfully and return 404 NotFound code because of illegal id format",
      userId:           "00.100",
      expectedRespCode: http.StatusNotFound,
    },
  }

  for testNum, tc := range tests {
    req, err := http.NewRequest("GET", fmt.Sprintf("/profile/%s/friends", tc.userId), strings.NewReader(""))

    if err != nil {
      t.Errorf("Error while request initialization")
    }

    rr := httptest.NewRecorder()
    handler := newGorilaFriendServerMock(getFriends(fservice))
    handler.ServeHTTP(rr, req)

    if rr.Code != tc.expectedRespCode {
      t.Errorf("\nTEST CASE #%d\nExpected response with code %d. \nGot %d code",
        testNum + 1, tc.expectedRespCode, rr.Code)
    }
  }
}

func TestGetUsersWithRequests(t *testing.T) {
  tests := []struct {
    description      string
    userId           string
    expectedRespCode int
  } {
    {
      description:      "should perform successfully and return 200 OK code",
      userId:           "1",
      expectedRespCode: http.StatusOK,
    },
    {
      description:      "should perform unsuccessfully and return 400 BadRequest code because of illegal id",
      userId:           "0",
      expectedRespCode: http.StatusBadRequest,
    },
    {
      description:      "should perform unsuccessfully and return 404 NotFound code because of illegal id format",
      userId:           "text",
      expectedRespCode: http.StatusNotFound,
    },
  }

  for testNum, tc := range tests {
    req, err := http.NewRequest("GET", fmt.Sprintf("/profile/%s/friends/requests", tc.userId), strings.NewReader(""))

    if err != nil {
      t.Errorf("Error while request initialization")
    }

    rr := httptest.NewRecorder()
    handler := newGorilaFriendServerMock(getUsersWithRequests(fservice))
    handler.ServeHTTP(rr, req)

    if rr.Code != tc.expectedRespCode {
      t.Errorf("\nTEST CASE #%d\nExpected response with code %d. \nGot %d code",
        testNum + 1, tc.expectedRespCode, rr.Code)
    }
  }
}

func TestNewFriendRequest(t *testing.T) {
  tests := []struct {
    description      string
    requestBody      string
    expectedRespCode int
  } {
    {
      description:      "should perform successfully and return 200 OK code",
      requestBody:      `{
        "friend_user_id": 1,
        "user_id": 2,
        "connection_status": "waiting"
        }`,
      expectedRespCode: http.StatusOK,
    },
    {
      description:      "should perform unsuccessfully and return 400 BadRequest code because of empty req body",
      expectedRespCode: http.StatusBadRequest,
    },
  }

  for testNum, tc := range tests {
    req, err := http.NewRequest("POST", fmt.Sprintf("/friend"), strings.NewReader(tc.requestBody))

    if err != nil {
      t.Errorf("Error while request initialization")
    }

    rr := httptest.NewRecorder()
    handler := newGorilaFriendServerMock(newFriendRequest(fservice))
    handler.ServeHTTP(rr, req)

    if rr.Code != tc.expectedRespCode {
      t.Errorf("\nTEST CASE #%d\nExpected response with code %d. \nGot %d code",
        testNum + 1, tc.expectedRespCode, rr.Code)
    }
  }
}

func TestReplyToFriendRequest(t *testing.T) {
  tests := []struct {
    description      string
    requestBody      string
    expectedRespCode int
  } {
    {
      description:      "should perform successfully and return 200 OK code",
      requestBody:      `{
        "friend_user_id": 1,
        "user_id": 2,
        "connection_status": "approved"
        }`,
      expectedRespCode: http.StatusOK,
    },
    {
      description:      "should perform successfully and return 200 OK code",
      requestBody:      `{
        "friend_user_id": 1,
        "user_id": 2,
        "connection_status": "rejected"
        }`,
      expectedRespCode: http.StatusOK,
    },
    {
      description:      "should perform successfully and return 400 Bad Request code [Invalid request reply]",
      requestBody:      `{
        "friend_user_id": 1,
        "user_id": 2,
        "connection_status": ""
        }`,
      expectedRespCode: http.StatusBadRequest,
    },
    {
      description:      "should perform unsuccessfully and return 400 BadRequest code because of empty req body",
      expectedRespCode: http.StatusBadRequest,
    },
  }

  for testNum, tc := range tests {
    req, err := http.NewRequest("PUT", fmt.Sprintf("/friend"), strings.NewReader(tc.requestBody))

    if err != nil {
      t.Errorf("Error while request initialization")
    }

    rr := httptest.NewRecorder()
    handler := newGorilaFriendServerMock(replyToFriendRequest(fservice))
    handler.ServeHTTP(rr, req)

    if rr.Code != tc.expectedRespCode {
      t.Errorf("\nTEST CASE #%d\nExpected response with code %d. \nGot %d code",
        testNum + 1, tc.expectedRespCode, rr.Code)
    }
  }
}

func TestDeleteFriendship(t *testing.T) {
  tests := []struct {
    description      string
    requestBody      string
    expectedRespCode int
  } {
    {
      description:      "should perform successfully and return 200 OK code",
      requestBody:      `{
        "friend_user_id": 1,
        "user_id": 2,
        "connection_status": ""
        }`,
      expectedRespCode: http.StatusOK,
    },
    {
      description:      "should perform unsuccessfully and return 400 BadRequest code because of empty req body",
      expectedRespCode: http.StatusBadRequest,
    },
  }

  for testNum, tc := range tests {
    req, err := http.NewRequest("DELETE", fmt.Sprintf("/friend"), strings.NewReader(tc.requestBody))

    if err != nil {
      t.Errorf("Error while request initialization")
    }

    rr := httptest.NewRecorder()
    handler := newGorilaFriendServerMock(deleteFriendship(fservice))
    handler.ServeHTTP(rr, req)

    if rr.Code != tc.expectedRespCode {
      t.Errorf("\nTEST CASE #%d\nExpected response with code %d. \nGot %d code",
        testNum + 1, tc.expectedRespCode, rr.Code)
    }
  }
}
