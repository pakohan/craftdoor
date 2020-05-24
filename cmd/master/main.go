package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pakohan/craftdoor/config"
	"github.com/pakohan/craftdoor/controller"
	"github.com/pakohan/craftdoor/lib"
	"github.com/pakohan/craftdoor/model"
	"github.com/pakohan/craftdoor/service"
	"periph.io/x/periph/host/rpi"
)

func main() {
	log.SetFlags(log.Llongfile)

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Panic(err)
	}

	db, err := openDB(cfg.SQLiteFile)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()

	c := make(chan os.Signal, 1)
	signal.Notify(c)
	go func() {
		defer wg.Done()
		<-c
		log.Printf("shutting down DB")
		e := db.Close()
		if e != nil {
			log.Printf("err closing db: %s", e.Error())
		}
	}()

	err = start(cfg, db, wg)
	if err != nil {
		c <- os.Interrupt
		log.Panic(err)
	}

}

func start(cfg config.Config, db *sqlx.DB, wg *sync.WaitGroup) error {
	cl := lib.NewChangeListener()

	var r lib.Reader
	if rpi.Present() {
		var err error
		r, err = lib.NewRC522Reader(cfg, cl)
		if err != nil {
			return err
		}
	} else {
		r =lib.NewDummyReader()
	}

	m := model.New(db)
	s := service.New(m, r, cl)
	c := controller.New(m, s)

	srv := http.Server{
		Addr:    cfg.ListenHTTP,
		Handler: c,
	}

	wg.Add(1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig)
	go func() {
		defer wg.Done()
		<-sig
		log.Printf("shutting down HTTP server")
		e := srv.Shutdown(context.Background())
		if e != nil {
			log.Printf("err closing HTTP server: %s", e.Error())
		}
	}()

	log.Printf("listening on %s", cfg.ListenHTTP)
	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}

func openDB(filename string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	err = lib.InitDBSchema(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
