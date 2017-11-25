.PHONY: client server

all: client server

client:
	go install .

server:
	@mkdir -p pkg
	env CGO_ENABLED=0 go build -o pkg/plex-static --ldflags '-extldflags "-static"' .
