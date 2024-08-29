package main

import "github.com/go-pg/pg/v10"

func initPostgres() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "admin",
		Password: "root",
		Database: "miniurl",
	})

	return db
}
