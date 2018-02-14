package controllers

import (
  "go-team-room/models/dao/entity"
  "errors"
  "testing"
  "go-team-room/models/dto"
  "strings"
  "database/sql"
)

type friendDaoMock struct {}

func (fd friendDaoMock) InsertConnection(connection *entity.Connection) error {
  return nil
}

func (fd friendDaoMock) UpdateStatus(connection *entity.Connection) error {
  return nil
}

func (fd friendDaoMock) Delete(friendship *entity.Connection) error {
  return nil
}

func (fd friendDaoMock) FriendsByUserID(id int64) ([]entity.User, error) {
  if id < 1 {
    return nil, errors.New("Invalid id")
  }

  return []entity.User{}, nil
}

func (fd friendDaoMock) UsersWithRequestsTo(id int64) ([]entity.User, error) {
  if id < 1 {
    return nil, errors.New("Invalid id")
  }

  return []entity.User{}, nil
}

func (fd friendDaoMock) FindConnection(connection *entity.Connection) (entity.Connection, error) {
  return entity.Connection{}, sql.ErrNoRows
}

var friendService = FriendService{friendDaoMock{}, userDaoMock{}}

func TestGetFriends(t *testing.T) {
  tests := []struct{
    description    string
    userId         int64
    expectedResult []dto.ShortUser
  } {
    {
      description:    "should perform successfully",
      userId:         1,
      expectedResult: []dto.ShortUser{},
    },
  }

  for testNum, tc := range tests {
    testResult, err := friendService.GetFriends(tc.userId)
    if err != nil {
      t.Errorf("\nTEST CASE #%d error occured %s", testNum + 1, err.Error())
    }

    for index, shortUser := range testResult {
      if strings.EqualFold(shortUser.String(), tc.expectedResult[index].String()) {
        t.Errorf("\nTEST CASE #%d\nExpected result %v\nGot %v", tc.expectedResult, testResult)
      }
    }
  }
}

func TestGetUsersWithRequests(t *testing.T) {
  tests := []struct{
    description    string
    userId         int64
    expectedResult []dto.ShortUser
  } {
    {
      description:    "should perform successfully",
      userId:         1,
      expectedResult: []dto.ShortUser{},
    },
  }

  for testNum, tc := range tests {
    testResult, err := friendService.GetUsersWithRequests(tc.userId)
    if err != nil {
      t.Errorf("\nTEST CASE #%d error occured %s", testNum + 1, err.Error())
    }

    for index, shortUser := range testResult {
      if strings.EqualFold(shortUser.String(), tc.expectedResult[index].String()) {
        t.Errorf("\nTEST CASE #%d\nExpected result %v\nGot %v", tc.expectedResult, testResult)
      }
    }
  }
}

func TestGetFriendIds(t *testing.T) {
  tests := []struct{
    description    string
    userId         int64
    expectedResult []int64
  } {
    {
      description:    "should perform successfully",
      userId:         1,
      expectedResult: []int64{},
    },
  }

  for testNum, tc := range tests {
    testResult, err := friendService.GetFriendIds(tc.userId)
    if err != nil {
      t.Errorf("\nTEST CASE #%d error occured %s", testNum + 1, err.Error())
    }

    for index, userId := range testResult {
      if userId != tc.expectedResult[index] {
        t.Errorf("\nTEST CASE #%d\nExpected result %v\nGot %v", tc.expectedResult, testResult)
      }
    }
  }
}

func TestNewFriendRequest(t *testing.T) {
  tests := []struct{
    description    string
    input          *entity.Connection
    expectedResult error
  } {
    {
      description:    "should perform successfully",
      input:  &entity.Connection{1, 2, entity.Waiting},
      expectedResult: nil,
    },
  }

  for testNum, tc := range tests {
    testResult := friendService.NewFriendRequest(tc.input)
    if testResult != tc.expectedResult {
      t.Errorf("\nTEST CASE #%d\nExpected result: %v\nGot: %v", testNum + 1, tc.expectedResult, testResult)
    }
  }
}

func TestApproveFriendRequest(t *testing.T) {
  tests := []struct{
    description    string
    input          *entity.Connection
    expectedResult string
  } {
    {
      description:    "should perform returning return error",
      input:  &entity.Connection{1, 2, entity.Waiting},
      expectedResult: "Unable to approve non existing connection",
    },
  }

  for testNum, tc := range tests {
    testResult := friendService.ApproveFriendRequest(tc.input)
    if testResult.Error() != tc.expectedResult {
      t.Errorf("\nTEST CASE #%d\nExpected result: %v\nGot: %v", testNum + 1, tc.expectedResult, testResult)
    }
  }
}

func TestRejectFriendRequest(t *testing.T) {
  tests := []struct{
    description    string
    input          *entity.Connection
    expectedResult string
  } {
    {
      description:    "should perform returning return error",
      input:  &entity.Connection{1, 2, entity.Waiting},
      expectedResult: "Unable to reject non existing connection",
    },
  }

  for testNum, tc := range tests {
    testResult := friendService.RejectFriendRequest(tc.input)
    if testResult.Error() != tc.expectedResult {
      t.Errorf("\nTEST CASE #%d\nExpected result: %v\nGot: %v", testNum + 1, tc.expectedResult, testResult)
    }
  }
}

func TestDeleteFriendship(t *testing.T) {
  tests := []struct{
    description    string
    input          *entity.Connection
    expectedResult string
  } {
    {
      description:    "should perform returning return error",
      input:  &entity.Connection{1, 2, entity.Approved},
      expectedResult: "Unable to delete non existing friendship",
    },
  }

  for testNum, tc := range tests {
    testResult := friendService.DeleteFriendship(tc.input)
    if testResult.Error() != tc.expectedResult {
      t.Errorf("\nTEST CASE #%d\nExpected result: %v\nGot: %v", testNum + 1, tc.expectedResult, testResult)
    }
  }
}
