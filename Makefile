run:
	@go run .

build:
	@go build -o bin/versions .

test:
	@go test -v ./...

install:
	@go install

uninstall:
	@go clean -i
