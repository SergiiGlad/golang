package server

import (
  "testing"
  "go-team-room/models/dto"
  "net/http"
  "strings"
  "net/http/httptest"
  "github.com/gorilla/mux"
  "errors"
  "go-team-room/models/dao/entity"
  "fmt"
)

type UserServiceMock struct {
}

func (usm UserServiceMock) CreateUser(userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {
  dao := dto.RequestUserDtoToEntity(userDto)
  resp := dto.UserEntityToResponseDto(&dao)
  return resp, nil
}

func (usm UserServiceMock) UpdateUser(id int64, userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {
  respUser := dto.ResponseUserDto{}

  if id < 1 {
    return respUser, errors.New("negative id")
  }

  dao := dto.RequestUserDtoToEntity(userDto)
  dao.ID = id
  respUser = dto.UserEntityToResponseDto(&dao)

  return respUser, nil
}


func (usm UserServiceMock) DeleteUser(id int64) (dto.ResponseUserDto, error) {
  respUser := dto.ResponseUserDto{}

  if id < 1 {
    return respUser, errors.New("negative id")
  }

  respUser = dto.UserEntityToResponseDto(&entity.User{ID:id})
  return respUser, nil
}

func (usm UserServiceMock) GetUserFriends(id int64) ([]int64, error) {

  if id < 0 {
    return []int64{}, errors.New("negative id")
  }

  return []int64{id}, nil
}

func TestNewProfileHandler(t *testing.T) {
    tests := []struct {
      description        string
      handlerFunc        http.HandlerFunc
      expectedStatusCode int
      reqBody            string
      expectRespBody     string
    }{
      {
        description:        "Creating user [Should return 200 OK]",
        handlerFunc:        createProfileByAdmin(UserServiceMock{}),
        expectedStatusCode: http.StatusOK,
        reqBody: `{
        "email": "string",
        "first_name": "string",
        "last_name": "string",
        "phone": "string",
        "password": "0"
        }`,
        expectRespBody:
          `{"id":0,"email":"string","first_name":"string","last_name":"string","phone":"string","avatar_ref":"","friends_num":0}`,

      },
    }

    for _, tc := range tests {

      //method and path can have any valid values. We test handlers, not routers.
      req, err := http.NewRequest("GET", "/any", strings.NewReader(tc.reqBody))

      if err != nil {
        t.Fatal(err)
      }

      rr := httptest.NewRecorder()
      handler := http.HandlerFunc(tc.handlerFunc)
      handler.ServeHTTP(rr, req)

      if respBody := rr.Body.String();
        rr.Code != tc.expectedStatusCode || respBody != tc.expectRespBody {
        t.Errorf("\nDecsription: %s\nExpected response code %v with body %s.\nGot code %v with body %s",
          tc.description, tc.expectedStatusCode, tc.expectRespBody, rr.Code, respBody)
      }
    }
}

func newGorilaUserServerMock(hf http.HandlerFunc) http.Handler {
  r := mux.NewRouter()
  r.HandleFunc("/admin/profile/{user_id:[0-9]+}", hf).Methods("PUT", "DELETE")
  return r
}

func TestUpdateProfileHandler(t *testing.T) {
  tests := []struct {
    description        string
    handlerFunc        http.HandlerFunc
    expectedStatusCode int
    pathUserId         int
    reqBody            string
    expectRespBody     string
  }{
    {
      description:        "UpdateStatus user [Should return 200 OK]",
      handlerFunc:        updateProfileByAdmin(&UserServiceMock{}),
      expectedStatusCode: http.StatusOK,
      pathUserId:         1,
      reqBody: `{
        "email": "string",
        "first_name": "string",
        "last_name": "string",
        "phone": "string",
        "password": "0"
        }`,
      expectRespBody:
      `{"id":1,"email":"string","first_name":"string","last_name":"string","phone":"string","avatar_ref":"","friends_num":0}`,
    },
  }

  for _, tc := range tests {

    //method and path can have any valid values. We test handlers, not routers.
    req, err := http.NewRequest("PUT", fmt.Sprintf("/admin/profile/%d", tc.pathUserId), strings.NewReader(tc.reqBody))
    if err != nil {
      t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := newGorilaUserServerMock(tc.handlerFunc)
    handler.ServeHTTP(rr, req)

    if respBody := rr.Body.String();
      rr.Code != tc.expectedStatusCode || respBody != tc.expectRespBody {
      t.Errorf("\nDecsription: %s\nExpected response code %v with body %s.\nGot code %v with body %s",
        tc.description, tc.expectedStatusCode, tc.expectRespBody, rr.Code, respBody)
    }
  }
}

func TestDeleteProfileHandler(t *testing.T) {
  tests := []struct {
    description        string
    handlerFunc        http.HandlerFunc
    expectedStatusCode int
    pathUserId         int
    expectRespBody     string
  }{
    {
      description:        "Deleting user [Should return 200 OK]",
      handlerFunc:        deleteProfileByAdmin(&UserServiceMock {}),
      expectedStatusCode: http.StatusOK,
      pathUserId:         1,
      expectRespBody:
      `{"id":1,"email":"","first_name":"","last_name":"","phone":"","avatar_ref":"","friends_num":0}`,
    },
  }

  for _, tc := range tests {

    //method and path can have any valid values. We test handlers, not routers.
    req, err := http.NewRequest("DELETE", fmt.Sprintf("/admin/profile/%d", tc.pathUserId), nil)

    if err != nil {
      t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := newGorilaUserServerMock(tc.handlerFunc)
    handler.ServeHTTP(rr, req)

    if respBody := rr.Body.String();
      rr.Code != tc.expectedStatusCode || respBody != tc.expectRespBody {
      t.Errorf("\nDecsription: %s\nExpected response code %v with body %s.\nGot code %v with body %s",
        tc.description, tc.expectedStatusCode, tc.expectRespBody, rr.Code, respBody)
    }
  }
}
