GOPATH := /go
GOBIN  := $(GOPATH)/bin
PATH   := $(GOROOT)/bin:$(PATH)

all: deps organizer

deps: $(DEPS)
	GOPATH=$(GOPATH) glide install

test: deps
		GOPATH=$(GOPATH) go test -v $(shell glide novendor)

organizer: deps 
    # always format code
		GOPATH=$(GOPATH) go fmt $(glide novendor)
    # vet it
		GOPATH=$(GOPATH) go tool vet *.go
    # binary
		GOPATH=$(GOPATH) go build -o $@ -v $(glide novendor)
		touch $@

linux64: deps
    # always format code
		GOPATH=$(GOPATH) go fmt $(glide novendor)
    # vet it
		GOPATH=$(GOPATH) go tool vet *.go
    # binary
		GOOS=linux GOARCH=amd64 GOPATH=$(GOPATH) go build -o organizer-linux-amd64.bin -v $(glide novendor)
		touch organizer-linux-amd64.bin

win64: deps
    # always format code
		GOPATH=$(GOPATH) go fmt $(glide novendor)
    # vet it
		GOPATH=$(GOPATH) go tool vet *.go
    # binary
		GOOS=windows GOARCH=amd64 GOPATH=$(GOPATH) go build -o organizer-win-amd64.exe -v $(glide novendor)
		touch organizer-win-amd64.exe

.PHONY: $(DEPS) clean

clean:
		rm -rf organizer organizer-win-amd64.exe organizer-linux-amd64.bin .glide vendor

