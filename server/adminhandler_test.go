package server

import (
  "testing"
  "go-team-room/models/dto"
  "net/http"
  "strings"
  "net/http/httptest"
  "go-team-room/models/dao"
  "github.com/gorilla/mux"
  "errors"
)

type UserServiceMock struct {
}

func (usm UserServiceMock) CreateUser(userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {
  dao := dto.RequestUserDtoToEntity(userDto)
  resp := dto.UserEntityToResponseDto(&dao)
  resp.Friends = []int64{}
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
  respUser.Friends = []int64{}

  return respUser, nil
}


func (usm UserServiceMock) DeleteUser(id int64) (dto.ResponseUserDto, error) {
  respUser := dto.ResponseUserDto{}

  if id < 1 {
    return respUser, errors.New("negative id")
  }

  respUser = dto.UserEntityToResponseDto(&dao.User{})
  respUser.Friends = []int64{}
  return respUser, nil

  return respUser, errors.New("user could be found")
}

func (usm UserServiceMock) GetUserFriends(id int64) ([]int64, error) {

  if id < 0 {
    return []int64{}, errors.New("negative id")
  }

  return []int64{id}, nil
}


func TestUserDtoFromReq(t *testing.T) {

  tests := []struct {
    description string
    reqBody     string
    expectDto   dto.RequestUserDto
  }{
    {
      description: "Should perform successfully",
      reqBody:  `{
        "email": "string",
        "firstName": "string",
        "lastName": "string",
        "phone": "string",
        "password": "0"
        }`,
        expectDto: dto.RequestUserDto{
          Email: "string",
          FirstName: "string",
          LastName: "string",
          Phone: "string",
          CurrentPass: "0",
        },
    },
    {
      description: "Should return empty struct",
      reqBody:  `{
        "email": "string",
        "firstName": "string",
        "lastName": "string",
        "phone": "string",
        "password": "0"`,
      expectDto: dto.RequestUserDto{},
    },
    {
    description: "Should return empty pass field",
    reqBody:  `{
        "email": "string",
        "firstName": "string",
        "lastName": "string",
        "phone": "string",
        "password": 0
    }`,
      expectDto: dto.RequestUserDto{
        Email: "string",
        FirstName: "string",
        LastName: "string",
        Phone: "string",
      },
  },
  }

  for _, tc := range tests {
    req, err := http.NewRequest("GET", "any", strings.NewReader(tc.reqBody))

    if err != nil {
      t.Fatal(err)
    }

    dto, err := userDtoFromReq(req)

    if dto.String() != tc.expectDto.String() {
      t.Errorf("\nExpected response dto \n%s\nGot\n%s",
        tc.expectDto, dto)
    }
  }
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
        handlerFunc:         createProfile(UserServiceMock{}),
        expectedStatusCode: http.StatusOK,
        reqBody: `{
        "email": "string",
        "firstName": "string",
        "lastName": "string",
        "phone": "string",
        "password": "0"
        }`,
        expectRespBody:
          `{"id":0,"email":"string","firstName":"string","lastName":"string","phone":"string","friends":[]}`,

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

func newGorilaServerMock(hf http.HandlerFunc) http.Handler {
  r := mux.NewRouter()
  r.HandleFunc("/admin/profile/{user_id:[0-9]+}", hf).Methods("PUT", "DELETE")
  return r
}

func TestUpdateProfileHandler(t *testing.T) {
}

func TestDeleteProfileHandler(t *testing.T) {
  tests := []struct {
    description        string
    handlerFunc        http.HandlerFunc
    expectedStatusCode int
    expectRespBody     string
  }{
    {
      description:        "Deleting user [Should return 200 OK]",
      handlerFunc: deleteProfile(&UserServiceMock {}),
      expectedStatusCode: http.StatusOK,
      expectRespBody:
      `{"id":0,"email":"","firstName":"","lastName":"","phone":"","friends":[]}`,
    },
  }

  for _, tc := range tests {

    //method and path can have any valid values. We test handlers, not routers.
    req, err := http.NewRequest("DELETE", "/admin/profile/1", nil)

    if err != nil {
      t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := newGorilaServerMock(tc.handlerFunc)
    handler.ServeHTTP(rr, req)

    if respBody := rr.Body.String();
      rr.Code != tc.expectedStatusCode || respBody != tc.expectRespBody {
      t.Errorf("\nDecsription: %s\nExpected response code %v with body %s.\nGot code %v with body %s",
        tc.description, tc.expectedStatusCode, tc.expectRespBody, rr.Code, respBody)
    }
  }
}
