![Go](https://github.com/pakohan/craftdoor/workflows/Go/badge.svg)

# craftdoor

A RFID based access control system written in Go

# Setup

Run the program from the cmd/master directory. If one argument is provided, it assumes this is the path to a config JSON file.
If none is provided, it first looks up `/etc/craftdoor/master.json` (for production use) and then `./develop.json` in the cmd/master directory.
That second file is already present for development and provides all defaults needed for local execution.

After starting, it checks whether the sqlite db file mentioned in the config has any tables. If not, it executes `schema.sql` from the cmd/master directory to set up the schema.

Next, it checks whether [periph.io/x/periph/host/rpi#Present](pkg.go.dev/periph.io/x/periph/host/rpi#Present) returns true. If so, it uses the GPI pins and spi device file mentioned in the config file to connect to the RFID reader.
If not, it uses a dummy reader doing nothing for local development.

In a nutshell: `cd cmd/master && go run main.go` should setup the DB, detect the platform and start a webserver listening on :8080 for easy development
