package main

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/pakohan/craftdoor/config"
	"github.com/pakohan/craftdoor/controller"
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
	s, err := service.New(cfg, m)
	if err != nil {
		return err
	}
	c := controller.New(m, s)

	return http.ListenAndServe(cfg.ListenHTTP, c)
}
