#!/bin/bash
#
# Build and release
#

set -e

SRC="cmd/master/"
DST="release/"

# cd into project root.
PROJECT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"/../
cd "${PROJECT_ROOT}"

# Clean the release/ directory.
rm -rf "${DST}"
mkdir -p "${DST}"

# Build for ARM.
env GOOS=linux \
    GOARCH=arm \
    GOARM=5 \
    CGO_ENABLED=1 \
    CC=arm-linux-gnueabi-gcc \
    go build -o "${DST}/main" "${SRC}/main.go"

# Copy auxiliary files.
cp ${SRC}/develop.json ${DST}/
cp ${SRC}/schema.sql ${DST}/

echo "Finished building '${SRC}'. Copy '${DST}' to your RPi and run it. For example,"
echo "$ rsync -r release/ pi@raspberrypi:/home/pi/craftdoor"
echo "$ ssh pi@raspberrypi 'cd /home/pi/craftdoor && ./main develop.json'"
