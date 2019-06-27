.PHONY: test

test:
	go test -race -coverprofile c.out ./...