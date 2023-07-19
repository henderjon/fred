################################################################################
#### INSTALLATION VARS
################################################################################
PREFIX=$(HOME)

################################################################################
#### BUILD VARS
################################################################################
BIN=ffred
BINDIR=bin
HEAD=$(shell git describe --dirty --long --tags 2> /dev/null  || git rev-parse --short HEAD)
TIMESTAMP=$(shell TZ=UTC date '+%FT%T %Z')
TEST_COVER_FILE=$(BIN)-test-coverage.out
# TIMESTAMP=$(shell date '+%Y-%m-%dT%H:%M:%S %z %Z')

LDFLAGS="-X 'main.binName=$(BIN)' -X 'main.buildVersion=$(HEAD)' -X 'main.buildTimestamp=$(TIMESTAMP)' -X 'main.compiledBy=$(shell go version)'"
LDFLAGS_STATIC="-X 'main.binName=$(BIN)' -X 'main.buildVersion=$(HEAD)' -X 'main.buildTimestamp=$(TIMESTAMP)' -X 'main.compiledBy=$(shell go version)' -extldflags=-static"
STATIC_TAGS="osusergo,netgo,sqlite_omit_load_extension"

all: local

################################################################################
#### HOUSE CLEANING
################################################################################

.PHONY: _setup
_setup:
	mkdir -p $(BINDIR)

.PHONY: clean
clean:
	rm -f $(BIN) $(BIN)-* $(BINDIR)/$(BIN) $(BINDIR)/$(BIN)-*

.PHONY: dep
dep:
	go mod tidy

.PHONY: version
version:
	@printf "\n\n%s\n\n" $(HEAD)

.PHONY: check
check: _setup
	golint
	goimports -w ./
	gofmt -w ./
	go vet

################################################################################
#### INSTALL
################################################################################

.PHONY: install
install: local
	mkdir -p $(PREFIX)/$(BINDIR)
	mv $(BINDIR)/$(BIN) $(PREFIX)/$(BINDIR)/$(BIN)
	@echo "\ninstalled $(BIN) to $(PREFIX)/$(BINDIR)\n"


.PHONY: uninstall
uninstall:
	rm -f $(PREFIX)/$(BINDIR)/$(BIN)

################################################################################
#### TESING
################################################################################

.PHONY: test
test: dep check
	go test -covermode=count ./...

.PHONY: test-cover
test-cover:
	go mod tidy
	go test -covermode=count -coverprofile $(TEST_COVER_FILE) ./...
	go tool cover -html=$(TEST_COVER_FILE)

################################################################################
#### MACOS BUILDS
################################################################################

.PHONY: local
local: check static
	go build -ldflags $(LDFLAGS) -o $(BINDIR)/$(BIN)

.PHONY: prod
prod: check static
	GOWORK=off go build -ldflags $(LDFLAGS) -o $(BINDIR)/$(BIN)

################################################################################
#### CROSS COMPILE TO LINUX +CGO BUILDS
################################################################################

.PHONY: cgo
cgo: check static docker-cgo
	$(info docker build --tag $(HEAD) .)
	$(info docker run --rm -p 80:9000 --detach $(HEAD) ./$(BIN)-linux64-$(HEAD) -debug -stand-alone -local-port 9000)

.PHONY: docker-cgo
docker-cgo: clean
# bind mount $PWD to /usr/local/src; set the working dir to /usr/local/src; run make linux64-cgo
	docker run --rm --volume "$(PWD)":/usr/local/src --workdir /usr/local/src golang:latest make linux64-cgo

.PHONY: linux64-cgo
linux64-cgo: _setup
	env CGO_ENABLED=1 go build -a -ldflags $(LDFLAGS_STATIC) -tags $(STATIC_TAGS) -o $(BINDIR)/$(BIN)-linux64-$(HEAD) .
