# Spomac

A lightweight terminal UI for controlling Spotify on macOS.

<p align="center">
  <img src="assets/screenshot.png" alt="Spomac Screenshot" width="600">
</p>

## Features

- Track info: name, artist, album
- Progress bar with elapsed/total time and percentage
- Play/Pause, Next, Previous controls
- Live updates every second
- Adapts to terminal width
- No API keys or OAuth required (uses macOS AppleScript)

## Requirements

- macOS
- Spotify desktop app

## Installation

### Homebrew (recommended)

```bash
brew tap eneslevent/tap
brew install spomac
```

### From source

```bash
git clone https://github.com/eneslevent/spomac.git
cd spomac
go install .
```

## Usage

```bash
spomac
```

## Keybindings

| Key | Action |
|---|---|
| `Space` / `2` | Play / Pause |
| `1` / `←` | Previous track |
| `3` / `→` | Next track |
| `q` / `Ctrl+C` | Quit |

## How it works

Spomac communicates with the Spotify desktop app through macOS AppleScript. No network requests, no API tokens — everything runs locally.

## Built with

- [Go](https://go.dev)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Terminal styling

## License

MIT
