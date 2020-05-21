package main

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pakohan/craftdoor/config"
	"github.com/pakohan/craftdoor/controller"
	"github.com/pakohan/craftdoor/lib"
	"github.com/pakohan/craftdoor/model"
	"github.com/pakohan/craftdoor/service"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Panic(err)
	}

	db, err := sqlx.Connect("sqlite3", cfg.SQLiteFile)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = start(cfg, db)
	if err != nil {
		log.Panic(err)
	}
}

func start(cfg config.Config, db *sqlx.DB) error {
	m := model.New(db)
	s := service.New(cfg, m)
	c := controller.New(m, s)

	_, err := lib.NewReader(cfg, s)
	if err != nil {
		return err
	}

	return http.ListenAndServe(cfg.ListenHTTP, c)
}
