package mysql

import (
  "go-team-room/models/dao"
  "log"
  "fmt"
)

type UserDaoImpl struct {
}

func (dao UserDaoImpl) Create(user *dao.User) error {
  query := "INSERT INTO " +
    "users_data (email, first_name, second_name, phone, current_password, role_in_network, account_status, avatar_ref) " +
    "VALUE (?, ?, ?, ?, ?, ?, ?, ?)"

  db := get()
  defer db.Close()

  statement, err := db.Prepare(query)

  if err != nil {
      log.Fatal(err)
      return err
  }

  defer statement.Close()

  result, err := statement.Exec(user.Email, user.FirstName, user.SecondName, user.Phone, user.CurrentPass, user.Role, user.AccStatus, user.AvatarRef)

  if err != nil {
    log.Fatal(err)
    return err
  }

  id, err := result.LastInsertId()

  if err != nil {
    log.Fatal(err)
    return err
  }

  user.ID = int(id)

  return nil
}

func (dao UserDaoImpl) Delete(id int) error {
  query := "UPDATE users_data SET account_status = 'deleted' WHERE user_id = ?"

  db := get()
  defer db.Close()

  statement, err := db.Prepare(query)

  if err != nil {
    log.Fatal(err)
    return err
  }

  defer statement.Close()

  _, err = statement.Exec(id)

  if err != nil {
    log.Fatal(err)
    return err
  }

  return nil
}

func (dao UserDaoImpl) Update(id int, user *dao.User) error {
  query := "UPDATE users_data SET " +
    "email = ?, first_name = ?, second_name = ?, phone = ?, current_password = ?, role_in_network = ?, account_status = ?, avatar_ref = ? " +
      "WHERE user_id = ?"

  db := get()
  defer db.Close()

  statement, err := db.Prepare(query)

  if err != nil {
    log.Fatal(err)
    return err
  }

  defer statement.Close()

  _, err = statement.Exec(user.Email, user.FirstName, user.SecondName, user.Phone, user.CurrentPass, user.Role,
    user.AccStatus, user.AvatarRef, id)

  if err != nil {
    log.Fatal(err)
    return err
  }

  return nil
}

func (dao UserDaoImpl)FindById(id int) (dao.User, error) {
  user, err := findByUniqueParameter("id", id)

  if err != nil {
    log.Fatal(err)
    return user, err
  }

  return user, nil
}

func (dao UserDaoImpl)FindByEmail(email string) (dao.User, error) {
  user, err := findByUniqueParameter("email", email)

  if err != nil {
    log.Fatal(err)
    return user, err
  }

  return user, nil
}

func (dao UserDaoImpl)FindByPhone(phone string) (dao.User, error) {
  user, err := findByUniqueParameter("phone", phone)

  if err != nil {
    log.Fatal(err)
    return user, err
  }

  return user, nil
}

func (dao UserDaoImpl) InitUsersTable() error {
  query := "CREATE TABLE IF NOT EXISTS users_data (" +
    "user_id SERIAL PRIMARY KEY, " +
      "first_name VARCHAR(50) NOT NULL, " +
        "second_name VARCHAR(50) NOT NULL, " +
          "email VARCHAR(100) NOT NULL, " +
            "phone VARCHAR(20)," +
              "current_password VARCHAR(200) NOT NULL," +
                "role_in_network ENUM('admin', 'user') NOT NULL," +
                  "account_status ENUM('active', 'deleted') NOT NULL," +
                    "avatar_ref MEDIUMTEXT);"

  db := get()
  defer db.Close()

  statement, err := db.Prepare(query)

  if err != nil {
    log.Fatal(err)
    return err
  }

  _, err = statement.Exec()
  if err != nil {
    log.Fatal(err)
    return err
  }

  return nil
}

func findByUniqueParameter(parameterName string, parameterValue interface{}) (dao.User, error) {
  query := fmt.Sprintf("SELECT * FROM users_data WHERE %s = ?", parameterName)

  db := get()
  defer db.Close()

  statement, err := db.Prepare(query)
  user := dao.User{}

  if err != nil {
    log.Fatal(err)
    return user, err
  }

  defer statement.Close()

  result := statement.QueryRow(parameterValue)

  err = result.Scan(
    &user.ID,
    &user.Email,
    &user.FirstName,
    &user.SecondName,
    &user.CurrentPass,
    &user.Role,
    &user.AccStatus,
    &user.AvatarRef,
  )

  if err != nil {
    return user, err
  }

  return user, nil
}
