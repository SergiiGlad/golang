package server

import (
  "go-team-room/models/dto"
  "strings"
  "testing"
  "net/http"
  "net/http/httptest"
  "encoding/json"
  "io/ioutil"
  "github.com/pkg/errors"
)

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
        "first_name": "string",
        "last_name": "string",
        "phone": "string",
        "password": "0"
        }`,
      expectDto: dto.RequestUserDto{
        Email:     "string",
        FirstName: "string",
        LastName:  "string",
        Phone:     "string",
        Password:  "0",
      },
    },
    {
      description: "Should return empty struct",
      //because of non closing bracket in the end oj this json
      reqBody:  `{
        "email": "string",
        "first_name": "string",
        "last_name": "string",
        "phone": "string",
        "password": "0"`,
      expectDto: dto.RequestUserDto{},
    },
    {
      description: "Should return empty pass field",
      //because of invalid type of password field
      reqBody:  `{
        "email": "string",
        "first_name": "string",
        "last_name": "string",
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

//func TestDtoFromReq(t *testing.T) {
//
//  tests := []struct {
//    description  string
//    reqBody      string
//    inputDtoType interface{}
//    expectDto    interface{}
//  }{
//    {
//      description: "Should perform successfully",
//      reqBody:  `{
//        "email": "string@stdasdas.com",
//        "first_name": "string",
//        "last_name": "string",
//        "phone": "+380509787555",
//        "password": "string"
//      }`,
//      inputDtoType: dto.RequestUserDto{},
//      expectDto: dto.RequestUserDto{
//        Email:     "string@stdasdas.com",
//        FirstName: "string",
//        LastName:  "string",
//        Phone:     "+380509787555",
//        //Role:      entity.UserRole,
//        Password:  "string",
//      },
//    },
//    {
//      description: "Should return empty struct",
//      reqBody:  `{
//        "email": "string",
//        "first_name": "string",
//        "last_name": "string",
//        "phone": "string",
//        "password": "0"`,
//      inputDtoType: dto.RequestUserDto{},
//      expectDto:    dto.RequestUserDto{},
//    },
//    {
//      description: "Should return empty pass field",
//      reqBody:  `{
//        "email": "string",
//        "first_name": "string",
//        "last_name": "string",
//        "phone": "string",
//        "password": 0
//    }`,
//      inputDtoType: dto.RequestUserDto{},
//      expectDto: dto.RequestUserDto{
//        Email: "string",
//        FirstName: "string",
//        LastName: "string",
//        Phone: "string",
//      },
//    },
//  }
//
//  for index, tc := range tests {
//    req, err := http.NewRequest("GET", "any", strings.NewReader(tc.reqBody))
//
//    if err != nil {
//      t.Fatal(err)
//    }
//
//    dtoType := tc.inputDtoType
//    err = dtoFromReq(req, dtoType)
//
//    switch tc.expectDto.(type) {
//    case dto.RequestUserDto:
//      dtoType := dtoType.(dto.RequestUserDto)
//      dtoExpected := tc.expectDto.(dto.RequestUserDto)
//      if !strings.EqualFold(dtoType.String(), dtoExpected.String()) {
//        t.Errorf("\nExpected response inputDtoType \n%s\nGot\n%s",
//          tc.expectDto, dtoType)
//      }
//    case dto.ShortUser:
//      dtoType := dtoType.(dto.ShortUser)
//      dtoExpected := tc.expectDto.(dto.ShortUser)
//      if !strings.EqualFold(dtoType.String(), dtoExpected.String()) {
//        t.Errorf("\nExpected response inputDtoType \n%s\nGot\n%s",
//          tc.expectDto, dtoType)
//      }
//    default:
//      t.Errorf("\nTest case #%d.\nUnknown dto type. To make this test passed with your type you need to add" +
//        "type assertion in switch expression here. As an type assertion example you can use already " +
//          "added types.", index + 1)
//    }
//  }
//}

func mockHandlerWrapper(err error, code int) http.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request) {
    responseError(w, err, code)
  }
}

func TestResponseError(t *testing.T) {

  tests := []struct {
    description     string
    expectErrorText string
    expectedCode    int
    expectRespError dto.ResponseError
  }{
    {
      description: "responseError should be called",
      expectErrorText: "some error",
      expectedCode: http.StatusForbidden,
      expectRespError:  dto.ResponseError{Reason: "some error"},
    },
    {
      description: "responseError should be called",
      expectErrorText: "some error 2",
      expectedCode: http.StatusBadRequest,
      expectRespError:  dto.ResponseError{Reason: "some error 2"},
    },
  }

  for _, tc := range tests {
    req, err := http.NewRequest("GET", "/any", strings.NewReader(""))

    if err != nil {
      t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(mockHandlerWrapper(errors.New(tc.expectErrorText), tc.expectedCode))
    handler.ServeHTTP(rr, req)

    respError := dto.ResponseError{}
    body, err := ioutil.ReadAll(rr.Result().Body)
    if err != nil {
      t.Errorf("\ncannot read body from responce")
    }

    err = json.Unmarshal(body, &respError)
    if err != nil {
      t.Errorf("\ncannot unmarshal error from responce body")
    }

    if rr.Code != tc.expectedCode || !strings.EqualFold(respError.String(), tc.expectRespError.String()) {
      t.Errorf("\nDecsription: %s\nExpected response code %v with body %s.\nGot code %v with body %s",
        tc.description, tc.expectedCode, tc.expectRespError, rr.Code, respError)
    }
  }
}
