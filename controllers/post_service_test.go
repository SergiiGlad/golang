package controllers

import (
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
  "testing"
  "fmt"
  "github.com/aws/aws-sdk-go/service/s3/s3iface"
  "strings"
)

type mockDynamoDBClient struct {
  dynamodbiface.DynamoDBAPI
}

type mockS3Client struct{
  s3iface.S3API
}



func (m *mockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error){
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

func TestGetPost(t *testing.T) {
  var expectedResult Post
  expectedResult.UserID = "user"
  expectedResult.Title = "title"
  expectedResult.Text = "text"
  expectedResult.PostID = "post id"
  expectedResult.FileLink = "NULL"
  expectedResult.LastUpdate = "last update"
  expectedResult.Like = "0"

  svc := &mockDynamoDBClient{}
  id := "post id"
  gotPost, err := GetPost(svc, id)

  fmt.Println(gotPost)
  fmt.Println(expectedResult)

  if gotPost != expectedResult || err != nil {
    fmt.Println("ERROR")
  }
}

func TestCreateNewPost(t *testing.T) {
  var expectedResult Post
  expectedResult.UserID = "user"
  expectedResult.Title = "title"
  expectedResult.Text = "text"
  expectedResult.PostID = "post id"
  expectedResult.FileLink = "NULL"
  expectedResult.LastUpdate = "last update"
  expectedResult.Like = "0"

  newpost := Post{
    Text: "text",
    Title: "title",
    PostID: "post id",
    UserID: "user id",
    FileLink: "NULL",
    LastUpdate: "last update",
    Like: "0",
  }

  svc := &mockDynamoDBClient{}

  result, err := CreateNewPost(svc, newpost)

  if err == nil || result != expectedResult {
    fmt.Println("Error")
  }

}

func TestGetPostByUserID(t *testing.T) {
  var expectedResultArr []Post
  var expectedResult Post
  expectedResult.UserID = "user"
  expectedResult.Title = "title"
  expectedResult.Text = "text"
  expectedResult.PostID = "post id"
  expectedResult.FileLink = "NULL"
  expectedResult.LastUpdate = "last update"
  expectedResult.Like = "0"

  expectedResultArr = append(expectedResultArr, expectedResult)

  svc := &mockDynamoDBClient{}

  id := "user"

  result, err := GetPostByUserID(svc, id)

  for i:=0; i< len(expectedResultArr); i++{
    if result[i] != expectedResultArr[i] || err != nil {
      fmt.Println("Error")
    }
  }
}

func TestDeletePost(t *testing.T) {
  expectedResult := "post id"

  svcd := &mockDynamoDBClient{}
  svcs := &mockS3Client{}
  result := DeletePost(svcd, svcs, "post_id")

  if !strings.EqualFold(result, expectedResult) {
    fmt.Println("Error")
  }
}

func TestUpdatePost(t *testing.T) {
  var expectedResult Post
  expectedResult.UserID = "user"
  expectedResult.Title = "title"
  expectedResult.Text = "text"
  expectedResult.PostID = "post id"
  expectedResult.FileLink = "NULL"
  expectedResult.LastUpdate = "last update"
  expectedResult.Like = "0"

  newpost := Post{
    Text: "text",
    Title: "title",
    PostID: "post id",
    UserID: "user id",
    FileLink: "NULL",
    LastUpdate: "last update",
    Like: "0",
  }

  svc := &mockDynamoDBClient{}

  result, _ := UpdatePost(svc, newpost)

  if result != expectedResult {
    fmt.Println("Error")
  }
}

