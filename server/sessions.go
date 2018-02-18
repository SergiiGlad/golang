package server

import (
  "github.com/gorilla/sessions"
  "net/http"
  "strconv"
  "go-team-room/models/Amazon"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws"
  "fmt"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "go-team-room/models/context"
)

var (
  store = sessions.NewCookieStore([]byte("secretkey"))
)

func Authorize(next http.Handler) http.Handler {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    // add session flashes when UI will be ready
    paths := []string{"/", "/login", "/logout", "/registration"} //paths no needed to check authorization

    for _, path := range paths {
      if r.URL.Path == path {
        next.ServeHTTP(w, r)
        return
      }
    }

    session, err := store.Get(r, "name")

    if err != nil {
      session.Options.MaxAge = -1
      session.Save(r, w)
      responseError(w, err, http.StatusForbidden)
      return
    }

    if session.Values["loginned"] == true && sessionIsValid(session){
      context.SetUserRoleToContext(r, session.Values["role"].(string))
      context.SetIdToContext(r, session.Values["user_id"].(int64))
      next.ServeHTTP(w, r)
      return
    } else {
      http.Error(w, "Forbidden", http.StatusForbidden)
      return
    }
  })
}

func sessionIsValid(s *sessions.Session) bool {
    res, err := Amazon.SVCD.GetItem(&dynamodb.GetItemInput{
    Key: map[string]*dynamodb.AttributeValue{
      "session_id": {
        S: aws.String(s.Values["session_id"].(string)),
      },
    },
    TableName: aws.String("UsersSessionsData"),
  })

  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case dynamodb.ErrCodeConditionalCheckFailedException:
        fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
      case dynamodb.ErrCodeProvisionedThroughputExceededException:
        fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
      case dynamodb.ErrCodeResourceNotFoundException:
        fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
      case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
        fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
      case dynamodb.ErrCodeInternalServerError:
        fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
      default:
        fmt.Println(aerr.Error())
      }
    } else {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      fmt.Println(err.Error())
    }
  }

  sessInst := res.Item

  if sessInst == nil {
    return false
  }

  //fmt.Println(*sessInst["session_id"].S, ",", *sessInst["user_id"].N)


  //if *sessInst["session_id"].S == s.Values["session_id"].(string) && *sessInst["user_id"].N == strconv.FormatInt(s.Values["user_id"].(int64), 10) {
  if *sessInst["user_id"].N == strconv.FormatInt(s.Values["user_id"].(int64), 10) {
    return true
  }
  return false
}
