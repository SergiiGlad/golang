package interfaces

//MySqlDal interface is general interface for mysql data access layer. It collects all dao methods so you
//can implement dao interfaces separately but use their concrete types together with one database connection.
type MySqlDal interface {
  UserDao
  PasswordDao
  UserTokenDao
}
