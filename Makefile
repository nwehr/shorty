.PHONY: shorty server

shorty:
	go build -o shorty cmd/shorty/**.go

server:
	go build -o shorty-server cmd/server/**.go

docker:
	docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/nwehr/shorty/server:$(shell git rev-parse --short=8 HEAD) -f server.Dockerfile --push .