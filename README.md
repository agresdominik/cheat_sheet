# Cheat Sheet

A small TUI tool that displays a pre-configured set of commands in a searchable
menu. Pick one and it's copied to your system clipboard — paste it where you
need it and edit before running.

I wrote this to look up commands I use often enough to need repeatedly, but not
often enough to remember exactly. Written in Go with [bubbletea](https://github.com/charmbracelet/bubbletea).

## Install

### Homebrew (recommended)

```bash
brew install agresdominik/repo/cheatsh
```

### From source

```bash
git clone https://github.com/agresdominik/cheat_sheet
cd cheat_sheet
make install              # builds and installs to ~/.local/bin
```

Make sure `~/.local/bin` is in your `PATH` (or change `PREFIX` in the `Makefile`).

## Configuration

cheatsh ships with a maintained command list and uses it by default. It lives
next to the binary and is **replaced on `brew upgrade`**, so updating the list
in this repo and releasing is all it takes to get the new commands on your
machines:

```
<prefix>/share/cheatsh/commands.json            # Homebrew, e.g. /opt/homebrew/...
~/.local/share/cheatsh/commands.json            # make install
data/commands.json                              # running from a source checkout
```

To maintain the default list, edit `data/commands.json` and release. The format
is a map of categories to command lists:

```json
{
  "git": [
    { "command": "git reset --soft HEAD~1", "desc": "Undo last commit, keep changes" }
  ]
}
```

### Custom list

If you'd rather keep your own list (one that survives upgrades untouched), copy
the template and pass it explicitly — `--config` always wins over the bundled
default:

```bash
cp "$(brew --prefix)/share/cheatsh/commands_template.json" ~/my-commands.json
cheatsh --config ~/my-commands.json
```

## Clipboard dependency

Copying to the clipboard needs a backend:

- **macOS** — works out of the box (`pbcopy`).
- **Linux/Wayland** — `wl-clipboard`.
- **Linux/X11** — `xsel` or `xclip`.

If no backend is found, cheatsh prints the command to stdout instead.

## Usage

```bash
cheatsh                  # browse and select; selection is copied to clipboard
cheatsh --config FILE    # use a specific config file
cheatsh --help
```

In the TUI: `/` to filter, `enter` to descend / select, `b` to go back, `q` to quit.

## Development

```bash
make build               # build to bin/cheatsh
make local               # build and run against data/commands.json
go test ./src/...
```
