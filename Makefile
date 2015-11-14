PORT ?= 4343
GO ?= GO15VENDOREXPERIMENT=1 go


build:
	$(GO) build

gin:
	$(GO) get github.com/codegangsta/gin
	gin --immediate --port=$(PORT) ./main.go
