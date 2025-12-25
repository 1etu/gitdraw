// DO NOT RUN THIS

package font

import (
	"fmt"
	"testing"
)

func TestRenderAllGlyphs(t *testing.T) {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fmt.Println("\n=== UPPERCASE ===")
	for _, ch := range chars {
		renderChar(rune(ch))
	}

	nums := "0123456789"
	fmt.Println("\n=== NUMBERS ===")
	for _, ch := range nums {
		renderChar(rune(ch))
	}

	symbols := " !.-_:/<>♥"
	fmt.Println("\n=== SYMBOLS ===")
	for _, ch := range symbols {
		renderChar(rune(ch))
	}
}

func renderChar(ch rune) {
	g := Get(ch)
	fmt.Printf("\n'%c':\n", ch)
	for row := 0; row < Height(); row++ {
		for col := 0; col < Width(); col++ {
			if g[row]&(1<<col) != 0 {
				fmt.Print("██")
			} else {
				fmt.Print("░░")
			}
		}
		fmt.Println()
	}
}
