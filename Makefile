NAME := tintfetch
VERSION ?= 0.1.0-dev
GO ?= go
GO_LDFLAGS ?= -buildid=
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
MANDIR ?= $(PREFIX)/share/man
DOCDIR ?= $(PREFIX)/share/doc/$(NAME)
LICENSEDIR ?= $(PREFIX)/share/licenses/$(NAME)
DESTDIR ?=
SOURCE_DATE_EPOCH ?= $(shell git log -1 --format=%ct 2>/dev/null || date +%s)
export GOTOOLCHAIN := local
export GOPROXY := off
export GOSUMDB := off

.PHONY: all build check install uninstall clean dist
all: build
build:
	mkdir -p build
	$(GO) build -buildvcs=false -trimpath -ldflags '$(GO_LDFLAGS) -X github.com/ffgan/tintfetch/cmd.Version=$(VERSION)' -o build/$(NAME) .
check:
	$(GO) test ./...
	$(GO) vet ./...
install: build
	install -Dm755 build/$(NAME) '$(DESTDIR)$(BINDIR)/$(NAME)'
	install -Dm644 docs/tintfetch.1 '$(DESTDIR)$(MANDIR)/man1/$(NAME).1'
	install -Dm644 LICENSE '$(DESTDIR)$(LICENSEDIR)/LICENSE'
	install -Dm644 THIRD_PARTY_NOTICES '$(DESTDIR)$(LICENSEDIR)/THIRD_PARTY_NOTICES'
	install -Dm644 README.md '$(DESTDIR)$(DOCDIR)/README.md'
	install -Dm644 docs/compatibility.md '$(DESTDIR)$(DOCDIR)/compatibility.md'
	install -Dm644 docs/provenance.md '$(DESTDIR)$(DOCDIR)/provenance.md'
	install -Dm644 docs/config.conf '$(DESTDIR)$(DOCDIR)/examples/config.conf'
uninstall:
	rm -f '$(DESTDIR)$(BINDIR)/$(NAME)' '$(DESTDIR)$(MANDIR)/man1/$(NAME).1'
	rm -f '$(DESTDIR)$(LICENSEDIR)/LICENSE' '$(DESTDIR)$(LICENSEDIR)/THIRD_PARTY_NOTICES'
	rm -f '$(DESTDIR)$(DOCDIR)/README.md' '$(DESTDIR)$(DOCDIR)/compatibility.md' '$(DESTDIR)$(DOCDIR)/provenance.md' '$(DESTDIR)$(DOCDIR)/examples/config.conf'
clean:
	rm -rf build
dist:
	mkdir -p build
	tar --sort=name --mtime='@$(SOURCE_DATE_EPOCH)' --owner=0 --group=0 --numeric-owner --exclude=./.git --exclude=./build --exclude=./.vscode --transform='s,^\.,$(NAME)-$(VERSION),' -cf build/$(NAME)-$(VERSION).tar .
	gzip -n -c build/$(NAME)-$(VERSION).tar > build/$(NAME)-$(VERSION).tar.gz
	rm -f build/$(NAME)-$(VERSION).tar
