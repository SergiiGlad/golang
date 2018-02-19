package controllers

import (
  "errors"
  "testing"
  "go-team-room/models/dto"
  "go-team-room/models/dao/interfaces"
  "database/sql"
  "go-team-room/models/dao/entity"
)

type userDaoMock struct {
  DB []entity.User
}
type friendServiceMock struct {}

func (fs friendServiceMock) GetFriends(id int64) ([]dto.ShortUser, error) {
  if id < 1 {
    return nil, errors.New("Invalid id")
  }

  return []dto.ShortUser{}, nil
}

func (fs friendServiceMock) GetUsersWithRequests(id int64) ([]dto.ShortUser, error) {
  if id < 1 {
    return nil, errors.New("Invalid id")
  }

  return []dto.ShortUser{}, nil
}

func (fs friendServiceMock) GetFriendIds(id int64) ([]int64, error) {
  if id < 1 {
    return nil, errors.New("Invalid id")
  }

  return []int64{}, nil
}

func (fs friendServiceMock) NewFriendRequest(connection *entity.Connection) error {
  return nil
}

func (fs friendServiceMock) ApproveFriendRequest(connection *entity.Connection) error {
  return nil
}

func (fs friendServiceMock) RejectFriendRequest(connection *entity.Connection) error {
  return nil
}

func (fs friendServiceMock) DeleteFriendship(friendship *entity.Connection) error {
  return nil
}

var friendServiceMocked = friendServiceMock{}
//user instance to be returned with errors
var errorUser entity.User

func (md userDaoMock) AddUser(user *entity.User) (entity.User, error) {
  md.DB = append(md.DB, *user)
  user.ID = int64(len(md.DB) - 1)
  return *user, nil
}

func (md userDaoMock) DeleteUser(id int64) error {
  if id < 0 {
    return errors.New("invalid id")
  }

  for indx, user := range md.DB {
    if user.ID == id {
      md.DB[indx].AccStatus = entity.Deleted
      return nil
    }
  }

  return errors.New("user could not be found")
}

func (md userDaoMock) ForceDeleteUser(id int64) error {
  if id < 0 {
    return errors.New("invalid id")
  }

  return errors.New("user could not be found")
}

func (md userDaoMock) UpdateUser(id int64, user *entity.User) (entity.User, error) {
  if id < 0 {
    return *user, errors.New("invalid id")
  }

  for indx, user := range md.DB {
    if user.ID == id {
      md.DB[indx] = user
      user.ID = id
      return user, nil
    }
  }

  return *user, errors.New("user could not be found")
}

func (md userDaoMock) CountByRole(role entity.Role) (int64, error) {

  counter := 0

  for _, user := range md.DB {
    if user.Role == role {
      counter++
    }
  }

  return int64(counter), nil
}

func (md userDaoMock) FindUserById(id int64) (entity.User, error) {

  if id < 0 {
    return errorUser, errors.New("invalid id")
  }

  for indx, user := range md.DB {
    if user.ID == id {
      return md.DB[indx], nil
    }
  }

  return errorUser, sql.ErrNoRows
}

func (md userDaoMock) FindUserByEmail(email string) (entity.User, error) {

  for indx, user := range md.DB {
    if user.Email == email {
      return md.DB[indx], nil
    }
  }

  return errorUser, sql.ErrNoRows
}

func (md userDaoMock) FindUserByPhone(phone string) (entity.User, error) {

  for indx, user := range md.DB {
    if user.Phone == phone {
      return md.DB[indx], nil
    }
  }

  return errorUser, sql.ErrNoRows
}

func (md userDaoMock) FriendsByUserID(id int64) ([]int64, error) {
  if id >= int64(len(md.DB)) || id < 0 {
    return nil, errors.New("invalid id")
  }

  return []int64{}, nil
}

type passDaoMock struct {}

func (md passDaoMock) InsertPass(pass *entity.Password) (int64, error) {
  return 0, nil
}

func (md passDaoMock) LastPassByUserId(id int64) (entity.Password, error) {
  return entity.Password{}, nil
}

func (md passDaoMock) PasswdsByUserId(id int64) ([]entity.Password, error) {
  return []entity.Password{}, nil
}

var userService = UserService{friendServiceMocked, passDaoMock{}, userDaoMock{}}

