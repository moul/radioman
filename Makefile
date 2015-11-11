PORT ?= 4343

gin:
	go get github.com/codegangsta/gin
	gin --immediate --port=$(PORT) ./main.go
