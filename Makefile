# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=go-lc3-vm
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build: 
				$(GOBUILD) -o $(BINARY_NAME) -v
test: 
				$(GOTEST) -v ./...
clean: 
				$(GOCLEAN)
				rm -f $(BINARY_NAME)
				rm -f $(BINARY_UNIX)
run:
				$(GOBUILD) -o $(BINARY_NAME) -v ./...
				./$(BINARY_NAME) prog/2048.obj
