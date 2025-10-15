EXECUTABLE = cheatsh
PACKAGES = ./src/.
BINDIR = bin
SYSCONFDIR = /etc/cheatsh
PREFIX ?= /usr/local

all: build

build:
	@mkdir -p bin
	go build -o bin/$(EXECUTABLE) $(PACKAGES)

install: build

	@echo "Installing binary in $(PREFIX)"
	install -Dm755 bin/$(EXECUTABLE) $(PREFIX)/bin/$(EXECUTABLE)

	@echo "Copying config files to $(SYSCONFDIR)"
	install -Dm644 data/commands.json $(SYSCONFDIR)/commands.json
	install -Dm644 data/commands_template.json $(SYSCONFDIR)/commands_template.json

uninstall:
	rm -f $(PREFIX)/bin/$(EXECUTABLE)
	rm -rf $(SYSCONFDIR)

clean:
	rm -rf bin/

docker:
	docker build -t cheatsh-test .
	docker run -it --rm cheatsh-test /bin/sh
