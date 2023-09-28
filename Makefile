################################################################################
#### INSTALLATION VARS
################################################################################
PREFIX=$(HOME)

################################################################################
#### BUILD VARS
################################################################################
BIN=fred
BINDIR=bin
RELEASEDIR=release
HEAD=$(shell git describe --dirty --long --tags 2> /dev/null  || git rev-parse --short HEAD)
COHASH=$(shell git rev-parse --short HEAD)
TIMESTAMP=$(shell TZ=UTC date '+%FT%TZ')
TEST_COVER_FILE=$(BIN)-test-coverage.out
# TIMESTAMP=$(shell date '+%Y-%m-%dT%H:%M:%S %z %Z')

LDFLAGS="-X 'main.binName=$(BIN)' -X 'main.buildVersion=$(HEAD)' -X 'main.buildTimestamp=$(TIMESTAMP)' -X 'main.compiledBy=$(shell go version)'"
# LDFLAGS="-X 'main.buildVersion=$(HEAD)' -X 'main.buildTimestamp=$(TIMESTAMP)'"

all: local

.PHONY: version
version:
	@printf "\n\n%s\n\n" $(HEAD)

################################################################################
#### HOUSE CLEANING
################################################################################

.PHONY: dep
dep:
	go mod tidy

.PHONY: check
check: _setup
	goimports -w ./
	go vet

################################################################################
#### UN/INSTALL
################################################################################

.PHONY: _setup
_setup:
	mkdir -p $(BINDIR)

.PHONY: clean
clean:
	rm -f $(BIN) $(BIN)-* $(BINDIR)/$(BIN) $(BINDIR)/$(BIN)-*

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

.PHONY: demo
demo: local
	./$(BINDIR)/$(BIN) -file "example/short" -debug

.PHONY: test
test: dep check
	GOEXPERIMENT=loopvar go test -tags memory -covermode=count ./...

.PHONY: test-cover
test-cover:
	go mod tidy
	go test -tags memory -covermode=count -coverprofile $(TEST_COVER_FILE) ./...
	go tool cover -html=$(TEST_COVER_FILE)

################################################################################
#### MACOS BUILDS
################################################################################

.PHONY: local
local: dep check
	GOEXPERIMENT=loopvar go build -ldflags $(LDFLAGS) -tags "$(FREDMODE)" -o $(BINDIR)/$(BIN)

.PHONY: release
release: dep check
	GOEXPERIMENT=loopvar go build -ldflags $(LDFLAGS) -tags "$(FREDMODE)" -o $(BINDIR)/$(BIN)
	tar -czf $(RELEASEDIR)/$(BIN)-$(COHASH)-$(GOOS)-$(GOARCH).tgz -C $(BINDIR) $(BIN)

.PHONY: prod
prod: dep check
	GOEXPERIMENT=loopvar GOWORK=off go build -ldflags $(LDFLAGS) -o $(BINDIR)/$(BIN)
