build:
	@go build -o bin/gobank

run: build
	@./bin/gobank

seed: build
	@./bin/gobank --seed

test:
	@go test -v ./testing/unit_test.go
