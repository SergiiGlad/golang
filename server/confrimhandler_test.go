package server

import (
  "github.com/gorilla/mux"
  "net/http"
  //"net/http/httptest"
  "github.com/pkg/errors"
  "testing"
  "fmt"
  "net/http/httptest"
)

func newGorilaConfrimHandlerMock(hf http.HandlerFunc) http.Handler {
  r := mux.NewRouter()
  r.HandleFunc("/confirm/email/{token}", hf).Methods("GET")
  return r
}

type TokenGeneratorMock struct {
}

func (tg TokenGeneratorMock) ApproveUser(token string) (bool, error) {
  if token == "badToken" {
    return false, errors.New("Badd token")
  }
  if token == "usedToken" {
    return false, nil
  }
  return true, nil
}

func (tg TokenGeneratorMock) GenerateTokenForEmail(email string) (string, error) {
  return "", nil
}

func TestConfirmAccount(t *testing.T) {
  tests := []struct {
    description        string
    handlerFunc        http.HandlerFunc
    expectedStatusCode int
    pathToken          string
  }{
    {
      description:
      "If token not used and no error occurred [Should return 200 OK]",
      handlerFunc:        ConfirmAccount(&TokenGeneratorMock{}),
      expectedStatusCode: http.StatusOK,
      pathToken:          "token",
    },
    {
      description:
      "If error occurred [Should return 400 Bad Request]",
      handlerFunc:        ConfirmAccount(&TokenGeneratorMock{}),
      expectedStatusCode: http.StatusBadRequest,
      pathToken:          "badToken",
    },
    {
      description:
      "If token was used [Should return 400  Bad Request]",
      handlerFunc:        ConfirmAccount(&TokenGeneratorMock{}),
      expectedStatusCode: http.StatusBadRequest,
      pathToken:          "usedToken",
    },
  }

  for _, tc := range tests {

    //method and path can have any valid values. We test handlers, not routers.
    req, err := http.NewRequest("GET", fmt.Sprintf("/confirm/email/%s", tc.pathToken), nil)
    if err != nil {
      t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := newGorilaConfrimHandlerMock(tc.handlerFunc)
    handler.ServeHTTP(rr, req)

    if respBody := rr.Body.String();
      rr.Code != tc.expectedStatusCode{
      t.Errorf("\nDecsription: %s\nExpected response code %v .\nGot code %v with body %s",
        tc.description, tc.expectedStatusCode, rr.Code, respBody)
    }
  }

}
