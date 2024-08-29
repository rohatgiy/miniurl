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

func checkIfSlugExists(db *pg.DB, slug string) (bool, error) {
	urlQueryResult := &Slug{}
	err := db.Model(urlQueryResult).Column("slug").Where("slug = ?", slug).Select()

	return err == nil, err
}

func saveSlug(db *pg.DB, slug *Slug) (pg.Result, error) {
	return db.Model(slug).Insert()
}
