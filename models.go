package main

type Slug struct {
	Url  string `pg:"url"`
	Slug string `pg:"slug"`
}
