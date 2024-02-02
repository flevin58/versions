ifeq ($(OS),Windows_NT)
BINFILE=versions.exe
else
BINFILE=versions
endif

run:
	@go run .

build:
	@go build -o bin/$(BINFILE) .

test:
	@go test -v ./...

install:
	@go install

uninstall:
	@go clean -i
