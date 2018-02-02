package store

import (
	"database/sql"
	"fmt"
	"os"
	"tracy/log"
)

/*DB is the one global used to gain access to the database from this package.
 * Other packages, like testing, might choose to not use this database and instead
 * will supply their own. */
var DB gorm.DB

/*Open the database and create the tables if they aren't already created.
 * Errors indicate something incorrectly happened while
 * connecting. Don't forget to close this DB when finished using it. */
func Open(driver, path string) error {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Error.Printf(err.Error())
	}

	db.Exec("PRAGMA foreign_keys = ON")
	db.Exec("PRAGMA journal_mode = WAL")
	db.LogMode(true)
	db.AutoMigrate(
		&types.Tracer{},
		&types.TracerEvent{},
		&types.DOMContext{},
		&types.Request{})

	/* We want to disable the goroutine thread pool that is used by default since this application doesn't need it and will
	 * cause performance issues. https://stackoverflow.com/questions/35804884/sqlite-concurrent-writing-performance */
	db.DB().SetMaxOpenConns(1)

	/* Return the database and nil, indicating we made a sound connection. */
	DB = db

	return err
}
