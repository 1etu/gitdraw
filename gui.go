//go:build gui

package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/1etu/gitdraw/draw"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed gui
var assets embed.FS

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type Point struct {
	Week int `json:"week"`
	Day  int `json:"day"`
}

func (a *App) TextToPoints(text string) string {
	grid := draw.Text(strings.ToUpper(text))
	points := grid.Points()
	
	result := make([]Point, len(points))
	for i, p := range points {
		result[i] = Point{Week: p.Week, Day: p.Day}
	}
	
	jsonData, _ := json.Marshal(result)
	return string(jsonData)
}

func (a *App) Generate(pointsJSON string, year, intensity int, fillBg bool, remoteURL string) string {
	var points []Point
	if err := json.Unmarshal([]byte(pointsJSON), &points); err != nil {
		return "error: invalid points data (corrupted)"
	}

	tmpDir, err := os.MkdirTemp("", "gitdraw-*")
	if err != nil {
		return "error: failed to create temp directory (probs no space left)"
	}

	if err := gitInit(tmpDir); err != nil {
		return "error: git init failed (probs git not installed)"
	}

	name, email := getGitUser()
	if name == "" {
		name = "gitdraw"
	}
	if email == "" {
		return "error: git user.email not configured (try 'git config --global user.email')"
	}

	fgDates := pointsToDates(points, year)
	var bgDates []time.Time

	if fillBg {
		bgDates = backgroundDates(points, year)
	}

	if err := fastImport(tmpDir, bgDates, 1, fgDates, intensity, name, email); err != nil {
		return "error: commit generation failed"
	}

	if remoteURL != "" {
		if !strings.Contains(remoteURL, "://") && !strings.HasPrefix(remoteURL, "git@") {
			remoteURL = "https://" + remoteURL
		}

		if err := gitAddRemote(tmpDir, remoteURL); err != nil {
			return "error: failed to add remote"
		}

		if err := gitPush(tmpDir); err != nil {
			return "error: push failed - check your credentials"
		}
	}

	os.RemoveAll(tmpDir)
	return "success"
}

func graphStart(year int) time.Time {
	dec31 := time.Date(year, 12, 31, 12, 0, 0, 0, time.UTC)
	daysSinceSunday := int(dec31.Weekday())
	lastSunday := dec31.AddDate(0, 0, -daysSinceSunday)
	return lastSunday.AddDate(0, 0, -52*7)
}

func pointsToDates(points []Point, year int) []time.Time {
	start := graphStart(year)
	now := time.Now()
	var dates []time.Time

	for _, p := range points {
		d := start.AddDate(0, 0, p.Week*7+p.Day)
		if d.Year() == year && !d.After(now) {
			dates = append(dates, d)
		}
	}
	return dates
}

func backgroundDates(fgPoints []Point, year int) []time.Time {
	start := graphStart(year)
	now := time.Now()

	fgSet := make(map[string]bool)
	for _, p := range fgPoints {
		fgSet[fmt.Sprintf("%d-%d", p.Week, p.Day)] = true
	}

	var dates []time.Time
	for week := 0; week < 53; week++ {
		for day := 0; day < 7; day++ {
			if !fgSet[fmt.Sprintf("%d-%d", week, day)] {
				d := start.AddDate(0, 0, week*7+day)
				if d.Year() == year && !d.After(now) {
					dates = append(dates, d)
				}
			}
		}
	}
	return dates
}

func gitInit(path string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	return cmd.Run()
}

func getGitUser() (string, string) {
	name, _ := exec.Command("git", "config", "user.name").Output()
	email, _ := exec.Command("git", "config", "user.email").Output()
	return strings.TrimSpace(string(name)), strings.TrimSpace(string(email))
}

func fastImport(repoPath string, bgDates []time.Time, bgIntensity int, fgDates []time.Time, fgIntensity int, name, email string) error {
	cmd := exec.Command("git", "fast-import", "--quiet")
	cmd.Dir = repoPath

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	total := len(bgDates)*bgIntensity + len(fgDates)*fgIntensity
	count := 0
	var parentMark int

	writeCommit := func(d time.Time) {
		count++
		ts := d.Unix() + int64(count)
		content := fmt.Sprintf("%d\n", ts)
		msg := fmt.Sprintf("draw %d/%d", count, total)

		blobMark := count * 2
		commitMark := count*2 + 1

		fmt.Fprintf(stdin, "blob\nmark :%d\ndata %d\n%s\n", blobMark, len(content), content)
		fmt.Fprintf(stdin, "commit refs/heads/main\nmark :%d\n", commitMark)
		fmt.Fprintf(stdin, "author %s <%s> %d +0000\n", name, email, d.Unix())
		fmt.Fprintf(stdin, "committer %s <%s> %d +0000\n", name, email, d.Unix())
		fmt.Fprintf(stdin, "data %d\n%s\n", len(msg), msg)

		if parentMark > 0 {
			fmt.Fprintf(stdin, "from :%d\n", parentMark)
		}
		fmt.Fprintf(stdin, "M 100644 :%d gitdraw.txt\n\n", blobMark)
		parentMark = commitMark
	}

	for _, d := range bgDates {
		for i := 0; i < bgIntensity; i++ {
			writeCommit(d.Add(time.Duration(i) * time.Hour))
		}
	}

	for _, d := range fgDates {
		for i := 0; i < fgIntensity; i++ {
			writeCommit(d.Add(time.Duration(i) * time.Hour))
		}
	}

	stdin.Close()
	return cmd.Wait()
}

func gitAddRemote(path, url string) error {
	cmd := exec.Command("git", "remote", "add", "origin", url)
	cmd.Dir = path
	return cmd.Run()
}

func gitPush(path string) error {
	branch := exec.Command("git", "branch", "-M", "main")
	branch.Dir = path
	branch.Run()

	push := exec.Command("git", "push", "-u", "origin", "main", "--force")
	push.Dir = path
	return push.Run()
}

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "--cli", "-c":
			runCLI()
			return
		case "--help", "-h":
			printHelpGUI()
			return
		case "--version", "-v":
			printVersionGUI()
			return
		}
	}

	runGUI()
}

func printHelpGUI() {
	fmt.Println(`
  gitdraw — contribution graph art

  Usage:
    gitdraw          Launch GUI (default)
    gitdraw --cli    Run interactive CLI

  Flags:
    -c, --cli      Use command-line interface
    -h, --help     Show help
    -v, --version  Show version
`)
}

func printVersionGUI() {
	fmt.Println("gitdraw v1.0.0")
}

func runGUI() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "GitDraw",
		Width:     900,
		Height:    720,
		MinWidth:  700,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: nil,
		},
		BackgroundColour: &options.RGBA{R: 13, G: 13, B: 13, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Mica,
			Theme:                windows.Dark,
		},
	})

	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}
