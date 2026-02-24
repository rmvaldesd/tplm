BINARY  := tplm
PREFIX  := /usr/local/bin

.PHONY: build install uninstall clean test lint

build:
	go build -o $(BINARY) .

install: build
	install -m 755 $(BINARY) $(PREFIX)/$(BINARY)

uninstall:
	rm -f $(PREFIX)/$(BINARY)

clean:
	rm -f $(BINARY)

test:
	go test ./...

lint:
	golangci-lint run
