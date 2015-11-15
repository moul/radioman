FROM golang:1.5


# Configure environment and rebuild stdlib
ENV CGO_ENABLED=1 GO15VENDOREXPERIMENT=1
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
COPY . /go/src/app
WORKDIR /go/src/app
RUN go build -o radiomand .
ENTRYPOINT /go/src/app/radiomand
