FROM golang:1.12


# Configure environment and rebuild stdlib
ENV CGO_ENABLED=1 GO15VENDOREXPERIMENT=1 GODIR=github.com/moul/radioman WEBDIR=/web
RUN go install -a std


# Install deps
RUN go get github.com/tools/godep \
 && go get github.com/golang/lint/golint \
 && go get golang.org/x/tools/cmd/goimports \
 && go get golang.org/x/tools/cmd/stringer


# Install dependencies
RUN apt-get update \
 && apt-get install -y -qq libtagc0-dev pkg-config \
 && apt-get clean


# Run the project

WORKDIR /go/src/${GODIR}/radioman
COPY pkg /go/src/${GODIR}/radioman/pkg
COPY vendor /go/src/${GODIR}/radioman/vendor
COPY cmd /go/src/${GODIR}/radioman/cmd
COPY Godeps /go/src/${GODIR}/radioman/Godeps
RUN go install ./cmd/radiomand
COPY web ${WEBDIR}
ENTRYPOINT ["radiomand"]
