PREFIX ?= $(HOME)/.local
BINDIR = $(PREFIX)/bin
GOCMD = go
GOBUILD = $(GOCMD) build
LDFLAGS = -s -w

.PHONY: build install clean test vet

build:
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o octk .

install: build
	mkdir -p $(BINDIR)
	install -m 755 octk $(BINDIR)/octk

clean:
	rm -f octk

vet:
	$(GOCMD) vet ./...

test:
	$(GOCMD) test ./...
