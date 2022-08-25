.PHONY: ctl server

ctl:
	go build -o shortyctl cmd/shortyctl/**.go

server:
	go build -o server cmd/server/**.go
