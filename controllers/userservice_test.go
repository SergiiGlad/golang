package controllers
//
//import (
//  "errors"
//  "testing"
//  "go-team-room/models/dto"
//  "go-team-room/models/dao/interfaces"
//  "database/sql"
//  "go-team-room/models/dao/entity"
//)
//
//type mockDb struct {
//  DB []entity.User
//}
//
////user instance to be returned with errors
//var errorUser entity.User
//
//func (md mockDb) AddUser(user *entity.User) (entity.User, error) {
//  md.DB = append(md.DB, *user)
//  user.ID = int64(len(md.DB) - 1)
//  return *user, nil
//}
//
//func (md mockDb) DeleteUser(id int64) error {
//  if id < 0 {
//    return errors.New("invalid id")
//  }
//
//  for indx, user := range md.DB {
//    if user.ID == id {
//      md.DB[indx].AccStatus = entity.Deleted
//      return nil
//    }
//  }
//
//  return errors.New("user could not be found")
//}
//
//func (md mockDb) ForceDeleteUser(id int64) error {
//  if id < 0 {
//    return errors.New("invalid id")
//  }
//
//  return errors.New("user could not be found")
//}
//
//func (md mockDb) UpdateUser(id int64, user *entity.User) (entity.User, error) {
//  if id < 0 {
//    return *user, errors.New("invalid id")
//  }
//
//  for indx, user := range md.DB {
//    if user.ID == id {
//      md.DB[indx] = user
//      user.ID = id
//      return user, nil
//    }
//  }
//
//  return *user, errors.New("user could not be found")
//}
//
//func (md mockDb) CountByRole(role entity.Role) (int64, error) {
//
//  counter := 0
//
//  for _, user := range md.DB {
//    if user.Role == role {
//      counter++
//    }
//  }
//
//  return int64(counter), nil
//}
//
//func (md mockDb) FindUserById(id int64) (entity.User, error) {
//
//  if id < 0 {
//    return errorUser, errors.New("invalid id")
//  }
//
//  for indx, user := range md.DB {
//    if user.ID == id {
//      return md.DB[indx], nil
//    }
//  }
//
//  return errorUser, sql.ErrNoRows
//}
//
//func (md mockDb) FindUserByEmail(email string) (entity.User, error) {
//
//  for indx, user := range md.DB {
//    if user.Email == email {
//      return md.DB[indx], nil
//    }
//  }
//
//  return errorUser, sql.ErrNoRows
//}
//
//func (md mockDb) FindUserByPhone(phone string) (entity.User, error) {
//
//  for indx, user := range md.DB {
//    if user.Phone == phone {
//      return md.DB[indx], nil
//    }
//  }
//
//  return errorUser, sql.ErrNoRows
//}
//
//func (md mockDb) FriendsByUserID(id int64) ([]int64, error) {
//  if id >= int64(len(md.DB)) || id < 0 {
//    return nil, errors.New("invalid id")
//  }
//
//  return []int64{}, nil
//}
//
//func (md mockDb) InsertPass(pass *entity.Password) (int64, error) {
//  return 0, nil
//}
//
//func (md mockDb) LastPassByUserId(id int64) (entity.Password, error) {
//  return entity.Password{}, nil
//}
//
//func (md mockDb) PasswdsByUserId(id int64) ([]entity.Password, error) {
//  return []entity.Password{}, nil
//}
//
//var userService = UserService{}
//
//func TestUserServiceCreate(t *testing.T) {
//  tests := [] struct {
//    description  string
//    db           interfaces.MySqlDal
//    newUser      dto.RequestUserDto
//    expectReturn dto.ResponseUserDto
//  }{
//    {
//      description: "CreateNewUser [Should perform successfully]",
//      db: mockDb{[]entity.User{}},
//      newUser: dto.RequestUserDto{
//        Email:     "email@gmail.com",
//        FirstName: "Name",
//        LastName:  "surname",
//        Phone:     "+380509684212",
//        Password:  "123456",
//      },
//      expectReturn: dto.ResponseUserDto{
//        ID:        0,
//        Email:     "email@gmail.com",
//        FirstName: "Name",
//        LastName:  "Surname",
//        Phone:     "+380509684212",
//        Friends: []int64{},
//      },
//    },
//    {
//      description: "CreateNewUser [Should return empty resp]",
//      db: mockDb{[]entity.User{}},
//      newUser: dto.RequestUserDto{
//        Email:     "email@gmail",
//        FirstName: "Name",
//        LastName:  "surname",
//        Phone:     "+380509684212",
//        Password:  "123456",
//      },
//    },
//    {
//      description: "CreateNewUser [Should return empty resp]",
//      db: mockDb{[]entity.User{}},
//      newUser: dto.RequestUserDto{
//        Email:     "email@gmail.com",
//        FirstName: "name",
//        LastName:  "Surname",
//        Phone:     "+380509",
//        Password:  "123456",
//      },
//    },
//    {
//      description: "CreateNewUser [Should return empty resp]",
//      db:          mockDb{[]entity.User{}},
//      newUser: dto.RequestUserDto{
//        Email:     "email@gmail.com",
//        FirstName: "name",
//        LastName:  "surname",
//        Phone:     "+380509684212",
//        Password:  "1",
//      },
//    },
//  }
//
//  for _, tc := range tests {
//    userService.UserDao = tc.db
//
//    respDto, _ := userService.CreateUser(&tc.newUser)
//
//    if respDto.String() != tc.expectReturn.String() {
//      t.Errorf("\nExpected: %s\nGot: %s", tc.expectReturn, respDto)
//    }
//  }
//}
//
//func TestUserServiceUpdate(t *testing.T) {
//  tests := [] struct {
//    description  string
//    db           interfaces.MySqlDal
//    newUser      dto.RequestUserDto
//    expectReturn dto.ResponseUserDto
//  }{
//    {
//      description: "UpdateStatus user [Should perform successfully]",
//      db: mockDb{[]entity.User{
//        entity.User{
//          ID:        0,
//          Email:     "email@gmail.com",
//          FirstName: "Name",
//          LastName:  "surname",
//          Phone:     "+380509684212",
//        },
//      },
//      },
//      newUser: dto.RequestUserDto{
//        Email:     "newemail@gmail.com",
//        FirstName: "Name",
//        LastName:  "surname",
//        Phone:     "+380509684211",
//        Password:  "123456",
//      },
//      expectReturn: dto.ResponseUserDto{
//        ID:        0,
//        Email:     "newemail@gmail.com",
//        FirstName: "Name",
//        LastName:  "Surname",
//        Phone:     "+380509684211",
//        Friends: []int64{},
//      },
//    },
//    {
//      description: "UpdateStatus user [Should return unique error]",
//      db: mockDb{[]entity.User{
//        entity.User{
//          ID:        0,
//          Email:     "email@gmail.com",
//          FirstName: "Name",
//          LastName:  "surname",
//          Phone:     "+380509684212",
//        },
//      },
//      },
//      newUser: dto.RequestUserDto{
//        Email:     "email@gmail.com",
//        FirstName: "Name",
//        LastName:  "surname",
//        Phone:     "+380509684212",
//        Password:  "123456",
//      },
//    },
//  }
//
//  for _, tc := range tests {
//    userService.UserDao = tc.db
//
//    respDto, _ := userService.UpdateUser(0, &tc.newUser)
//
//    if respDto.String() != tc.expectReturn.String() {
//      t.Errorf("\nExpected: %s\nGot: %s", tc.expectReturn, respDto)
//    }
//  }
//}
//
//func TestUserServiceDelete(t *testing.T) {
//  tests := [] struct {
//    description  string
//    db           interfaces.MySqlDal
//    id           int64
//    expectReturn dto.ResponseUserDto
//  }{
//    {
//      description: "UpdateStatus user [Should perform successfully]",
//      db: mockDb{[]entity.User{
//        entity.User{
//          ID:        0,
//          Email:     "email@gmail.com",
//          FirstName: "Name",
//          LastName:  "surname",
//          Phone:     "+380509684212",
//        },
//      },
//      },
//      id: 0,
//      expectReturn: dto.ResponseUserDto{
//        ID:        0,
//        Email:     "email@gmail.com",
//        FirstName: "Name",
//        LastName:  "surname",
//        Phone:     "+380509684212",
//        Friends: []int64{},
//      },
//    },
//  }
//
//  for _, tc := range tests {
//    userService.UserDao = tc.db
//
//    respDto, _ := userService.DeleteUser(0)
//
//    if respDto.String() != tc.expectReturn.String() {
//      t.Errorf("\nExpected: %s\nGot: %s", tc.expectReturn, respDto)
//    }
//  }
//}
//
