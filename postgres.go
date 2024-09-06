package main

import (
	"os"

	"github.com/go-pg/pg/v10"
)

func getPostgresAddr() string {
	env := os.Getenv("ENV")

	if env == "prod" {
		return "postgres:5432"
	}

	return "localhost:5432"
}

func initPostgres() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "admin",
		Password: "root",
		Database: "miniurl",
		Addr:     getPostgresAddr(),
	})

	return db
}

func checkIfSlugExists(db *pg.DB, slug string) (bool, error) {
	urlQueryResult := &Slug{}
	err := db.Model(urlQueryResult).Column("slug").Where("slug = ?", slug).Select()

	return err == nil, err
}

func saveSlug(db *pg.DB, slug *Slug) (pg.Result, error) {
	return db.Model(slug).Insert()
}
