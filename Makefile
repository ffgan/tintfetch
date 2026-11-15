NAME := tintfetch
VERSION ?= 0.3.7-dev
GO ?= go
PYTHON ?= python3
INSTALL ?= install
GO_LDFLAGS ?= -buildid=
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
MANDIR ?= $(PREFIX)/share/man
DOCDIR ?= $(PREFIX)/share/doc/$(NAME)
LICENSEDIR ?= $(PREFIX)/share/licenses/$(NAME)
DESTDIR ?=
SOURCE_DATE_EPOCH ?= $(shell git log -1 --format=%ct 2>/dev/null || echo 0)
export SOURCE_DATE_EPOCH
export GOTOOLCHAIN := local
export GOPROXY := off
export GOSUMDB := off

.PHONY: all build check install uninstall clean dist check-provenance
all: build
build:
	mkdir -p build
	$(GO) build -buildvcs=false -trimpath -ldflags '$(GO_LDFLAGS) -X github.com/ffgan/tintfetch/cmd.Version=$(VERSION)' -o build/$(NAME) .
check: check-provenance
	$(GO) test ./...
	$(GO) vet ./...
check-provenance:
	$(PYTHON) scripts/check-provenance.py
	$(PYTHON) scripts/generate-arm-ids.py --check
install: build
	$(INSTALL) -Dm755 build/$(NAME) '$(DESTDIR)$(BINDIR)/$(NAME)'
	$(INSTALL) -Dm644 docs/tintfetch.1 '$(DESTDIR)$(MANDIR)/man1/$(NAME).1'
	$(INSTALL) -Dm644 LICENSE '$(DESTDIR)$(LICENSEDIR)/LICENSE'
	$(INSTALL) -Dm644 COPYING '$(DESTDIR)$(LICENSEDIR)/COPYING'
	$(INSTALL) -Dm644 THIRD_PARTY_NOTICES '$(DESTDIR)$(LICENSEDIR)/THIRD_PARTY_NOTICES'
	$(INSTALL) -Dm644 README.md '$(DESTDIR)$(DOCDIR)/README.md'
	$(INSTALL) -Dm644 docs/compatibility.md '$(DESTDIR)$(DOCDIR)/compatibility.md'
	$(INSTALL) -Dm644 docs/provenance.md '$(DESTDIR)$(DOCDIR)/provenance.md'
	$(INSTALL) -Dm644 docs/provenance.json '$(DESTDIR)$(DOCDIR)/provenance.json'
	$(INSTALL) -Dm644 docs/assets.sha256 '$(DESTDIR)$(DOCDIR)/assets.sha256'
	$(INSTALL) -Dm644 docs/config.conf '$(DESTDIR)$(DOCDIR)/examples/config.conf'
uninstall:
	rm -f '$(DESTDIR)$(BINDIR)/$(NAME)' '$(DESTDIR)$(MANDIR)/man1/$(NAME).1'
	rm -f '$(DESTDIR)$(LICENSEDIR)/LICENSE' '$(DESTDIR)$(LICENSEDIR)/COPYING' '$(DESTDIR)$(LICENSEDIR)/THIRD_PARTY_NOTICES'
	rm -f '$(DESTDIR)$(DOCDIR)/README.md' '$(DESTDIR)$(DOCDIR)/compatibility.md' '$(DESTDIR)$(DOCDIR)/provenance.md' '$(DESTDIR)$(DOCDIR)/provenance.json' '$(DESTDIR)$(DOCDIR)/assets.sha256' '$(DESTDIR)$(DOCDIR)/examples/config.conf'
clean:
	rm -rf build
dist: check-provenance
	$(PYTHON) scripts/create-source-archive.py '$(VERSION)' 'build/$(NAME)-$(VERSION).tar.gz'
