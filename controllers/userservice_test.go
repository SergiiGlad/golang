package controllers

import (
  "go-team-room/models/dao"
  "errors"
  "testing"
  "go-team-room/models/dto"
  "go-team-room/models/dao/interfaces"
  "database/sql"
)

type mockDb struct {
  DB []dao.User
}

//user instance to be returned with errors
var errorUser dao.User

func (md mockDb) AddUser(user *dao.User) (int64, error) {
  md.DB = append(md.DB, *user)
  user.ID = int64(len(md.DB) - 1)
  return user.ID, nil
}

func (md mockDb) DeleteUser(id int64) error {
  if id < 0 {
    return errors.New("invalid id")
  }

  for indx, user := range md.DB {
    if user.ID == id {
      md.DB[indx].AccStatus = dao.Deleted
      return nil
    }
  }

  return errors.New("user could not be found")
}

func (md mockDb) UpdateUser(id int64, user *dao.User) error {
  if id < 0 {
    return errors.New("invalid id")
  }

  for indx, user := range md.DB {
    if user.ID == id {
      md.DB[indx] = user
      user.ID = id
      return nil
    }
  }

  return errors.New("user could not be found")
}

func (md mockDb) FindUserById(id int64) (dao.User, error) {

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

func (md mockDb) FindUserByEmail(email string) (dao.User, error) {

  for indx, user := range md.DB {
    if user.Email == email {
      return md.DB[indx], nil
    }
  }

  return errorUser, sql.ErrNoRows
}

func (md mockDb) FindUserByPhone(phone string) (dao.User, error) {

  for indx, user := range md.DB {
    if user.Phone == phone {
      return md.DB[indx], nil
    }
  }

  return errorUser, sql.ErrNoRows
}

func (md mockDb) FriendsByUserID(id int64) ([]int64, error) {
  if id >= int64(len(md.DB)) || id < 0 {
    return nil, errors.New("invalid id")
  }

  return []int64{}, nil
}

func (md mockDb) InsertPass(pass *dao.Password) (int64, error) {
  return 0, nil
}

func (md mockDb) LastPassByUserId(id int64) (dao.Password, error) {
  return dao.Password{}, nil
}

func (md mockDb) PasswdsByUserId(id int64) ([]dao.Password, error) {
  return []dao.Password{}, nil
}

var userService = UserService{}

func TestUserServiceCreate(t *testing.T) {
  tests := [] struct {
    description  string
    db           interfaces.MySqlDal
    newUser      dto.RequestUserDto
    expectReturn dto.ResponseUserDto
  }{
    {
      description: "CreateNewUser [Should perform successfully]",
      db: mockDb{[]dao.User{}},
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
        Friends: []int64{},
      },
    },
    {
      description: "CreateNewUser [Should return empty resp]",
      db: mockDb{[]dao.User{}},
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
      db: mockDb{[]dao.User{}},
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
      db:          mockDb{[]dao.User{}},
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
    userService.DB = tc.db

    respDto, _ := userService.CreateUser(&tc.newUser)

    if respDto.String() != tc.expectReturn.String() {
      t.Errorf("\nExpected: %s\nGot: %s", tc.expectReturn, respDto)
    }
  }
}

func TestUserServiceUpdate(t *testing.T) {
  tests := [] struct {
    description  string
    db           interfaces.MySqlDal
    newUser      dto.RequestUserDto
    expectReturn dto.ResponseUserDto
  }{
    {
      description: "Update user [Should perform successfully]",
      db: mockDb{[]dao.User{
        dao.User{
          ID:        0,
          Email:     "email@gmail.com",
          FirstName: "Name",
          LastName:  "surname",
          Phone:     "+380509684212",
        },
      },
      },
      newUser: dto.RequestUserDto{
        Email:     "newemail@gmail.com",
        FirstName: "Name",
        LastName:  "surname",
        Phone:     "+380509684211",
        Password:  "123456",
      },
      expectReturn: dto.ResponseUserDto{
        ID:        0,
        Email:     "newemail@gmail.com",
        FirstName: "Name",
        LastName:  "Surname",
        Phone:     "+380509684211",
        Friends: []int64{},
      },
    },
    {
      description: "Update user [Should return unique error]",
      db: mockDb{[]dao.User{
        dao.User{
          ID:        0,
          Email:     "email@gmail.com",
          FirstName: "Name",
          LastName:  "surname",
          Phone:     "+380509684212",
        },
      },
      },
      newUser: dto.RequestUserDto{
        Email:     "email@gmail.com",
        FirstName: "Name",
        LastName:  "surname",
        Phone:     "+380509684212",
        Password:  "123456",
      },
    },
  }

  for _, tc := range tests {
    userService.DB = tc.db

    respDto, _ := userService.UpdateUser(0, &tc.newUser)

    if respDto.String() != tc.expectReturn.String() {
      t.Errorf("\nExpected: %s\nGot: %s", tc.expectReturn, respDto)
    }
  }
}

func TestUserServiceDelete(t *testing.T) {
  tests := [] struct {
    description  string
    db           interfaces.MySqlDal
    newUser      dto.RequestUserDto
    expectReturn dto.ResponseUserDto
  }{
    {
      description: "Update user [Should perform successfully]",
      db: mockDb{[]dao.User{
        dao.User{
          ID:        0,
          Email:     "email@gmail.com",
          FirstName: "Name",
          LastName:  "surname",
          Phone:     "+380509684212",
        },
      },
      },
      newUser: dto.RequestUserDto{
        Email:     "newemail@gmail.com",
        FirstName: "Name",
        LastName:  "surname",
        Phone:     "+380509684211",
        Password:  "123456",
      },
      expectReturn: dto.ResponseUserDto{
        ID:        0,
        Email:     "newemail@gmail.com",
        FirstName: "Name",
        LastName:  "Surname",
        Phone:     "+380509684211",
        Friends: []int64{},
      },
    },
    {
      description: "Update user [Should return unique error]",
      db: mockDb{[]dao.User{
        dao.User{
          ID:        0,
          Email:     "email@gmail.com",
          FirstName: "Name",
          LastName:  "surname",
          Phone:     "+380509684212",
        },
      },
      },
      newUser: dto.RequestUserDto{
        Email:     "email@gmail.com",
        FirstName: "Name",
        LastName:  "surname",
        Phone:     "+380509684212",
        Password:  "123456",
      },
    },
  }

  for _, tc := range tests {
    userService.DB = tc.db

    respDto, _ := userService.UpdateUser(0, &tc.newUser)

    if respDto.String() != tc.expectReturn.String() {
      t.Errorf("\nExpected: %s\nGot: %s", tc.expectReturn, respDto)
    }
  }
}
