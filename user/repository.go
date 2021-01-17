package user

import (
	"fmt"
	"nevergo/db"
)

// Execute :: desc
func Execute(params map[interface{}]interface{}) {
	sql := params["sql"].(string)
	fmt.Println(sql)
	result, err := db.SqliteHandler.Conn.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

// Query :: desc
func Query(params map[interface{}]interface{}) []User {
	sql := params["sql"].(string)
	fmt.Println(sql)
	var users []User
	rows, err := db.SqliteHandler.Conn.Query(sql)
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
