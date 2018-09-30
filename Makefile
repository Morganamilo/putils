.PHONY: FORCE default build doc all install install_doc \
	install_all uninstall uninstall_doc uninstall_all\
	clean

PREFIX := /usr/local
DESTDIR :=

BIN := pedit
DOC := pedit.1
GO := go

default: build
FORCE:
build: ${BIN}

doc:
	for f in ${DOC}; do \
		a2x -d manpage -f manpage doc/$$f.asciidoc; \
	done

all: build doc

%: %.go FORCE
	${GO} build -o bin/$@ $<

install:
	for f in ${BIN}; do \
		install -Dm755 bin/$$f $(DESTDIR)$(PREFIX)/bin/$$f; \
	done

install_doc:
	for f in ${DOC}; do \
		install -Dm755 doc/$$f $(DESTDIR)$(PREFIX)/share/man/man1/$$f; \
	done

install_all: install install_doc

uninstall:
	for f in ${BIN}; do \
		rm -f $(DESTDIR)$(PREFIX)/bin/$$f; \
	done

uninstall_doc:
	for f in ${DOC}; do \
		rm -f $(DESTDIR)$(PREFIX)/share/man/man1/$$f; \
	done

uninstall_all: uninstall uninstall_doc

clean:
	for f in ${BIN}; do \
		rm -f bin/$$f; \
	done

	for f in ${DOC}; do \
		rm -f doc/$$f; \
	done
