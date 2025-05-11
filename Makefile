run:
	go run ./cmd/api

build 
	go build -o bin/api ./cmd/api

test
	go test -v ./tests/...