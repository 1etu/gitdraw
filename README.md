# gitdraw

turn your github contribution graph into pixel art.

draw on a graph or type text, pick colors, push to github. uses `git fast-import` to generate thousands of commits in seconds.

## showcase / demo

### text-to-graph
![text2graph demo](https://github.com/user-attachments/assets/4d12010f-c379-45af-a3c6-da8d827d7429)
![drawing tool](https://github.com/user-attachments/assets/2e6c4adc-d5d1-4f7d-be95-ae4b2897d912)


## install

### gui (recommended)

download from [releases](https://github.com/1etu/gitdraw/releases/latest):

- macOS (Apple Silicon): `GitDraw-macOS-AppleSilicon.zip`
- macOS (Intel): `GitDraw-macOS-Intel.zip`
- Windows: `GitDraw-Windows-x64.exe`

### homebrew

```bash
brew install 1etu/tap/gitdraw
```

### go install (cli only)

```bash
go install github.com/1etu/gitdraw@latest
```

### build from source

```bash
git clone https://github.com/1etu/gitdraw
cd gitdraw

# cli only
go build

# gui (requires wails)
wails build -tags gui
```

## usage

### gui

launch the app, draw on the graph or type text, pick intensity levels, enter your repo url, hit generate & push.

the gui build also supports cli mode:

```bash
./GitDraw --cli
./GitDraw --help
```

### cli

```bash
$ gitdraw

  gitdraw — contribution graph art

  Text to draw HELLO
  
  Preview:

  Sun ██░░░░░░██░░██████████░░██░░░░░░░░░░██░░░░░░░░░░░░██████░░
  Mon ██░░░░░░██░░██░░░░░░░░░░██░░░░░░░░░░██░░░░░░░░░░██░░░░░░██
  Tue ██░░░░░░██░░██░░░░░░░░░░██░░░░░░░░░░██░░░░░░░░░░██░░░░░░██
  Wed ██████████░░██████████░░██░░░░░░░░░░██░░░░░░░░░░██░░░░░░██
  Thu ██░░░░░░██░░██░░░░░░░░░░██░░░░░░░░░░██░░░░░░░░░░██░░░░░░██
  Fri ██░░░░░░██░░██░░░░░░░░░░██░░░░░░░░░░██░░░░░░░░░░██░░░░░░██
  Sat ██░░░░░░██░░██████████░░██████████░░██████████░░░░██████░░

  Continue with this design? (y/n) y
  Target year (2025) 
  Fill background? (y/n) y
  Text intensity (15) 
  
  ✓ Repository ready
  Configure GitHub remote? (y/n) y
```

### supported characters

```
A-Z  0-9  space  !  .  -  _  :  /  <  >
```

---

## how it works

github's contribution graph is a 7×52 grid (days × weeks).

1. renders text/drawing using a custom 5×7 pixel font
2. maps each pixel to a specific date in the target year
3. creates backdated commits via `git fast-import`
4. pushes to your github repo

commit intensity controls color shade:

- 1-3 commits → light green
- 4-6 commits → medium green
- 7-9 commits → darker green
- 10+ commits → darkest green

## authentication

uses your system's git credentials:

- ssh keys (recommended)
- macos keychain
- git credential helper
- github cli (`gh auth login`)

---

## project structure

```
gitdraw/
├── cli.go          # cli entry point
├── cli_gui.go      # cli code for gui build
├── gui.go          # gui entry point (wails)
├── draw/           # grid and text rendering
├── font/           # 5x7 pixel font definitions
├── git/            # git operations
├── gui/            # frontend (html/css/js)
└── build.sh        # release build script
```

## contributing

see [CONTRIBUTING.md](CONTRIBUTING.md).

ideas:
- image-to-graph converter
- more fonts
- animation support (multi-year)
- web app version

## license

MIT
