package user

import (
	"fmt"
	"nevergo/db"
)

// Execute :: desc
func execute(params map[interface{}]interface{}) {
	sql := params["sql"].(string)
	fmt.Println(sql)
	_, err := db.SqliteHandler.Conn.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(result)
}

// Query :: desc
func query(params map[interface{}]interface{}) []User {
	sql := params["sql"].(string)
	fmt.Println(sql)
	var users []User
	rows, err := db.SqliteHandler.Conn.Query(sql)
	defer rows.Close()
	for rows.Next() {
		var (
			id       int
			username string
			fname    string
			lname    string
			email    string
			tel      string
		)
		err := rows.Scan(&id, &username, &fname, &lname, &email, &tel)
		if err != nil {
			fmt.Println(err)
		}
		var user User
		user.ID = id
		user.Username = username
		user.Fname = fname
		user.Lname = lname
		user.Email = email
		user.Telephone = tel
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	return users
}
