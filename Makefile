run:
	@go run .

build:
	@go build -o bin .

test:
	@go test -v ./...

install:
	@go install

uninstall:
	@go clean -i
