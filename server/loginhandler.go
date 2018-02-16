package server

import (
  "net/http"
  "go-team-room/controllers"
  "io/ioutil"
  "go-team-room/models/dto"
  "encoding/json"
  //"github.com/gorilla/sessions"
  //"github.com/google/uuid"
  //"github.com/aws/aws-sdk-go/aws"
  //"github.com/aws/aws-sdk-go/aws/awserr"
  //"github.com/aws/aws-sdk-go/service/dynamodb"
  "fmt"
  //"go-team-room/models/Amazon"
  //"strconv"
  //"time"
)

func loginhandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)

  if err != nil {
    responseError(w, err, http.StatusBadRequest)
    return
  }

  var login dto.LoginDto

  err = json.Unmarshal(body, &login)

  if err != nil {
    responseError(w, err, http.StatusBadRequest)
    return
  }

  session, err := store.Get(r, "name")

  if err != nil {
    responseError(w, err, http.StatusBadRequest)
    return
  }

  user, err := controllers.Login(login.PhoneOrEmail, login.Password)

  if err != nil {
    responseError(w, err, http.StatusForbidden)
    //session.AddFlash(errors.New("Wrong credentials"), "_errors")
    //session.Save(r, w)
    return
  }

  var userRes dto.ResponseUserDto

  userRes = dto.UserEntityToResponseDto(user)
  userResMarsh, err := json.Marshal(userRes)

  if err != nil {
    responseError(w, err, http.StatusBadRequest)
    return
  }



  session.Values["loginned"] = true
  //session.Values["user_id"] = user.ID
  //session.Values["role"] = user.Role
  //session.Values["session_id"] = uuid.New().String()
  //
  //res, err := Amazon.SVCD.PutItem(&dynamodb.PutItemInput {
  //  Item: map[string]*dynamodb.AttributeValue {
  //    "session_id": {
  //      S: aws.String(session.Values["session_id"].(string)),
  //    },
  //    "user_id": {
  //      N: aws.String(strconv.FormatInt(session.Values["user_id"].(int64), 10)),
  //    },
  //    "TTL": {
  //      N: aws.String(strconv.Itoa(int(time.Now().Add(time.Hour * 24).Unix()))),
  //    },
  //  },
  //  TableName: aws.String("UsersSessionsData"),
  //  ReturnConsumedCapacity: aws.String("TOTAL"),
  //})
  //
  //if err != nil {
  //  if aerr, ok := err.(awserr.Error); ok {
  //    switch aerr.Code() {
  //    case dynamodb.ErrCodeConditionalCheckFailedException:
  //      fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
  //    case dynamodb.ErrCodeProvisionedThroughputExceededException:
  //      fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
  //    case dynamodb.ErrCodeResourceNotFoundException:
  //      fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
  //    case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
  //      fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
  //    case dynamodb.ErrCodeInternalServerError:
  //      fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
  //    default:
  //      fmt.Println(aerr.Error())
  //    }
  //  } else {
  //    // Print the error, cast err to awserr.Error to get the Code and
  //    // Message from an error.
  //    fmt.Println(err.Error())
  //  }
  //}
  //
  //fmt.Println(res)

  //store.Options = &sessions.Options {
  //  MaxAge:   24*60*60,
  //  Secure:   true,
  //  HttpOnly: true,
  //    }

  session.Save(r, w)
  fmt.Println(session.Values["loginned"])
  w.Write(userResMarsh)
}