func TestUserServiceCreate(t *testing.T) {
  tests := [] struct {
    description  string
    db           interfaces.UserDao
    newUser      dto.RequestUserDto
    expectReturn dto.ResponseUserDto
  }{
    {
      description: "CreateNewUser [Should perform successfully]",
      db:          userDaoMock{[]entity.User{}},
      newUser: dto.RequestUserDto{
        Email:     "email@gmail.com",
        FirstName: "Name",
        LastName:  "surname",
        Phone:     "+380509684212",
        Password:  "123456",
      },
      expectReturn: dto.ResponseUserDto{
        ID:        0,
        Email:     "email@gmail.com",
        FirstName: "Name",
        LastName:  "Surname",
        Phone:     "+380509684212",
        Friends: 0,
      },
    },
    {
      description: "CreateNewUser [Should return empty resp]",
      db:          userDaoMock{[]entity.User{}},
      newUser: dto.RequestUserDto{
        Email:     "email@gmail",
        FirstName: "Name",
        LastName:  "surname",
        Phone:     "+380509684212",
        Password:  "123456",
      },
    },
    {
      description: "CreateNewUser [Should return empty resp]",
      db:          userDaoMock{[]entity.User{}},
      newUser: dto.RequestUserDto{
        Email:     "email@gmail.com",
        FirstName: "name",
        LastName:  "Surname",
        Phone:     "+380509",
        Password:  "123456",
      },
    },
    {
      description: "CreateNewUser [Should return empty resp]",
      db:          userDaoMock{[]entity.User{}},
      newUser: dto.RequestUserDto{
        Email:     "email@gmail.com",
        FirstName: "name",
        LastName:  "surname",
        Phone:     "+380509684212",
        Password:  "1",
      },
    },
  }

  for _, tc := range tests {
    userService.UserDao = tc.db

    respDto, _ := userService.CreateUser(&tc.newUser)

    if respDto.String() != tc.expectReturn.String() {
      t.Errorf("\nExpected: %s\nGot: %s", tc.expectReturn, respDto)
    }
  }
}

func TestUserServiceUpdate(t *testing.T) {
  tests := [] struct {
    description  string
    db           interfaces.UserDao
    newUser      dto.RequestUserDto
    expectReturn dto.ResponseUserDto
  }{
    {
      description: "UpdateStatus user [Should perform successfully]",
      db: userDaoMock{[]entity.User{
        entity.User{
          ID:        0,
          Email:     "email@gmail.com",
          FirstName: "Name",
          LastName:  "surname",
          Phone:     "+380509684212",
          AvatarRef: "",
        },
      },
      },
      newUser: dto.RequestUserDto{
        Email:     "newemail@gmail.com",
        FirstName: "Name",
        LastName:  "surname",
        Phone:     "+380509684211",
        Password:  "123456",
        Avatar:    "",
      },
      expectReturn: dto.ResponseUserDto{
        ID:        0,
        Email:     "newemail@gmail.com",
        FirstName: "Name",
        LastName:  "Surname",
        Phone:     "+380509684211",
        Avatar:    "",
        Friends: 0,
      },
    },
    {
      description: "UpdateStatus user [Should return unique error]",
      db: userDaoMock{[]entity.User{
        entity.User{
          ID:        0,
          Email:     "email@gmail.com",
          FirstName: "Name",
          LastName:  "surname",
          Phone:     "+380509684212",
          AvatarRef: "",
        },
      },
      },
      newUser: dto.RequestUserDto{
        Email:     "email@gmail.com",
        FirstName: "Name",
        LastName:  "surname",
        Phone:     "+380509684212",
        Password:  "123456",
        Avatar:    "",
      },
    },
  }

  for _, tc := range tests {
    userService.UserDao = tc.db

    respDto, _ := userService.UpdateUser(0, &tc.newUser)

    if respDto.String() != tc.expectReturn.String() {
      t.Errorf("\nExpected: %s\nGot: %s", tc.expectReturn, respDto)
    }
  }
}

func TestUserServiceDelete(t *testing.T) {
  tests := [] struct {
    description  string
    db           interfaces.UserDao
    id           int64
    expectReturn dto.ResponseUserDto
  }{
    {
      description: "UpdateStatus user [Should perform successfully]",
      db: userDaoMock{[]entity.User{
        entity.User{
          ID:        0,
          Email:     "email@gmail.com",
          FirstName: "Name",
          LastName:  "surname",
          Phone:     "+380509684212",
          AvatarRef: "",
        },
      },
      },
      id: 0,
      expectReturn: dto.ResponseUserDto{
        ID:        0,
        Email:     "email@gmail.com",
        FirstName: "Name",
        LastName:  "surname",
        Phone:     "+380509684212",
        Avatar:    "",
        Friends: 0,
      },
    },
  }

  for _, tc := range tests {
    userService.UserDao = tc.db

    respDto, _ := userService.DeleteUser(0)

    if respDto.String() != tc.expectReturn.String() {
      t.Errorf("\nExpected: %s\nGot: %s", tc.expectReturn, respDto)
    }
  }
}
