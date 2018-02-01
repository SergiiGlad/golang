package controllers

import (
  "go-team-room/models/dao"
  "errors"
  "testing"
  "go-team-room/models/dto"
)

type mockDb struct {
  Dao []dao.User
}

func (md mockDb) AddUser(user *dao.User) (int64, error) {
  md.Dao = append(md.Dao, *user)
  user.ID = int64(len(md.Dao) - 1)
  return user.ID, nil
}

func (md mockDb) DeleteUser(id int64) error {
  if id < 0 {
    return errors.New("invalid id")
  }

  for indx, user := range md.Dao {
    if user.ID == id {
      md.Dao[indx].AccStatus = dao.Deleted
      return nil
    }
  }

  return errors.New("user could not be found")
}

func (md mockDb) UpdateUser(id int64, user *dao.User) error {
  if id < 0 {
    return errors.New("invalid id")
  }

  for indx, user := range md.Dao {
    if user.ID == id {
      md.Dao[indx] = user
      user.ID = id
      return nil
    }
  }

  return errors.New("user could not be found")
}

func (md mockDb) FindUserById(id int64) (*dao.User, error) {

  if id < 0 {
    return nil, errors.New("invalid id")
  }

  for indx, user := range md.Dao {
    if user.ID == id {
      return &md.Dao[indx], nil
    }
  }

  return nil, errors.New("user could not be found")
}

func (md mockDb) FindUserByEmail(email string) (*dao.User, error) {

  for indx, user := range md.Dao {
    if user.Email == email {
      return &md.Dao[indx], nil
    }
  }

  return nil, errors.New("user could not be found")
}

func (md mockDb) FindUserByPhone(phone string) (*dao.User, error) {

  for indx, user := range md.Dao {
    if user.Phone == phone {
      return &md.Dao[indx], nil
    }
  }

  return nil, errors.New("user could not be found")
}

func (md mockDb) FriendsByUserID(id int64) ([]int64, error) {
  if id >= int64(len(md.Dao)) || id < 0 {
    return nil, errors.New("invalid id")
  }

  return []int64{}, nil
}

var userService = UserService{}

func TestUserServiceCreate(t *testing.T) {
  tests := [] struct {
    description  string
    db           mockDb
    newUser      dto.RequestUserDto
    expectReturn dto.ResponseUserDto
  }{
    {
      description: "CreateNewUser [Should perform successfully]",
      db: mockDb{[]dao.User{}},
      newUser: dto.RequestUserDto{
        Email: "email@gmail.com",
        FirstName: "Name",
        LastName: "surname",
        Phone: "+380509684212",
        CurrentPass: "123456",
      },
    },
  }

  for _, tc := range tests {
    userService.Dao = tc.db

    respDto, _ := userService.CreateUser(&tc.newUser)

    if respDto.String() != tc.expectReturn.String() {
      t.Errorf("\nExpected: %s\nGot:%s", tc.expectReturn, respDto)
    }
  }
}
