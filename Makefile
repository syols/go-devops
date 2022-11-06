build:
	echo "Build agent"
	go build -o bin/agent cmd/agent/main.go
	echo "Build agent"
	go build -o bin/server cmd/server/main.go

server:
	go run cmd/server/main.go

agent:
	go run cmd/agent/main.go

get:
	go get golang.org/x/tools/cmd/goimports
	go get github.com/kisielk/errcheck
	go get golang.org/x/lint
	go get github.com/tools/godep

imports:
	goimports -l -w .

fmt:
	go fmt ./...

lint:
	golint ./...

vet:
	go vet -v ./...

errors:
	errcheck -ignoretests -blank ./...

deps:
	godep restore

test: deps
	go test -v ./...

coverage: deps
	gocov test ./... > $(CURDIR)/coverage.out 2>/dev/null
	gocov report $(CURDIR)/coverage.out
	if test -z "$$CI"; then \
	  gocov-html $(CURDIR)/coverage.out > $(CURDIR)/coverage.html; \
	  if which open &>/dev/null; then \
	    open $(CURDIR)/coverage.html; \
	  fi; \
	fi

run: server

