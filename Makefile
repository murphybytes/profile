.PHONY: test lint

test:
	go test -race -coverprofile c.out ./...

lint:
	golint ./...
