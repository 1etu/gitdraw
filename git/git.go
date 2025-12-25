package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Repo struct {
	Path string
}

func Init(path string) (*Repo, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(abs, 0755); err != nil {
		return nil, err
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = abs
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git init: %s", out)
	}

	return &Repo{Path: abs}, nil
}

func (r *Repo) Commit(date time.Time, msg string) error {
	filename := filepath.Join(r.Path, "gitdraw.txt")
	content := fmt.Sprintf("%s\n%s\n", date.Format(time.RFC3339), msg)

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return err
	}

	add := exec.Command("git", "add", ".")
	add.Dir = r.Path
	if out, err := add.CombinedOutput(); err != nil {
		return fmt.Errorf("git add: %s", out)
	}

	dateStr := date.Format("2006-01-02T15:04:05")
	commit := exec.Command("git", "commit", "-m", msg, "--date", dateStr)
	commit.Dir = r.Path
	commit.Env = append(os.Environ(),
		"GIT_AUTHOR_DATE="+dateStr,
		"GIT_COMMITTER_DATE="+dateStr,
	)

	if out, err := commit.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit: %s", out)
	}

	return nil
}

func (r *Repo) FastImport(dates []time.Time, intensity int, progress func(int, int)) error {
	return r.FastImportLayers(nil, 0, dates, intensity, progress)
}

func (r *Repo) FastImportLayers(bgDates []time.Time, bgIntensity int, fgDates []time.Time, fgIntensity int, progress func(int, int)) error {
	name, email := getGitUser()
	if name == "" {
		name = "gitdraw"
	}
	if email == "" {
		email = "gitdraw@local"
	}

	cmd := exec.Command("git", "fast-import", "--quiet")
	cmd.Dir = r.Path

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
		progress(count, total)
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

func getGitUser() (string, string) {
	name, _ := exec.Command("git", "config", "user.name").Output()
	email, _ := exec.Command("git", "config", "user.email").Output()
	return strings.TrimSpace(string(name)), strings.TrimSpace(string(email))
}

func (r *Repo) CommitAll(dates []time.Time) error {
	total := len(dates)
	for i, d := range dates {
		msg := fmt.Sprintf("gitdraw: %d/%d", i+1, total)
		if err := r.Commit(d, msg); err != nil {
			return err
		}
		fmt.Printf("\r[%d/%d] commits", i+1, total)
	}
	fmt.Println()
	return nil
}

func (r *Repo) AddRemote(url string) error {
	cmd := exec.Command("git", "remote", "add", "origin", url)
	cmd.Dir = r.Path
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s", out)
	}
	return nil
}

func (r *Repo) Push() error {
	cmd := exec.Command("git", "branch", "-M", "main")
	cmd.Dir = r.Path
	cmd.CombinedOutput()

	push := exec.Command("git", "push", "-u", "origin", "main")
	push.Dir = r.Path
	if out, err := push.CombinedOutput(); err != nil {
		return fmt.Errorf("%s", out)
	}
	return nil
}
