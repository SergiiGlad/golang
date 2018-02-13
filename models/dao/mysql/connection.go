package mysql

import (
	"go-team-room/conf"
	"go-team-room/models/dao/interfaces"
  "database/sql"
  "fmt"
  "github.com/go-sql-driver/mysql"
  "database/sql/driver"
)

// Get instance of logger (Formatter, Hook，Level, Output ).
// If you want to use only your log message  It will need use own call logs example
var log = conf.GetLog()

//this var defines global mysql db connection that can be used and accessed in any project part
var (
  Conn          *sql.DB
  UserDao       interfaces.UserDao
  PasswordDao   interfaces.PasswordDao
  FriendshipDao interfaces.FriendshipDao
  TokenDao      interfaces.UserTokenDao
)

func init() {

  var err error

  Conn, err = newMySqlConnection()
  if err != nil {
    log.Fatal(err)
  }

  UserDao, err = newMySqlUserDao(Conn)
  if err != nil {
    log.Fatal(err)
  }

  PasswordDao, err = newMySqlPassDao(Conn)
  if err != nil {
    log.Fatal(err)
  }

  FriendshipDao, err = newMySqlFriendshipDao(Conn)
  if err != nil {
    log.Fatal(err)
  }

  TokenDao, err = newMySqlTokenDao(Conn)
  if err != nil {
    log.Fatal(err)
  }

}

func newMySqlConnection() (*sql.DB, error) {

  // Check database and table exists. If not, create it.
  if err := ensureTablesExist(); err != nil {
    return nil, err
  }

  conn, err := sql.Open("mysql", conf.MysqlDsn + conf.MysqlDBName)

  if err != nil {
    return nil, fmt.Errorf("mysql: could not get a connection: %v", err)
  }

  if err := conn.Ping(); err != nil {
    conn.Close()
    return nil, fmt.Errorf("mysql: could not establish a good connection with database: %v", err)
  }

  return conn, nil
}

var createTableStatements = []string{
  fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`, conf.MysqlDBName),

  fmt.Sprintf(`USE %s;`, conf.MysqlDBName),

  `CREATE TABLE IF NOT EXISTS users_data (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    phone VARCHAR(20),
    role_in_network ENUM('admin', 'user') NOT NULL,
    account_status ENUM('inactive', 'active', 'deleted') NOT NULL,
    avatar_ref MEDIUMTEXT
  );`,

  `CREATE TABLE IF NOT EXISTS users_passwords (
    password_id SERIAL PRIMARY KEY,
    password VARCHAR(200) NOT NULL,
    password_created TIMESTAMP NOT NULL,
    user_id INTEGER REFERENCES users_data(user_id)
  );`,

  `CREATE TABLE IF NOT EXISTS user_tokens (
  token_id SERIAL PRIMARY KEY,
  token VARCHAR(128) NOT NULL,
  email VARCHAR(100) NOT NULL,
  is_active BOOLEAN,
  user_id INTEGER REFERENCES users_data(user_id)
  );`,

  `CREATE TABLE friend_list (
  friend_user_id INTEGER REFERENCES users_data(user_id),
  user_id INTEGER REFERENCES users_data(user_id),
  connection_status ENUM('approved', 'rejected', 'waiting') NOT NULL,
  user_id_equals_friend_id CHAR(0) AS (CASE WHEN friend_user_id NOT IN (user_id) THEN '' END) VIRTUAL NOT NULL
  );`,
}

func ensureTablesExist() error {
  conn, err := sql.Open("mysql", conf.MysqlDsn)

  if err != nil {
    return fmt.Errorf("mysql: could not get a connection: %v", err)
  }
  defer conn.Close()

  if conn.Ping() == driver.ErrBadConn {
    return fmt.Errorf("mysql: could not connect to the database. " +
      "could be bad address, or this address is not whitelisted for access.")
  }

  if  _, err := conn.Exec(fmt.Sprintf("USE %s", conf.MysqlDBName)); err != nil {
    // MySQL error 1049 is "database does not exist"
    if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1049 {
      return createAllTables(conn)
    }
  }

  if _, err := conn.Exec("DESCRIBE users_data"); err != nil {
    // MySQL error 1146 is "table does not exist"
    if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
      return fmt.Errorf("mysql: could not connect to the database table: users_data")
    }

    return fmt.Errorf("mysql: could not connect to the database: %v", err)
  }

  if _, err := conn.Exec("DESCRIBE users_passwords"); err != nil {
    // MySQL error 1146 is "table does not exist"
    if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
      return fmt.Errorf("mysql: could not connect to the database table: users_passwords")
    }

    return fmt.Errorf("mysql: could not connect to the database: %v", err)
  }

  if _, err := conn.Exec("DESCRIBE friend_list"); err != nil {
    // MySQL error 1146 is "table does not exist"
    if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
      return fmt.Errorf("mysql: could not connect to the database table: friend_list")
    }

    return fmt.Errorf("mysql: could not connect to the database: %v", err)
  }

  if _, err := conn.Exec("DESCRIBE user_tokens"); err != nil {
    // MySQL error 1146 is "table does not exist"
    if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
      return fmt.Errorf("mysql: could not connect to the database table: user_tokenst")
    }

    return fmt.Errorf("mysql: could not connect to the database: %v", err)
  }

  return nil
}

func createAllTables(conn *sql.DB) error {
  for _, stmt := range createTableStatements {
    _, err := conn.Exec(stmt)
    if err != nil {
      return err
    }
  }

  return nil
}

// rowScanner is implemented by sql.Row and sql.Rows
type rowScanner interface {
  Scan(dest ...interface{}) error
}

// execAffectingOneRow executes a given statement, expecting one row to be affected.
func execAffectingOneRow(stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
  r, err := stmt.Exec(args...)
  if err != nil {
    return r, fmt.Errorf("mysql: could not execute statement: %v", err)
  }
  rowsAffected, err := r.RowsAffected()
  if err != nil {
    return r, fmt.Errorf("mysql: could not get rows affected: %v", err)
  } else if rowsAffected != 1 {
    return r, fmt.Errorf("mysql: expected 1 row affected, got %d", rowsAffected)
  }
  return r, nil
}
