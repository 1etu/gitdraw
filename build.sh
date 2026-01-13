#!/bin/bash
set -e

v="${1:-1.0.0}"
d=./dist
wails=$(command -v wails || echo ~/go/bin/wails)

rm -rf "$d" && mkdir -p "$d"

for p in darwin/arm64:macOS-AppleSilicon darwin/amd64:macOS-Intel windows/amd64:Windows-x64; do
	plat=${p%:*} name=${p#*:}
	echo "gui: $name"
	$wails build -tags gui -platform "$plat" -clean -skipbindings 2>/dev/null ||
		$wails build -tags gui -platform "$plat" -clean
	ext=$([[ $plat == windows* ]] && echo .exe || echo .app)
	mv ./build/bin/gitdraw$ext "$d/GitDraw-$name$ext"
done

for t in darwin/arm64:macos-arm64 darwin/amd64:macos-amd64 windows/amd64:windows-x64 linux/amd64:linux-amd64; do
	os=${t%/*} arch=${t#*/}; arch=${arch%:*}; name=${t#*:}
	echo "cli: $name"
	ext=$([[ $os == windows ]] && echo .exe || echo "")
	GOOS=$os GOARCH=$arch go build -o "$d/gitdraw-cli-$name$ext" .
done

cd "$d"
for f in *; do
	[[ -e "$f" ]] && zip -rq "${f%.*}-v$v.zip" "$f"
done

echo && ls -lh *.zip
