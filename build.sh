#!/bin/bash

# client
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o watchdog-ddns-client ./main-code/client/watchdog-ddns-client.go
tar -czvf watchdog-ddns-client.linux-amd64.tar.gz watchdog-ddns-client
rm -f watchdog-ddns-client

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o watchdog-ddns-client.exe ./main-code/client/watchdog-ddns-client.go
tar -czvf watchdog-ddns-client.windows-amd64.tar.gz watchdog-ddns-client.exe
rm -f watchdog-ddns-client.exe

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s" -o watchdog-ddns-client ./main-code/client/watchdog-ddns-client.go
tar -czvf watchdog-ddns-client.darwin-amd64.tar.gz watchdog-ddns-client
rm -f watchdog-ddns-client

# mips
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags "-w -s" -o watchdog-ddns-client ./main-code/client/watchdog-ddns-client.go
tar -czvf watchdog-ddns-client.linux-mips64le.tar.gz watchdog-ddns-client
rm -f watchdog-ddns-client

CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags "-w -s" -o watchdog-ddns-client ./main-code/client/watchdog-ddns-client.go
tar -czvf watchdog-ddns-client.linux-mips64.tar.gz watchdog-ddns-client
rm -f watchdog-ddns-client

CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "-w -s" -o watchdog-ddns-client ./main-code/client/watchdog-ddns-client.go
tar -czvf watchdog-ddns-client.linux-mipsle.tar.gz watchdog-ddns-client
rm -f watchdog-ddns-client

CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "-w -s" -o watchdog-ddns-client ./main-code/client/watchdog-ddns-client.go
tar -czvf watchdog-ddns-client.linux-mips.tar.gz watchdog-ddns-client
rm -f watchdog-ddns-client

# arm
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-w -s" -o watchdog-ddns-client ./main-code/client/watchdog-ddns-client.go
tar -czvf watchdog-ddns-client.linux-arm64.tar.gz watchdog-ddns-client
rm -f watchdog-ddns-client

CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-w -s" -o watchdog-ddns-client ./main-code/client/watchdog-ddns-client.go
tar -czvf watchdog-ddns-client.linux-arm_v7.tar.gz watchdog-ddns-client
rm -f watchdog-ddns-client

# server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o watchdog-ddns-server ./main-code/server/watchdog-ddns-server.go
tar -czvf watchdog-ddns-server.linux-amd64.tar.gz watchdog-ddns-server
rm -f watchdog-ddns-server

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o watchdog-ddns-server.exe ./main-code/server/watchdog-ddns-server.go
tar -czvf watchdog-ddns-server.windows-amd64.tar.gz watchdog-ddns-server.exe
rm -f watchdog-ddns-server.exe

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s" -o watchdog-ddns-server ./main-code/server/watchdog-ddns-server.go
tar -czvf watchdog-ddns-server.darwin-amd64.tar.gz watchdog-ddns-server
rm -f watchdog-ddns-server

# mips
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags "-w -s" -o watchdog-ddns-server ./main-code/server/watchdog-ddns-server.go
tar -czvf watchdog-ddns-server.linux-mips64le.tar.gz watchdog-ddns-server
rm -f watchdog-ddns-server

CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags "-w -s" -o watchdog-ddns-server ./main-code/server/watchdog-ddns-server.go
tar -czvf watchdog-ddns-server.linux-mips64.tar.gz watchdog-ddns-server
rm -f watchdog-ddns-server

CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "-w -s" -o watchdog-ddns-server ./main-code/server/watchdog-ddns-server.go
tar -czvf watchdog-ddns-server.linux-mipsle.tar.gz watchdog-ddns-server
rm -f watchdog-ddns-server

CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "-w -s" -o watchdog-ddns-server ./main-code/server/watchdog-ddns-server.go
tar -czvf watchdog-ddns-server.linux-mips.tar.gz watchdog-ddns-server
rm -f watchdog-ddns-server

# arm
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-w -s" -o watchdog-ddns-server ./main-code/server/watchdog-ddns-server.go
tar -czvf watchdog-ddns-server.linux-arm64.tar.gz watchdog-ddns-server
rm -f watchdog-ddns-server

CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-w -s" -o watchdog-ddns-server ./main-code/server/watchdog-ddns-server.go
tar -czvf watchdog-ddns-server.linux-arm_v7.tar.gz watchdog-ddns-server
rm -f watchdog-ddns-server
