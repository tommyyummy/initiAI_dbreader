#!/bin/sh
cwd=$(pwd)

env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dbreader *.go

# admin123
scp dbreader tommy@172.27.177.11:/home/tommy/ttu/dbreader.temp
scp -r config/* tommy@172.27.177.11:/home/tommy/ttu/config/