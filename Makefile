GOCMD=GO111MODULE=on go

test:
	$(GOCMD) test -cover -race ./...

bench:
	$(GOCMD) test -run=NONE -bench=. -benchmem  ./...

lint:
	golangci-lint run

.PHONY: lint test bench