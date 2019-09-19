SHELL 	   := $(shell which bash)

## BOF define block
PROJECT    := scp-aws-lambda
BINARY     := scp-aws-lambda

ROOT_DIR   := $(shell git rev-parse --show-toplevel)
BIN_DIR    := $(ROOT_DIR)/bin
REL_DIR    := $(ROOT_DIR)/release
SRC_DIR    := $(ROOT_DIR)/cmd/scp-aws-lambda

VERSION    :=`git describe --tags 2>/dev/null`
COMMIT     :=`git rev-parse --short HEAD 2>/dev/null`
DATE       :=`date "+%FT%T%z"`

LDBASE     := github.com/gertd/$(PROJECT)
LDFLAGS    := -ldflags "-w -s -X $(LDBASE)/cmd.version=${VERSION} -X $(LDBASE)/cmd.date=${DATE} -X $(LDBASE)/cmd.commit=${COMMIT}"

<<<<<<< HEAD
PLATFORMS  := linux darwin
=======
PLATFORMS  := linux
>>>>>>> master
OS         = $(word 1, $@)

GOARCH     ?= amd64
GOOS       :=
ifeq ($(OS),Windows_NT)
	GOOS = windows
else 
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		GOOS = linux
	endif
	ifeq ($(UNAME_S),Darwin)
		GOOS = darwin
	endif
endif

LINTER     := $(BIN_DIR)/golangci-lint
LINTVERSION:= v1.15.0
TESTRUNNER := $(GOPATH)/bin/gotestsum

NO_COLOR   :=\033[0m
OK_COLOR   :=\033[32;01m
ERR_COLOR  :=\033[31;01m
WARN_COLOR :=\033[36;01m
ATTN_COLOR :=\033[33;01m

## EOF define block

.PHONY: all
all: build test lint

deps:
	@echo -e "$(ATTN_COLOR)==> download dependencies $(NO_COLOR)"
	@GO111MODULE=on go mod download

.PHONY: build
build: deps 
	@echo -e "$(ATTN_COLOR)==> build GOOS=$(GOOS) GOARCH=$(GOARCH) VERSION=$(VERSION) COMMIT=$(COMMIT) DATE=$(DATE) $(NO_COLOR)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) GO111MODULE=on go build $(LDFLAGS) -v -o $(BIN_DIR)/$(BINARY)-$(GOOS)-$(GOARCH) $(SRC_DIR)

$(TESTRUNNER):
	@echo -e "$(ATTN_COLOR)==> get gotestsum test runner  $(NO_COLOR)"
	@go get -u gotest.tools/gotestsum 

.PHONY: test 
test: $(TESTRUNNER) gen
	@echo -e "$(ATTN_COLOR)==> test $(NO_COLOR)"
	@gotestsum --format short-verbose ./...

$(LINTER):
	@echo -e "$(ATTN_COLOR)==> get  $(NO_COLOR)"
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s $(LINTVERSION)
 
.PHONY: lint
lint: $(LINTER)
	@echo -e "$(ATTN_COLOR)==> lint $(NO_COLOR)"
	@$(LINTER) run --enable-all
	@echo -e "$(NO_COLOR)\c"

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	@mkdir -p $(REL_DIR)
	
	@echo -e "$(ATTN_COLOR)==> release to $(REL_DIR) $(NO_COLOR)"

	@echo -e "$(ATTN_COLOR)==> build GOOS=$(@:release-%=%) GOARCH=$(GOARCH) VERSION=$(VERSION) COMMIT=$(COMMIT) DATE=$(DATE) $(NO_COLOR)"
	@GOOS=$(@:release-%=%) GOARCH=$(GOARCH) GO111MODULE=on go build $(LDFLAGS) -v -o $(REL_DIR)/$(BINARY)-$(@:release-%=%)-$(GOARCH) $(SRC_DIR)

	@echo -e "$(ATTN_COLOR)==> zip $(REL_DIR)/$(BINARY)-$(@:release-%=%)-$(VERSION).zip $(NO_COLOR)"
	@zip -j $(REL_DIR)/$(BINARY)-$(@:release-%=%)-$(VERSION).zip $(REL_DIR)/$(BINARY)-$(@:release-%=%)-$(GOARCH) >/dev/null

.PHONY: release
release: $(PLATFORMS)

.PHONY: install
install:
	@echo -e "$(ATTN_COLOR)==> install $(NO_COLOR)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) GO111MODULE=on go install $(LDFLAGS) $(SRC_DIR)

.PHONY: clean
clean:
	@echo -e "$(ATTN_COLOR)==> clean $(NO_COLOR)"
	@rm -rf $(BIN_DIR)
	@rm -rf $(REL_DIR)

.PHONY: gen
<<<<<<< HEAD
gen: deps
=======
gen:
>>>>>>> master
	@echo -e "$(ATTN_COLOR)==> generate $(NO_COLOR)"
	@go generate ./...
