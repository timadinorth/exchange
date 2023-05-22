build:
	go build -o bin/ex 

run: build
	./bin/ex

test:
	go test -v ./...