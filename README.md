![Go](https://github.com/pakohan/craftdoor/workflows/Go/badge.svg)

# craftdoor

A RFID based access control system written in Go for Raspberry Pi + MFRCC522
tag readers + MIFARE RFID tags.

# Project Overview

Run the program from the `cmd/master` directory. If one argument is provided, it
assumes this is the path to a config JSON file.  If none is provided, it first
looks up `/etc/craftdoor/master.json` (for production use) and then
`./develop.json` in the `cmd/master` directory.  That second file is already
present for development and provides all defaults needed for local execution.

After starting, it checks whether the sqlite db file mentioned in the config
has any tables. If not, it executes `schema.sql` from the `cmd/master` directory
to set up the schema.

Next, it checks whether
[periph.io/x/periph/host/rpi#Present](http://pkg.go.dev/periph.io/x/periph/host/rpi#Present)
returns true. If so, it uses the GPI pins and spi device file mentioned in the
config file to connect to the RFID reader.  If not, it uses a dummy reader
doing nothing for local development.

In a nutshell: `cd cmd/master && go run main.go` should setup the DB, detect
the platform and start a webserver listening on :8080 for easy development


# Installation

1. Download `golang` from https://golang.org. You may follow the instructions
   [here](https://golang.org/doc/install#install).
1. Run `main.go`. This will launch a webserver listening on port 8080.
  ```
  $ cd cmd/master
  $ go run main.go develop.json
  ```

# Usage

Once `main.go` is launched, the following endpoints are available via the HTTP
webserver,

- `POST /init`: Writes a sector on an active RFID tag. Only for testing.
- `GET /state?id=<UUID>`: Get the state of a particular RFID tag as soon as
  it's put in front of the reader. You must specify which tag you want by its
  UUID.

For doors,

- `GET /doors`: list doors in database
- `POST /doors`: Create a new door.
- `PUT /doors/<id>`: Update an existing door.
- `DELETE /doors/<id>`: Delete an existing door.

Similar to doors, one can query and manager members and roles via `/members` and `/roles`.

# Code Organization

```
cmd/
  master/
    develop.db       # sqlite database used during development
    develop.json     # JSON config used during development
    main.go          # main binary for this project
    schema.sql       # SQL for initializing develop.db
config/
  config.go          # JSON config file API
controller/
  controller.go      # HTTP request handling logic.
  ...
lib/
  change_listener.go
  db.go              # initialize database schema
  rc522.go           # reader implementation for RC522 RFID reader.
  reader.go          # interfaces for an RFID reader
  state.go           # State of an RFID tag.
model/               # database definitions, API
  model.go
  ...
service/             # business logic for adding/removing keys, doors, etc
vendor/              # third-party code
```
