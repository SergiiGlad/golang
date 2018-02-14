package main

import (
	"database/sql" // пакет для форматированного ввода вывода
	"encoding/json"
	"fmt"
	"net/http" // пакет для поддержки HTTP протокола
	//"strings" // пакет для работы с  UTF-8 строками
	"log" // пакет для логирования

	_ "github.com/go-sql-driver/mysql" // Дополнительный модуль для mysql
	"gopkg.in/hlandau/passlib.v1/hash/bcrypt"
)

// Данные юзера для восстановления пароля
type UserToRestore struct {
	userID int
	email  string
}

type UserToRegister struct {
	email     string `json:"email"`
	firstName string `json:"firstName"`
	lastName  string `json:"lastName"`
	phone     string `json:"phone"`
	password  string `json:"-"`
}

func scanRows(db *sql.DB, searchmail string) UserToRestore {
	var user UserToRestore
	rows, err := db.Query("select user_id, email from users_data where email = ?", searchmail)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.userID, &user.email)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(user.userID, user.email)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return user
}

func openDataBase(db **sql.DB) {
	var err error
	*db, err = sql.Open("mysql",
		"root:123456@tcp(127.0.0.1:3306)/gohum")
	if err != nil {
		log.Fatal(err)
	}
}

func insertNewPass(db *sql.DB, user *UserToRestore, hashPass string){
	query := "INSERT INTO users_passwords (password) VALUES(?) where user_id = user.UserID"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(hashPass)
	if err != nil {
		log.Fatal(err)
	}
}

func writeGenerPass(user UserToRestore) {
	pass := NewPassword(6)
	fmt.Println(pass)
	hashPass, err := bcrypt.Crypter.Hash(pass)
	if err != nil {
		log.Fatal(err)
	}
	var db *sql.DB
	openDataBase(&db)
	defer db.Close()
	insertNewPass(db, &user, hashPass)
	fmt.Println(hashPass)

}

func restoreByEmail(email string) bool {
	var db *sql.DB
	openDataBase(&db)
	defer db.Close()
	user := scanRows(db, email)
	if user.email == "" {
		return false
	} else {

		writeGenerPass(user)
		return true
	}
}

func insertNewUser(db *sql.DB, user *UserToRegister) {
	query := "INSERT INTO users_data (first_name, second_name, email, phone, current_password) VALUES(?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(user.firstName, user.lastName, user.email, user.phone, user.password)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
}

func saveToDatabase(user *UserToRegister) {
	var db *sql.DB
	openDataBase(&db)
	defer db.Close()
	insertNewUser(db, user)
}

func RegisterNewUserHandler(rw http.ResponseWriter, req *http.Request) {

	var user UserToRegister
	var values interface{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&user)
	if err != nil {
		values = map[string]string{"status": "error"}
	} else {
		saveToDatabase(&user)
		values = map[string]string{"status": "success"}
	}
	jsonValue, _ := json.Marshal(values)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonValue)
	defer req.Body.Close()
}

func RestoreRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //анализ аргументов,
	//fmt.Println(r.Form) // ввод информации о форме на стороне сервера
	//fmt.Println("path", r.URL.Path)
	//fmt.Println("scheme", r.URL.Scheme)
	//fmt.Println(r.Form["url_long"])
	email := r.Form["email"][0]
	if restoreByEmail(email) {
		values := map[string]string{"status": "success"}
		jsonValue, _ := json.Marshal(values)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonValue)
	} else {
		values := map[string]string{"status": "error"}
		jsonValue, _ := json.Marshal(values)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write(jsonValue)
	}
}

func main() {
	http.HandleFunc("/restore", RestoreRouterHandler)
	http.HandleFunc("/register", RegisterNewUserHandler) // установим роутер
	err := http.ListenAndServe(":9000", nil)             // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
