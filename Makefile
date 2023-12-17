run-unit-tests:
	go test -v -count 1 ./...

stop: 
	docker compose down

build:
	rm -rf bin/*
	go mod tidy
	go build -o ./bin/api ./cmd/main.go

run-api:
	./bin/api