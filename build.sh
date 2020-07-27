#!/bin/bash

# client
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o ddns-client ./main-code/client/ddns-client.go
tar -czvf ddns-client.linux-amd64.tar.gz ddns-client
rm -f ddns-client

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o ddns-client.exe ./main-code/client/ddns-client.go
tar -czvf ddns-client.windows-amd64.tar.gz ddns-client.exe
rm -f ddns-client.exe

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s" -o ddns-client ./main-code/client/ddns-client.go
tar -czvf ddns-client.darwin-amd64.tar.gz ddns-client
rm -f ddns-client

# mips
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags "-w -s" -o ddns-client ./main-code/client/ddns-client.go
tar -czvf ddns-client.linux-mips64le.tar.gz ddns-client
rm -f ddns-client

CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags "-w -s" -o ddns-client ./main-code/client/ddns-client.go
tar -czvf ddns-client.linux-mips64.tar.gz ddns-client
rm -f ddns-client

CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "-w -s" -o ddns-client ./main-code/client/ddns-client.go
tar -czvf ddns-client.linux-mipsle.tar.gz ddns-client
rm -f ddns-client

CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "-w -s" -o ddns-client ./main-code/client/ddns-client.go
tar -czvf ddns-client.linux-mips.tar.gz ddns-client
rm -f ddns-client

# arm
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-w -s" -o ddns-client ./main-code/client/ddns-client.go
tar -czvf ddns-client.linux-arm64.tar.gz ddns-client
rm -f ddns-client

CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-w -s" -o ddns-client ./main-code/client/ddns-client.go
tar -czvf ddns-client.linux-arm_v7.tar.gz ddns-client
rm -f ddns-client

# server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o ddns-server ./main-code/server/ddns-server.go
tar -czvf ddns-server.linux-amd64.tar.gz ddns-server
rm -f ddns-server

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o ddns-server.exe ./main-code/server/ddns-server.go
tar -czvf ddns-server.windows-amd64.tar.gz ddns-server.exe
rm -f ddns-server.exe

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s" -o ddns-server ./main-code/server/ddns-server.go
tar -czvf ddns-server.darwin-amd64.tar.gz ddns-server
rm -f ddns-server

# mips
CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags "-w -s" -o ddns-server ./main-code/server/ddns-server.go
tar -czvf ddns-server.linux-mips64le.tar.gz ddns-server
rm -f ddns-server

CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags "-w -s" -o ddns-server ./main-code/server/ddns-server.go
tar -czvf ddns-server.linux-mips64.tar.gz ddns-server
rm -f ddns-server

CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "-w -s" -o ddns-server ./main-code/server/ddns-server.go
tar -czvf ddns-server.linux-mipsle.tar.gz ddns-server
rm -f ddns-server

CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "-w -s" -o ddns-server ./main-code/server/ddns-server.go
tar -czvf ddns-server.linux-mips.tar.gz ddns-server
rm -f ddns-server

# arm
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-w -s" -o ddns-server ./main-code/server/ddns-server.go
tar -czvf ddns-server.linux-arm64.tar.gz ddns-server
rm -f ddns-server

CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-w -s" -o ddns-server ./main-code/server/ddns-server.go
tar -czvf ddns-server.linux-arm_v7.tar.gz ddns-server
rm -f ddns-server
