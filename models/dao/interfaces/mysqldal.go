package interfaces

type MySqlDal interface {
  UserDao
  PasswordDao
  UserTokenDao
}
