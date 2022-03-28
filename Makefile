GOCMD=GO111MODULE=on go

test:
	$(GOCMD) test -cover -race ./...

bench:
	$(GOCMD) test -run=NONE -bench=. -benchmem  ./...

.PHONY: linters-install lint test bench