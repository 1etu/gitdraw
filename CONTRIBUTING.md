# contributing

## setup

```bash
# requirements
go 1.21+
wails (for gui)

# clone
git clone https://github.com/1etu/gitdraw
cd gitdraw

# cli only
go build && ./gitdraw

# gui
wails dev -tags gui
```

## build tags

- cli only: `go build`
- gui: `wails build -tags gui`

same codebase produces both. gui build includes cli via `--cli` flag.

## adding a character

edit `font/font.go`:

```go
'X': {
    {1, 0, 0, 0, 1},
    {0, 1, 0, 1, 0},
    {0, 0, 1, 0, 0},
    {0, 1, 0, 1, 0},
    {1, 0, 0, 0, 1},
    {0, 0, 0, 0, 0},
    {0, 0, 0, 0, 0},
},
```

each character is a 5×7 grid. `1` = filled pixel.

## tests

```bash
go test ./...
```

## pull requests

- one feature per pr
- test your changes
- run `go fmt` before committing

## ideas

- image-to-graph converter
- more fonts
- animation (multi-year)
- web app
- qr code generator
- undo/redo in gui
- shape tools (rect, circle, line)

see [IDEAS.md](IDEAS.md) for full roadmap.

## questions

open an issue.
