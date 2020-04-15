#!/bin/bash


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o ddns-server ./ddns-server.go
tar -czvf linux_amd64.tar.gz ddns-server conf
rm -f ddns-server

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o ddns-server.exe ./ddns-server.go
tar -czvf windows_amd64.tar.gz ddns-server.exe conf
rm -f ddns-server.exe
