package draw

import (
	"strings"
	"time"

	"github.com/1etu/gitdraw/font"
)

const (
	Rows  = 7
	Weeks = 53
)

type Grid [Rows][Weeks]int

type Point struct {
	Week int
	Day  int
}

func Text(text string) Grid {
	var grid Grid
	col := 1

	for _, ch := range text {
		if col >= Weeks {
			break
		}

		glyph := font.Get(ch)
		for x := 0; x < font.Width(); x++ {
			if col+x >= Weeks {
				break
			}
			for y := 0; y < font.Height(); y++ {
				if glyph[y]&(1<<x) != 0 {
					grid[y][col+x] = 1
				}
			}
		}
		col += font.Width() + 1
	}

	return grid
}

func (g Grid) Points() []Point {
	var pts []Point
	for week := 0; week < Weeks; week++ {
		for day := 0; day < Rows; day++ {
			if g[day][week] > 0 {
				pts = append(pts, Point{Week: week, Day: day})
			}
		}
	}
	return pts
}

func graphStart(year int) time.Time {
	dec31 := time.Date(year, 12, 31, 12, 0, 0, 0, time.UTC)
	
	daysSinceSunday := int(dec31.Weekday())
	lastSunday := dec31.AddDate(0, 0, -daysSinceSunday)
	
	return lastSunday.AddDate(0, 0, -52*7)
}

func (g Grid) Dates(year int) []time.Time {
	start := graphStart(year)
	now := time.Now()
	pts := g.Points()
	dates := make([]time.Time, 0, len(pts))

	for _, p := range pts {
		d := start.AddDate(0, 0, p.Week*7+p.Day)
		if d.Year() == year && !d.After(now) {
			dates = append(dates, d)
		}
	}
	return dates
}

func (g Grid) BackgroundDates(year int) []time.Time {
	start := graphStart(year)
	now := time.Now()
	var dates []time.Time

	for week := 0; week < Weeks; week++ {
		for day := 0; day < Rows; day++ {
			if g[day][week] == 0 {
				d := start.AddDate(0, 0, week*7+day)
				if d.Year() == year && !d.After(now) {
					dates = append(dates, d)
				}
			}
		}
	}
	return dates
}

func AllDates(year int) []time.Time {
	start := time.Date(year, 1, 1, 12, 0, 0, 0, time.UTC)
	end := time.Date(year, 12, 31, 12, 0, 0, 0, time.UTC)
	now := time.Now()

	var dates []time.Time
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if !d.After(now) {
			dates = append(dates, d)
		}
	}
	return dates
}

func (g Grid) Render() string {
	var sb strings.Builder
	days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	for row := 0; row < Rows; row++ {
		sb.WriteString(days[row])
		sb.WriteString(" ")
		for col := 0; col < Weeks; col++ {
			if g[row][col] > 0 {
				sb.WriteString("██")
			} else {
				sb.WriteString("░░")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
