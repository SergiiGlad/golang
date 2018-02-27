package server

import (
	"fmt"
	"go-team-room/models/Amazon"
	"go-team-room/models/context"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/sessions"
)

var (
	store = sessions.NewCookieStore([]byte("secretkey"))
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// add session flashes when UI will be ready
		paths := []string{"/", "/api/login", "/api/logout", "/api/registration"} //paths no needed to check authorization

		for _, path := range paths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r)
				return
			}
		}

		session, err := store.Get(r, "name")

		if err != nil {
			// sometimes we can use expired session to login, in this case we will get an error
			// we should catch this error by this block and reset session. After this we should reload the page
			session.Options.MaxAge = -1
			session.Save(r, w)
			responseError(w, err, http.StatusForbidden)
			return
		}

		if session.Values["loginned"] == true && sessionIsValid(session) {
			// set context values to use them on the next wrappers (check user role and user_id)
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

// Check session for its validity:
// if session_id exists in table and belongs to definite user_id its ok
func sessionIsValid(s *sessions.Session) bool {
	res, err := Amazon.SVCD.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"session_id": {
				S: aws.String(s.Values["session_id"].(string)),
			},
		},
		// To store sessions in DynamoDB we should have a table there with "UsersSessionsData" name
		// with primary key "session_id" and enabled TTL option with name "TTL"
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

	// check an item for existence in DB
	if sessInst == nil {
		return false
	}

	if *sessInst["user_id"].N == strconv.FormatInt(s.Values["user_id"].(int64), 10) {
		return true
	}
	return false
}
