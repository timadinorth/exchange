all: run

clean:
	docker-compose down --volumes

build:
	go build -o bin/ex 

doc:
	swag init --parseDependency

run: doc build 
	docker compose up --build 

test: 
	ENV=test docker compose down --volumes
	ENV=test docker compose -f docker-compose.yml -f docker-compose.test.yml up --build --abort-on-container-exit

cover: 
	go test -coverprofile=coverage.out -v ./...
	go tool cover -html=coverage.out