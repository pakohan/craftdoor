package config

import (
	"encoding/json"
	"os"
)

func ReadConfig() (Config, error) {
	cfg := Config{}
	f, err := os.Open("/etc/craftdoor/master.json")
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	return cfg, json.NewDecoder(f).Decode(&cfg)
}

type Config struct {
	MasterKey  string `json:"master_key"`
	SQLiteFile string `json:"sqlite_file"`
	ListenHTTP string `json:"listen_http"`
	Device     string `json:"device"`
	RSTPin     string `json:"rst_pin"`
	IRQPin     string `json:"irq_pin"`
}
