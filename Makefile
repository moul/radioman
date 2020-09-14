PORT ?= 4343
GOENV ?= GO15VENDOREXPERIMENT=1
GO ?= $(GOENV) go
COMPOSE_TARGET ?=
ADMIN_PASSWORD ?= radioman
COMPOSE_ENV ?= ICECAST_SOURCE_PASSWORD="$(ADMIN_PASSWORD)" ICECAST_ADMIN_PASSWORD="$(ADMIN_PASSWORD)" ICECAST_PASSWORD="$(ADMIN_PASSWORD)" ICECAST_RELAY_PASSWORD="$(ADMIN_PASSWORD)"
SOURCES := $(shell find . -name "*.go")
DOCKER_HOST ?= tcp://127.0.0.1:2376
DOCKER_HOST_IP := $(shell echo $(DOCKER_HOST) | cut -d/ -f3 | cut -d: -f1)
DOCKER_COMPOSE ?= docker-compose -pradioman
RADIOMAND := radiomand-$(shell uname -s)-$(shell uname -m)

all: build

.PHONY: install
install:
	cd radioman && go install ./cmd/radiomand

.PHONY: liquidsoap-telnet
liquidsoap-telnet:
	nc -v localhost 2300

.PHONY: docker-telnet
docker-telnet:
	socat readline TCP:$(DOCKER_HOST_IP):2300

.PHONY: test-liquidsoap-config
test-liquidsoap-config:
	docker run -it -u liquidsoap --rm -v "$(PWD)/liquidsoap:/liquidsoap" moul/liquidsoap liquidsoap --verbose --debug /liquidsoap/main.liq

.PHONY: docker-exec-liquidsoap
docker-exec-liquidsoap:
	docker exec -it `docker-compose ps -q liquidsoap` /bin/bash

.PHONY: docker-exec-radiomand
docker-exec-radiomand:
	docker exec -it `docker-compose ps -q radiomand` /bin/bash

.PHONY: docker-exec-icecast
docker-exec-icecast:
	docker exec -it `docker-compose ps -q icecast` /bin/bash

.PHONY: compose
compose:
	$(DOCKER_COMPOSE) build radiomand
	$(DOCKER_COMPOSE) kill
	$(DOCKER_COMPOSE) rm -fv
	$(COMPOSE_ENV) $(DOCKER_COMPOSE) up -d $(COMPOSE_TARGET)
	$(DOCKER_COMPOSE) logs -f $(COMPOSE_TARGET)

.PHONY: clean
clean:
	rm -f radiomand-*-*
	find . -name gin-bin -delete
	$(DOCKER_COMPOSE) kill
	$(DOCKER_COMPOSE) rm -fv

.PHONY: tidy
tidy:
	cd radioman; go mod tidy

.PHONY: deps
deps:
	sudo apt install libtagc0-dev

docker.build:
	docker build -t moul/radioman ./radioman
