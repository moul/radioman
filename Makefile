PORT ?= 4343
GOENV ?= GO15VENDOREXPERIMENT=1
GO ?= $(GOENV) go
COMPOSE_TARGET ?=
PASSWORD ?= toor
COMPOSE_ENV ?= ICECAST_SOURCE_PASSWORD="$(PASSWORD)" ICECAST_ADMIN_PASSWORD="$(PASSWORD)" ICECAST_PASSWORD="$(PASSWORD)" ICECAST_RELAY_PASSWORD="$(PASSWORD)"
SOURCES := $(find . -name "*.go")
DOCKER_HOST ?= tcp://127.0.0.1:2376
DOCKER_HOST_IP := $(shell echo $(DOCKER_HOST) | cut -d/ -f3 | cut -d: -f1)


all: build


.PHONY: build
build: radioman


.PHONY: docker-telnet
docker-telnet:
	socat readline TCP:$(DOCKER_HOST_IP):2300


radioman: $(SOURCES)
	$(GO) build -o $@ .


.PHONY: compose
compose:
	docker-compose build radioman
	docker-compose kill
	docker-compose rm -f
	$(COMPOSE_ENV) docker-compose up -d $(COMPOSE_TARGET)
	docker-compose logs $(COMPOSE_TARGET)


.PHONY: gin
gin:
	$(GO) get github.com/codegangsta/gin
	gin --immediate --port=$(PORT) ./main.go


.PHONY: clean
clean:
	rm -f radioman gin-bin
