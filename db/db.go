package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 54322
	username = "admindb"
	password = "adminpass"
	dbname   = "echodemo"
)

type Resource struct {
	Conn *sql.DB
}

// Close :: to close database connection
func (r *Resource) Close() {
	fmt.Println("Closing all db connections")
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var sqliteHandler = new(Resource)

func Execute(params map[interface{}]interface{}) {

	sql := params["sql"].(string)
	fmt.Println(sql)
	result, err := sqliteHandler.Conn.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func Query(params map[interface{}]interface{}) []User {
	sql := params["sql"].(string)
	fmt.Println(sql)
	var users []User
	rows, err := sqliteHandler.Conn.Query(sql)
	defer rows.Close()
	for rows.Next() {
		var (
			id    int
			name  string
			email string
		)
		err := rows.Scan(&id, &name, &email)
		if err != nil {
			fmt.Println(err)
		}
		var user User
		user.ID = id
		user.Name = name
		user.Email = email
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	return users
}

func ConnDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	// defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected! DB")

	sqliteHandler.Conn = db
}
