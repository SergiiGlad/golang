package server

import (
  "net/http"
  "go-team-room/controllers"
  "io/ioutil"
  "go-team-room/models/dto"
  "encoding/json"
  "github.com/gorilla/sessions"
  "fmt"
  "strconv"
  "github.com/google/uuid"
  "github.com/aws/aws-sdk-go/aws"
  "time"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "go-team-room/models/Amazon"
)

func loginhandler(w http.ResponseWriter, r *http.Request) {
  // get user struct from request body
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

  // check if user exists and his password validity
  user, err := controllers.Login(login.PhoneOrEmail, login.Password)

  if err != nil {
    responseError(w, err, http.StatusForbidden)
    //session.AddFlash(errors.New("Wrong credentials"), "_errors")
    //session.Save(r, w)
    return
  }

  // if everything ok try to get session from store or create new there
  session, err := store.Get(r, "name")
  // set session options
  session.Options = &sessions.Options {
    MaxAge:   24*60*60,
    HttpOnly: true,
  }

  if err != nil {
    responseError(w, err, http.StatusBadRequest)
    return
  }

  var userRes dto.ResponseUserDto

  userRes = dto.UserEntityToResponseDto(user)
  userResMarsh, err := json.Marshal(userRes)

  if err != nil {
    responseError(w, err, http.StatusBadRequest)
    return
  }

  // set session values to operate with them later
  session.Values["loginned"] = true
  session.Values["user_id"] = user.ID
  session.Values["role"] = string(user.Role)
  session.Values["session_id"] = uuid.New().String()

  res, err := Amazon.SVCD.PutItem(&dynamodb.PutItemInput {
    Item: map[string]*dynamodb.AttributeValue {
      "session_id": {
        S: aws.String(session.Values["session_id"].(string)),
      },
      "user_id": {
        N: aws.String(strconv.FormatInt(session.Values["user_id"].(int64), 10)),
      },
      "TTL": {
        N: aws.String(strconv.Itoa(int(time.Now().Add(time.Hour * 24).Unix()))),
      },
    },
    TableName: aws.String("UsersSessionsData"),
    ReturnConsumedCapacity: aws.String("TOTAL"),
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

  fmt.Println(res)

  session.Save(r, w)
  w.Write(userResMarsh)
}
