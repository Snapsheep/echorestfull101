package db

import (
	"database/sql"
	"fmt"
	"os"
)

type Resource struct {
	Conn *sql.DB
}

// Close :: to close database connection
func (r *Resource) Close() {
	fmt.Println("Closing all db connections")
}

var SqliteHandler = new(Resource)

func ConnDB() {
	host := os.Getenv("HOST_ENDPOINT")
	port := os.Getenv("POSTGRES_PORT")
	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%v port=%v user=%v "+
		"password=%v dbname=%v sslmode=disable",
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

	SqliteHandler.Conn = db
}
