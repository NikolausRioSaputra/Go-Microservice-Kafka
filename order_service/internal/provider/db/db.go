package db


import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func GetConnection() (*sql.DB, error) {
	driver := "postgres"
	connString := "postgres://postgres:root@localhost:5432/orderKafka?sslmode=disable"

	DB, err := sql.Open(driver, connString)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	err = DB.Ping()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fmt.Println("Successfully Connected to the database")
	return DB, nil
}
