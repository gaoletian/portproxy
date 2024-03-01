export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all: fmt clean build

build: 
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/portproxy .

fmt:
	go fmt ./...
	
clean:
	rm -f ./bin/portproxy