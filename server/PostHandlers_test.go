package server

import (
  "testing"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "net/http/httptest"
  "strings"
)

type mockDynamoDBClient struct {
  dynamodbiface.DynamoDBAPI
}

func (m *mockDynamoDBClient) GetItem (input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error){
  test := make(map[string]*dynamodb.AttributeValue)
  s1 := "title"
  s2 := "text"
  s3 := "user"
  s4 := "post id"
  s5 := "NULL"
  s6 := "last update"
  s7 := "0"

  test["post_title"] = &dynamodb.AttributeValue{S: &s1}
  test["post_text"] = &dynamodb.AttributeValue{S: &s2}
  test["user_id"] = &dynamodb.AttributeValue{S: &s3}
  test["post_id"] = &dynamodb.AttributeValue{S: &s4}
  test["file_link"] = &dynamodb.AttributeValue{S: &s5}
  test["post_last_update"] = &dynamodb.AttributeValue{S: &s6}
  test["post_like"] = &dynamodb.AttributeValue{S: &s7}

  output := dynamodb.GetItemOutput{
    Item: test,
  }
  return &output, nil
}

func TestGetPostHandler(t *testing.T) {

  mockSVC := &mockDynamoDBClient{}
  tests := []struct {
    description        string
    handlerFunc        http.HandlerFunc
    expectedStatusCode int
    reqBody            string
    expectRespBody     string
    url                string
  }{
    {
      description:        "GET Post [Should return 200 OK]",
      handlerFunc:        GetPost(mockSVC),
      expectedStatusCode: http.StatusOK,
      reqBody: "",
      expectRespBody:
      `{"post_title":"title","post_text":"text","post_id":"post id","user_id":"user","post_like":"0","file_link":"NULL","post_last_update":"last update"}`,
      url: "/post/post id",
    },
  }

  for _, tc := range tests {

    //method and path can have any valid values. We test handlers, not routers.
    req, err := http.NewRequest("GET", tc.url, strings.NewReader(tc.reqBody))

    if err != nil {
      t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := newPostGorilaServerMock(tc.handlerFunc)
    handler.ServeHTTP(rr, req)

    if respBody := rr.Body.String();
      rr.Code != tc.expectedStatusCode || strings.EqualFold(respBody, tc.expectRespBody){
      t.Errorf("\nDecsription: %s\nExpected response code %v with body %s.\nGot code %v with body %s",
        tc.description, tc.expectedStatusCode, tc.expectRespBody, rr.Code, respBody)
    }
  }
}

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

