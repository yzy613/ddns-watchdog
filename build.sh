#!/bin/bash

# client
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o ddns-client ./main-code/client/ddns-client.go
tar -czvf linux_amd64_client.tar.gz ddns-client
rm -f ddns-client

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o ddns-client.exe ./main-code/client/ddns-client.go
tar -czvf windows_amd64_client.tar.gz ddns-client.exe
rm -f ddns-client.exe

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-w -s" -o ddns-client ./main-code/client/ddns-client.go
tar -czvf linux_arm64_client.tar.gz ddns-client
rm -f ddns-client

# server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o ddns-server ./main-code/server/ddns-server.go
tar -czvf linux_amd64_server.tar.gz ddns-server
rm -f ddns-server

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o ddns-server.exe ./main-code/server/ddns-server.go
tar -czvf windows_amd64_server.tar.gz ddns-server.exe
rm -f ddns-server.exe

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-w -s" -o ddns-server ./main-code/server/ddns-server.go
tar -czvf linux_arm64_server.tar.gz ddns-server
rm -f ddns-server
