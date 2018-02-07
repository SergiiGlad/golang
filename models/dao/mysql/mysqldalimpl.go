package mysql

import (
  "database/sql"
  "go-team-room/models/dao/interfaces"
  "go-team-room/conf"
  "fmt"
  //"database/sql/driver"
  "github.com/go-sql-driver/mysql"
  "database/sql/driver"
)

type MySqlDatabase struct {
  conn *sql.DB

  interfaces.UserDao
  interfaces.PasswordDao
}

func newMySQLDatabase() (MySqlDatabase, error) {

  db := MySqlDatabase{}

  // Check database and table exists. If not, create it.
  if err := ensureTableExists(); err != nil {
    return db, err
  }

  conn, err := sql.Open("mysql", conf.MysqlDsn + conf.MysqlDBName)

  if err != nil {
    return db, fmt.Errorf("mysql: could not get a connection: %v", err)
  }

  if err := conn.Ping(); err != nil {
    conn.Close()
    return db, fmt.Errorf("mysql: could not establish a good connection: %v", err)
  }

  userDao, err := newMySqlUserDao(conn)

  if err != nil {
    fmt.Errorf("mysql: could not establish connection with userDao: %s", err)
    return db, err
  }

  passwordDao, err := newMySqlPassDao(conn)

  if err != nil {
    fmt.Errorf("mysql: could not establish connection with userDao: %s", err)
    return db, err
  }

  db = MySqlDatabase{
    conn,
    userDao,
    passwordDao,
  }

  return db, nil
}

// Close closes the database, freeing up any resources.
func (db *MySqlDatabase) Close() {
  db.conn.Close()
}

var createTableStatements = []string{
  fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`, conf.MysqlDBName),

  fmt.Sprintf(`USE %s;`, conf.MysqlDBName),

  `CREATE TABLE IF NOT EXISTS users_data (
    user_id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    role_in_network ENUM('admin', 'user') NOT NULL,
    account_status ENUM('active', 'deleted') NOT NULL,
    avatar_ref MEDIUMTEXT
  );`,

  `CREATE TABLE IF NOT EXISTS users_passwords (
    password_id SERIAL PRIMARY KEY,
    password VARCHAR(200) NOT NULL,
    password_created TIMESTAMP NOT NULL,
    user_id INTEGER REFERENCES users_data(user_id)
  );`,

  `CREATE TABLE IF NOT EXISTS friend_list (
    friend_id INTEGER REFERENCES users_data(user_id),
    user_id INTEGER REFERENCES users_data(user_id),
    user_id_equals_friend_id CHAR(0) AS (CASE WHEN friend_id NOT IN (user_id) THEN '' END) VIRTUAL NOT NULL
  );`,
}

func ensureTableExists() error {
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
      return createAllTables(conn)
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
