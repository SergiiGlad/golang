package server

import (
  "testing"
  "net/http"
  "strings"
  "net/http/httptest"
  "github.com/gorilla/mux"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockDynamoDBClient struct {
  dynamodbiface.DynamoDBAPI
}

//func TestGetPost(t *testing.T) {
//  mockSvc := &mockDynamoDBClient{}
//  tests := []struct {
//    description        string
//    handlerFunc        http.HandlerFunc
//    expectedStatusCode int
//    reqBody            string
//    expectRespBody     string
//  }{
//    {
//      description:        "GET Post [Should return 200 OK]",
//      handlerFunc:        GetPost(mockSvc),
//      expectedStatusCode: http.StatusOK,
//      reqBody: `{
//        "post_id": "2018-02-06 19:41:46.8453473 +0200 EET m=+182.704141501",
//        }`,
//      expectRespBody:
//      `{"post_title":"41","post_text":"rqew","post_id":"2018-02-06 19:41:46.8453473 +0200 EET m=+182.704141501","user_id":"awesome","post_like":"0","file_link":"NULL","post_last_update":"2018-02-06 19:41:46.8453473 +0200 EET m=+182.704141501"}`,
//
//    },
//  }
//
//  for _, tc := range tests {
//
//    //method and path can have any valid values. We test handlers, not routers.
//    req, err := http.NewRequest("GET", "/post/2018-02-06 19:41:46.8453473 +0200 EET m=+182.704141501", strings.NewReader(tc.reqBody))
//
//    if err != nil {
//      t.Fatal(err)
//    }
//
//    rr := httptest.NewRecorder()
//    handler := newPostGorilaServerMock(tc.handlerFunc)
//    handler.ServeHTTP(rr, req)
//
//    if respBody := rr.Body.String();
//      rr.Code != tc.expectedStatusCode || strings.EqualFold(respBody, tc.expectRespBody){
//      t.Errorf("\nDecsription: %s\nExpected response code %v with body %s.\nGot code %v with body %s",
//        tc.description, tc.expectedStatusCode, tc.expectRespBody, rr.Code, respBody)
//    }
//  }
//}

//func TestGetPostByUserID(t *testing.T) {
//  tests := []struct {
//    description        string
//    handlerFunc        http.HandlerFunc
//    expectedStatusCode int
//    reqBody            string
//    expectRespBody     string
//  }{
//    {
//      description:        "Get Post by USER_ID [Should return 200 OK]",
//      handlerFunc:        GetPostByUserID,
//      expectedStatusCode: http.StatusOK,
//      reqBody: `{}`,
//      expectRespBody:
//      `{"post_title":"41","post_text":"rqew","post_id":"2018-02-06 19:41:46.8453473 +0200 EET m=+182.704141501","user_id":"awesome","post_like":"0","file_link":"NULL","post_last_update":"2018-02-06 19:41:46.8453473 +0200 EET m=+182.704141501"}`,
//
//    },
//  }
//
//  for _, tc := range tests {
//
//    //method and path can have any valid values. We test handlers, not routers.
//    req, err := http.NewRequest("GET", "/post/user/USER", strings.NewReader(tc.reqBody))
//
//    if err != nil {
//      t.Fatal(err)
//    }
//
//    rr := httptest.NewRecorder()
//    handler := newPostGorilaServerMock(tc.handlerFunc)
//    handler.ServeHTTP(rr, req)
//
//    if respBody := rr.Body.String();
//      rr.Code != tc.expectedStatusCode || strings.EqualFold(respBody, tc.expectRespBody){
//      t.Errorf("\nDecsription: %s\nExpected response code %v with body %s.\nGot code %v with body %s",
//        tc.description, tc.expectedStatusCode, tc.expectRespBody, rr.Code, respBody)
//    }
//  }
//}

//func TestDeletePost(t *testing.T) {
//
//
//  tests := []struct {
//    description        string
//    handlerFunc        http.HandlerFunc
//    expectedStatusCode int
//    reqBody            string
//    expectRespBody     string
//  }{
//    {
//      description:        "Delete Post [Should return 200 OK]",
//      handlerFunc:        DeletePost,
//      expectedStatusCode: http.StatusOK,
//      reqBody: `{
//        "post_id": "2018-02-05 18:42:01.418318894 +0200 EET m=+275.696160409",
//        }`,
//      expectRespBody:
//      `{"post_id": "2018-02-05 18:42:01.418318894 +0200 EET m=+275.696160409"}`,
//    },
//  }
//
//  for _, tc := range tests {
//
//    //method and path can have any valid values. We test handlers, not routers.
//    req, err := http.NewRequest("DELETE", "/post/2018-02-05 18:42:01.418318894 +0200 EET m=+275.696160409", strings.NewReader(tc.reqBody))
//
//    if err != nil {
//      t.Fatal(err)
//    }
//
//    rr := httptest.NewRecorder()
//    handler := newPostGorilaServerMock(tc.handlerFunc)
//    handler.ServeHTTP(rr, req)
//
//    if respBody := rr.Body.String();
//      rr.Code != tc.expectedStatusCode || strings.EqualFold(respBody, tc.expectRespBody){
//      t.Errorf("\nDecsription: %s\nExpected response code %v with body %s.\nGot code %v with body %s",
//        tc.description, tc.expectedStatusCode, tc.expectRespBody, rr.Code, respBody)
//    }
//  }
//}

func newPostGorilaServerMock(hf http.HandlerFunc) http.Handler {
  r := mux.NewRouter()
  r.HandleFunc("/post/{post_id}", hf).Methods("PUT", "DELETE", "GET")
  r.HandleFunc("/post/user/{user_id}", hf).Methods("GET")
  return r
}

