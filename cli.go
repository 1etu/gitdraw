//go:build !gui

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/1etu/gitdraw/draw"
	"github.com/1etu/gitdraw/git"
)

const (
	reset  = "\033[0m"
	dim    = "\033[2m"
	bold   = "\033[1m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
)

var reader = bufio.NewReader(os.Stdin)

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "--gui", "-g":
			fmt.Println("\n  This build does not include GUI support.")
			fmt.Println("  Build with: wails build -tags gui")
			fmt.Println()
			os.Exit(1)
		case "--help", "-h":
			printHelp()
			return
		case "--version", "-v":
			printVersion()
			return
		}
	}

	runCLI()
}

func printHelp() {
	fmt.Println(`
  gitdraw — contribution graph art

  Usage:
    gitdraw          Run interactive CLI
    gitdraw --help   Show this help

  Build with GUI:
    wails build -tags gui

  Flags:
    -h, --help     Show help
    -v, --version  Show version
`)
}

func printVersion() {
	fmt.Println("gitdraw v1.0.0")
}

func runCLI() {
	clearScreen()
	printHeader()

	text := ask("Text to draw")
	if text == "" {
		exit("text cannot be empty")
	}

	grid := draw.Text(strings.ToUpper(text))

	fmt.Println()
	printPreview(grid)

	if !confirm("Continue with this design") {
		fmt.Println()
		fmt.Println(dim + "Cancelled." + reset)
		return
	}

	year := askWithDefault("Target year", fmt.Sprintf("%d", time.Now().Year()))
	var yearInt int
	fmt.Sscanf(year, "%d", &yearInt)
	if yearInt < 2008 || yearInt > 2099 {
		yearInt = time.Now().Year()
	}

	dates := grid.Dates(yearInt)

	fillMode := confirm("Fill background? (creates contrast)")

	var bgDates []time.Time
	var bgIntensity int

	if fillMode {
		bgDates = grid.BackgroundDates(yearInt)
		bgIntensity = 1
	}

	intensity := askWithDefault("Text intensity", "15")
	var intensityInt int
	fmt.Sscanf(intensity, "%d", &intensityInt)
	if intensityInt < 1 {
		intensityInt = 1
	}
	if intensityInt > 50 {
		intensityInt = 50
	}

	totalCommits := len(dates)*intensityInt + len(bgDates)*bgIntensity

	fmt.Println()
	info("text pixels", fmt.Sprintf("%d", len(dates)))
	if fillMode {
		info("background pixels", fmt.Sprintf("%d", len(bgDates)))
	}
	info("total commits", fmt.Sprintf("%d", totalCommits))
	info("target year", fmt.Sprintf("%d", yearInt))

	repoPath := askWithDefault("Output directory", "gitdraw-repo")

	if dirExists(repoPath) {
		fmt.Println()
		warn("directory already exists: " + repoPath)
		if !confirm("Overwrite") {
			fmt.Println(dim + "Cancelled." + reset)
			return
		}
		os.RemoveAll(repoPath)
	}

	fmt.Println()
	spin("Initializing repository", func() error {
		_, err := git.Init(repoPath)
		return err
	})

	repo := &git.Repo{Path: repoPath}

	fmt.Println()
	if !confirm(fmt.Sprintf("Generate %d commits", totalCommits)) {
		fmt.Println(dim + "Repository created but empty." + reset)
		return
	}

	fmt.Println()
	width := 40
	var importErr error
	if fillMode {
		importErr = repo.FastImportLayers(bgDates, bgIntensity, dates, intensityInt, func(done, total int) {
			pct := float64(done) / float64(total)
			filled := int(pct * float64(width))
			bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
			fmt.Printf("\r  %s%s%s %s%3d%%%s", cyan, bar, reset, dim, int(pct*100), reset)
		})
	} else {
		importErr = repo.FastImport(dates, intensityInt, func(done, total int) {
			pct := float64(done) / float64(total)
			filled := int(pct * float64(width))
			bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
			fmt.Printf("\r  %s%s%s %s%3d%%%s", cyan, bar, reset, dim, int(pct*100), reset)
		})
	}
	fmt.Println()

	if importErr != nil {
		fmt.Printf("\n  %s!%s %s\n", bold+yellow, reset, "commit generation failed")
		return
	}

	fmt.Println()
	success("Repository ready")
	fmt.Println()

	if confirm("Configure GitHub remote") {
		configureRemote(repo, repoPath)
	} else {
		printNextSteps(repoPath)
	}
}

func configureRemote(repo *git.Repo, repoPath string) {
	fmt.Println()
	fmt.Println(dim + "  Create an empty repo at github.com/new (no README)" + reset)
	fmt.Println()

	remote := ask("GitHub URL (or enter to skip)")
	if remote == "" {
		printNextSteps(repoPath)
		return
	}

	if !strings.Contains(remote, "://") && !strings.HasPrefix(remote, "git@") {
		remote = "https://" + remote
	}

	if err := repo.AddRemote(remote); err != nil {
		warn("failed to add remote")
		printNextSteps(repoPath)
		return
	}

	fmt.Println()
	spin("Pushing", func() error {
		return repo.Push()
	})

	fmt.Println()
	success("Done")
	fmt.Println()
}

func printNextSteps(repoPath string) {
	fmt.Println()
	fmt.Println(dim + "  To push manually:" + reset)
	fmt.Println()
	fmt.Printf("    cd %s\n", repoPath)
	fmt.Println("    git remote add origin <url>")
	fmt.Println("    git push -u origin main")
	fmt.Println()
}

func printHeader() {
	fmt.Println()
	fmt.Println(bold + "  gitdraw" + reset + dim + " — contribution graph art" + reset)
	fmt.Println()
}

func printPreview(grid draw.Grid) {
	fmt.Println(dim + "  Preview:" + reset)
	fmt.Println()
	for _, line := range strings.Split(grid.Render(), "\n") {
		if line != "" {
			fmt.Println("  " + line)
		}
	}
}

func ask(prompt string) string {
	fmt.Printf("  %s%s%s ", bold, prompt, reset)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func askWithDefault(prompt, def string) string {
	fmt.Printf("  %s%s%s %s(%s)%s ", bold, prompt, reset, dim, def, reset)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return def
	}
	return text
}

func confirm(prompt string) bool {
	fmt.Printf("  %s%s?%s %s(y/n)%s ", bold, prompt, reset, dim, reset)
	text, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(text)) == "y"
}

func info(label, value string) {
	fmt.Printf("  %s%s:%s %s\n", dim, label, reset, value)
}

func success(msg string) {
	fmt.Printf("  %s%s✓%s %s\n", bold, green, reset, msg)
}

func warn(msg string) {
	fmt.Printf("  %s%s!%s %s\n", bold, yellow, reset, msg)
}

func exit(msg string) {
	fmt.Println()
	fmt.Printf("  %serror:%s %s\n", bold, reset, msg)
	fmt.Println()
	os.Exit(1)
}

func spin(msg string, fn func() error) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	done := make(chan error)

	go func() {
		done <- fn()
	}()

	i := 0
	for {
		select {
		case err := <-done:
			if err != nil {
				fmt.Printf("\r  %s%s✗%s %s\n", bold, yellow, reset, msg)
				return
			}
			fmt.Printf("\r  %s%s✓%s %s\n", bold, green, reset, msg)
			return
		default:
			fmt.Printf("\r  %s%s%s %s", cyan, frames[i%len(frames)], reset, msg)
			time.Sleep(80 * time.Millisecond)
			i++
		}
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
