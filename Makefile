HASGOLINTCI := $(shell command -v golintci-lint 2> /dev/null)

ifdef HASGOLINTCI
    GOLINT=golangci-lint
else
    GOLINT=bin/golangci-lint
endif

GOLANGCILINT_VERSION = 1.50.1

.PHONY: install
install:
	go install -v ./...

.PHONY: start
start: main.go
	go run main.go

.PHONY: serve
serve: main.go
	air

.PHONY: nodemon
nodemon: main.go
	nodemon --exec go run main.go --signal SIGTERM

.PHONY: clean
clean:
	@echo "==> Cleaning..."
	rm -rf ./**/*.{o,exe}

tidy:
	go fmt ./**/*.go

ifdef HASGOCILINT
golangci-lint:
	@echo "Skip this"
else
golangci-lint: golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} golangci-lint
endif

golangci-lint-${GOLANGCI_VERSION}:
	@echo "==> Installing golangci-lint ${GOLANGCI_VERSION}"
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v${GOLANGCI_VERSION}

.PHONY: lint
lint: golangci-lint ## Run linter
	@echo "==> Running linter..."
	$(GOLINT) run

.PHONY: fix
fix: golangci-lint ## Fix lint violations
	$(GOLINT) run --fix
	gofmt -s -w .
	goimports -w .

.PHONY: list-todo
list-todo: ## Detect FIXME, TODO and other comment keywords
	golangci-lint run --enable=godox --disable-all

# Fix lint violations with gofmt and goimports
.PHONY: fmt
fmt:
	gofmt -s -w ./...
	goimports -w ./...
