package main

import (
	"github.com/jackc/pgx"
	"github.com/usrpro/dotpgx"
	"github.com/usrpro/pgxmgr"
	log "gopkg.in/inconshreveable/log15.v2"
)

var db *dotpgx.DB

func initDb() {
	var dbConf = pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     flags.pgHost,
			User:     flags.pgUser,
			Database: flags.pgDatabase,
		},
		MaxConnections: 50,
	}
	var err error
	db, err = dotpgx.New(dbConf)
	if err != nil {
		panic(err)
	}
	if err = pgxmgr.Run(db, "sql/migrations"); err != nil {
		panic(err)
	}
	if flags.test {
		if err = loadTestData(); err != nil {
			log.Error("loadTestData failed", "error", err)
		}
	}
	if err = db.ParsePath("sql/app"); err != nil {
		panic(err)
	}
}

func loadTestData() (err error) {
	if err = db.ClearMap(); err != nil {
		return
	}
	if err = db.ParseFiles("sql/tests/cat_tree.sql"); err != nil {
		return
	}
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	for _, q := range db.List() {
		if _, err = tx.Exec(q); err != nil {
			return
		}
	}
	err = tx.Commit()
	return
}
