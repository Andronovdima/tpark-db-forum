package store

import (
	"database/sql"
	"io/ioutil"
)

func CreateTables(db *sql.DB) error {

	//userQ := `CREATE TABLE IF NOT EXISTS users (
	//	nickname varchar primary key,
	//	about varchar,
	//	fullname varchar,
	//	email varchar not null unique
	//);`
	//if _, err := db.Exec(userQ); err != nil {
	//	return err
	//}
	//
	//forumsQ := `CREATE TABLE IF NOT EXISTS forums (
	//	slug varchar not null primary key,
	//	posts bigint not null,
	//	threads bigint not null,
	//	title varchar not null,
	//	author varchar not null references users
	//);`
	//if _, err := db.Exec(forumsQ); err != nil {
	//	return err
	//}
	//
	//
	//threadsQ := `CREATE TABLE IF NOT EXISTS threads (
	//	id bigserial not null primary key,
	//	author varchar not null references users,
	//	created timestamptz DEFAULT now(),
	//	forum varchar not null references forums,
	//	message varchar not null,
	//	slug varchar,
	//	title varchar not null,
	//	votes int
	//);`
	//if _, err := db.Exec(threadsQ); err != nil {
	//	return err
	//}
	//
	//postQ := `CREATE TABLE IF NOT EXISTS posts (
	//	id bigserial not null primary key,
	//	author varchar not null references users,
	//	created timestamptz DEFAULT now(),
	//	forum varchar not null references forums,
	//	isEdited boolean not null DEFAULT false,
	//	message varchar not null,
	//	parent bigint not null DEFAULT NULL,
	//	thread bigint references threads,
	//	path bigint[]  NOT NULL DEFAULT '{0}'
	//);`
	//
	//if _, err := db.Exec(postQ); err != nil {
	//	return err
	//}
	//
	//voteQ := `CREATE TABLE IF NOT EXISTS votes (
	//	id bigserial not null primary key,
	//	nickname varchar references users,
	//	thread bigint references threads,
	//	voice int
	//);`
	//	if _, err := db.Exec(voteQ); err != nil {
	//		return err
	//	}
	//
	//fu := `CREATE TABLE IF NOT EXISTS forum_users (
	//	id bigserial not null primary key,
	//	forum varchar not null references forums,
	//	nickname varchar not null references users,
	//);`
	//if _, err := db.Exec(fu); err != nil {
	//	return err
	//}



	file, err := ioutil.ReadFile("store/main.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(file))
	if err != nil {
		return err
	}

	return nil
}

