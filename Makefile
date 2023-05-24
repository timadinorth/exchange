build:
	go build -o bin/ex 

run: build
	./bin/ex

test:
	go test -coverprofile=coverage.out -v ./...

cover: test
	go tool cover -html=coverage.out