package server

import (
  "testing"
  "strings"
  "net/http/httptest"
  "net/http"
  "go-team-room/controllers"
  "github.com/gorilla/mux"
  "go-team-room/models/dto"
)

type mock struct {
  controllers.UserServiceInterface
}

func (m mock) GetUser(id int64) (dto.ResponseUserDto, error) {
  user := dto.ResponseUserDto{
    ID: 1,
    Email: "qwer@asdf.com",
    FirstName: "Denis",
    LastName: "Kononenko",
    Phone: "+380994321345",
    Avatar: "/uploads/16d6cbe5-6c14-4432-b167-dd94872ce31c.jpg",
    Friends: 0,
  }
  return user, nil
}

func TestGetProfile(t *testing.T) {
  mockSVC := &mock{}
  tests := []struct {
    description        string
    handlerFunc        http.HandlerFunc
    expectedStatusCode int
    reqBody            string
    expectRespBody     string
    url                string
  }{
    {
    description:        "GET Profile [Should return 200 OK]",
    handlerFunc:        GetProfile(mockSVC),
    expectedStatusCode: http.StatusOK,
    reqBody:            "",
    expectRespBody:
    `{"id":1,"email":"qwer@asdf.com","first_name":"Denis","last_name":"Kononenko","phone":"+380994321345","avatar_ref":"/uploads/16d6cbe5-6c14-4432-b167-dd94872ce31c.jpg","friends_num":0}`,
    url: "/profile/1",
    },
  }


  for _, tc := range tests {
    //method and path can have any valid values. We test handlers, not routers.
    req, err := http.NewRequest("GET", tc.url, strings.NewReader(tc.reqBody))

    if err != nil {
      t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := newProfileGorilaServerMock(tc.handlerFunc)
    handler.ServeHTTP(rr, req)

    if respBody := rr.Body.String();
      rr.Code != tc.expectedStatusCode || strings.EqualFold(respBody, tc.expectRespBody){
      t.Errorf("\nDecsription: %s\nExpected response code %v with body %s.\nGot code %v with body %s",
        tc.description, tc.expectedStatusCode, tc.expectRespBody, rr.Code, respBody)
    }
  }
}

func newProfileGorilaServerMock(hf http.HandlerFunc) http.Handler {
  r := mux.NewRouter()
  r.HandleFunc("/profile/{user_id}", hf).Methods("GET")
  return r
}
