![Go](https://github.com/pakohan/craftdoor/workflows/Go/badge.svg)

# craftdoor

A RFID based access control system written in Go for Raspberry Pi + MFRCC522
tag readers + MIFARE RFID tags.

# Project Overview

Craftdoor is a software suite for an RFID-powered door access system on a
federation of Raspberry Pi devices. With the exception of the "master", each
Raspberry Pi is connected to an [RFID
reader](https://www.nxp.com/docs/en/data-sheet/MFRC522.pdf) and a door.
Registered members may tap the RFID reader with their MIFARE RFID tag to open
the adjacent door.

The system is administered via a WebUI interface and accompanying REST API
served by the master device. Persistent state is stored in a SQLite database on
the master device. See below for valid endpoints.

Instructions below for building, configuring, and launching the webserver.

**Note**: At time of writing, only a single "master" Raspberry Pi is supported.

# Installation

To start the software suite, do the following on the master Raspberry Pi device,

1. Connect RC522 to master Raspberry Pi's hardware SPI interface. Follow
   instructions [here](https://github.com/pakohan/craftdoor.git).
1. Download `golang` from https://golang.org. Follow installation instructions
   [here](https://golang.org/doc/install#install). Verify that go is installed
   by running `go version` in a terminal. Expect to see >= 1.14.
1. Install GCC cross-compiler,
  ```
  $ sudo apt install gcc-arm-linux-gnueabi libc6-armel-cross \
    libc6-dev-armel-cross binutils-arm-linux-gnueabi
  ```
1. Run `cmd/master/main.go`. This will launch a webserver listening on port 8080.
  ```
  $ git clone https://github.com/pakohan/craftdoor.git
  $ cd craftdoor/cmd/master
  $ go run main.go develop.json
  ```

**Note**: If the RC522 RFID reader is not
[detected](http://pkg.go.dev/periph.io/x/periph/host/rpi#Present), a fake,
dummy interface will be used. This dummy interface cannot interact with RFID
tags.

# Usage

Once `main.go` is launched, the following endpoints are available via the HTTP
webserver,

- `POST /init`: Writes a default key to an active RFID tag. Only for testing.
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
