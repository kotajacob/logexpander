# logexpander
# See LICENSE for copyright and license details.
.POSIX:

PREFIX ?= /usr/local
MANPREFIX ?= $(PREFIX)/share/man
GO ?= go
GOFLAGS ?=
RM ?= rm -f

all: logexpander

logexpander:
	$(GO) build $(GOFLAGS)

clean:
	$(RM) logexpander

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f logexpander $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/logexpander

uninstall:
	$(RM) $(DESTDIR)$(PREFIX)/bin/logexpander

.DEFAULT_GOAL := all

.PHONY: all logexpander clean install uninstall
