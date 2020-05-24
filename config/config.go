package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// ReadConfig reads the config file.
// If the service is started with an argument it's assuming this is a path to a config JSON file.
// If no argument is passed it checks:
//   * /etc/craftdoor/master.json
//   * ./develop.json
func ReadConfig() (Config, error) {
	if len(os.Args) > 1 {
		return readFile(os.Args[1])
	}

	for _, f := range []string{"/etc/craftdoor/master.json", "./develop.json"} {
		_, err := os.Stat(f)
		if err == nil {
			return readFile(f)
		}
	}

	return Config{}, fmt.Errorf("could not find config file")
}

func readFile(filename string) (Config, error) {
	log.Printf("reading config from '%s'", filename)
	f, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer func() {
		e := f.Close()
		if e != nil {
			log.Printf("failed closing config file: %s", e.Error())
		}
	}()

	cfg := Config{}
	return cfg, json.NewDecoder(f).Decode(&cfg)
}

// Config represents the config file contents
type Config struct {
	MasterKey  string `json:"master_key"`
	SQLiteFile string `json:"sqlite_file"`
	ListenHTTP string `json:"listen_http"`
	Device     string `json:"device"`
	RSTPin     string `json:"rst_pin"`
	IRQPin     string `json:"irq_pin"`
}
