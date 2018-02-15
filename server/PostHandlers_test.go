package server

import (
  "testing"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
  "net/http/httptest"
  "strings"
  "github.com/aws/aws-sdk-go/service/s3/s3iface"
  "github.com/aws/aws-sdk-go/service/dynamodb"
)

type mockDynamoDBClient struct {
  dynamodbiface.DynamoDBAPI
}

type mockS3Client struct {
  s3iface.S3API
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
func (m *mockDynamoDBClient) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error){
  test := make([]map[string]*dynamodb.AttributeValue, 1)
  temp := make(map[string]*dynamodb.AttributeValue)
  s1 := "title"
  s2 := "text"
  s3 := "user"
  s4 := "post id"
  s5 := "NULL"
  s6 := "last update"
  s7 := "0"

  temp["post_title"] = &dynamodb.AttributeValue{S: &s1}
  temp["post_text"] = &dynamodb.AttributeValue{S: &s2}
  temp["user_id"] = &dynamodb.AttributeValue{S: &s3}
  temp["post_id"] = &dynamodb.AttributeValue{S: &s4}
  temp["file_link"] = &dynamodb.AttributeValue{S: &s5}
  temp["post_last_update"] = &dynamodb.AttributeValue{S: &s6}
  temp["post_like"] = &dynamodb.AttributeValue{S: &s7}

  test = append(test, temp)
  output := dynamodb.ScanOutput{
    Items: test,
  }
  return &output, nil
}
func (m *mockDynamoDBClient) DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error){
  output := dynamodb.DeleteItemOutput{}
  return &output, nil
}
func (m *mockDynamoDBClient) UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error){
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

  output := dynamodb.UpdateItemOutput{
    Attributes: test,
  }
  return &output, nil
}
func (m *mockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error){
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

  output := dynamodb.PutItemOutput{
    Attributes: test,
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

func TestGetPostByUserID(t *testing.T) {
  mocksvc := &mockDynamoDBClient{}
  tests := []struct {
    description        string
    handlerFunc        http.HandlerFunc
    expectedStatusCode int
    reqBody            string
    expectRespBody     string
    url                string
  }{
    {
      description:        "Get Post by USER_ID [Should return 200 OK]",
      handlerFunc:        GetPostByUserID(mocksvc),
      expectedStatusCode: http.StatusOK,
      reqBody: `{}`,
      expectRespBody:
      `{"post_title":"title","post_text":"text","post_id":"post id","user_id":"user","post_like":"0","file_link":"NULL","post_last_update":"last update"}`,
    },
  }

  for _, tc := range tests {

    //method and path can have any valid values. We test handlers, not routers.
    req, err := http.NewRequest("GET", "/post/user/user", strings.NewReader(tc.reqBody))

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

func TestDeletePost(t *testing.T) {
  mocksvcd := &mockDynamoDBClient{}
  mocksvcs := &mockS3Client{}

  tests := []struct {
    description        string
    handlerFunc        http.HandlerFunc
    expectedStatusCode int
    reqBody            string
    expectRespBody     string
  }{
    {
      description:        "Delete Post [Should return 200 OK]",
      handlerFunc:        DeletePost(mocksvcd, mocksvcs),
      expectedStatusCode: http.StatusOK,
      reqBody: `{
        "post_id": "post id",
        }`,
      expectRespBody:
      `{"post_id": "post id"}`,
    },
  }

  for _, tc := range tests {

    //method and path can have any valid values. We test handlers, not routers.
    req, err := http.NewRequest("DELETE", "/post/post id", strings.NewReader(tc.reqBody))

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

func TestUpdatePost(t *testing.T) {
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
      handlerFunc:        UpdatePost(mockSVC),
      expectedStatusCode: http.StatusOK,
      reqBody: "",
      expectRespBody:
      `{"post_title":"title","post_text":"text","post_id":"post id","user_id":"user","post_like":"0","file_link":"NULL","post_last_update":"last update"}`,
      url: "/post/post id",
    },
  }

  for _, tc := range tests {

    //method and path can have any valid values. We test handlers, not routers.
    req, err := http.NewRequest("PUT", tc.url, strings.NewReader(tc.reqBody))

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

//func TestCreateNewPost(t *testing.T) {
//  mockSVC := &mockDynamoDBClient{}
//  mockSVCS := &mockS3Client{}
//  tests := []struct {
//    description        string
//    handlerFunc        http.HandlerFunc
//    expectedStatusCode int
//    reqBody            string
//    expectRespBody     string
//    url                string
//  }{
//    {
//      description:        "Create Post [Should return 200 OK]",
//      handlerFunc:        CreateNewPost(mockSVC, mockSVCS),
//      expectedStatusCode: http.StatusOK,
//      reqBody: "post_title=title, ",
//      expectRespBody:
//      `{"post_title":"title","post_text":"text","post_id":"post id","user_id":"user","post_like":"0","file_link":"NULL","post_last_update":"last update"}`,
//      url: "/post/post id",
//    },
//  }
//
//  for _, tc := range tests {
//
//    //method and path can have any valid values. We test handlers, not routers.
//    req, err := http.NewRequest("POST", tc.url, strings.NewReader(tc.reqBody))
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

