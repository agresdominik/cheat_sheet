EXECUTABLE = cheatsh
PACKAGES = ./src/.
BINDIR = bin
PREFIX ?= $(HOME)/.local
DATADIR = $(PREFIX)/share/cheatsh

all: build

build:
	@mkdir -p bin
	go build -o bin/$(EXECUTABLE) $(PACKAGES)

install: build

	@echo "Installing binary in $(PREFIX)/bin"
	install -Dm755 bin/$(EXECUTABLE) $(PREFIX)/bin/$(EXECUTABLE)

	@echo "Installing command lists in $(DATADIR)"
	install -Dm644 data/commands.json $(DATADIR)/commands.json
	install -Dm644 data/commands_template.json $(DATADIR)/commands_template.json

uninstall:
	rm -f $(PREFIX)/bin/$(EXECUTABLE)
	rm -rf $(DATADIR)

clean:
	rm -rf bin/

docker:
	docker build -t cheatsh-test .
	docker run -it --rm cheatsh-test /bin/sh

local: build
	./bin/cheatsh --config ./data/commands.json
