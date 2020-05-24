package lib

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

// InitDBSchema initializes the DB schema if the sqlite_master table has no entries
func InitDBSchema(db *sqlx.DB) error {
	var count int
	err := db.Get(&count, checkTables)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	log.Printf("didn't find any tables, will init db schema")

	f, err := os.Open("./schema.sql")
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(b))
	return err
}

const checkTables = `
SELECT COUNT(*)
FROM sqlite_master
WHERE type='table'
ORDER BY name`
