CREATE DATABASE miniurl OWNER admin;

\c miniurl

CREATE TABLE slugs (
	url TEXT NOT NULL,
	slug CHAR(10) NOT NULL
);