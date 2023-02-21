package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/stan.go"
)

func initDB() *sqlx.DB {

	dsn := "host=localhost user=postgres password=92usubop dbname=postgres port=5432 sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)

	if err != nil {
		fmt.Println("Error in postgres connection: ", err)
	}
	return db
}

func insertDb(key string, m *stan.Msg, db *sqlx.DB) {

	_, err := db.Exec(`INSERT INTO student (order_id,data) VALUES ($1,$2)`, key, m.Data)

	if err != nil {
		fmt.Println("erro insert to db")
	}
}
