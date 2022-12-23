.PHONY: shorty server

shorty:
	go build -o shorty cmd/shorty/**.go

server:
	CGO_ENABLED=0 go build -o shorty-server cmd/server/**.go

docker:
	docker build --platform linux/amd64 -t ghcr.io/nwehr/shorty/server:$(shell git rev-parse --short=8 HEAD) -f server.Dockerfile .
	docker push ghcr.io/nwehr/shorty/server:$(shell git rev-parse --short=8 HEAD)