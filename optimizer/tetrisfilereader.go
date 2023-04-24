package optimizer

import (
	"log"
	"os"
	"strings"
)

type Tetrominoe struct {
	squares []Square
	letter  rune
	height  int
	width   int
}

type Square struct {
	x int
	y int
}

func FileReader(filename string) []Tetrominoe {
	filebytes, _ := os.ReadFile(filename)
	filestring := strings.Split(string(filebytes), "\n")
	for i := 0; i < len(filestring); i++ {
		if i%5 == 4 && i != 0 && filestring[i] != "" {
			log.Fatalln("ERROR")
		}
	}
	filestring = strings.Fields(string(filebytes))
	tetrominoes := make([]Tetrominoe, 0)
	squares := make([]Square, 0)
	runeCounter := 65
	sideCounter := 0
	il := 0
	for i := range filestring {
		if len(filestring[i]) != 0 {
			for j := range filestring[i] {
				if filestring[i][j] == '#' {
					squares = append(squares, Square{((i) % 4), (j)})
					sideCounter += AdjacentSquareCounter(filestring, i, j, il)
				}
			}
		}
		if i%4 == 3 && i != 0 {
			if len(squares) != 4 || sideCounter < 6 {
				log.Fatalln("ERROR")
			}
			sideCounter = 0
			tetrominoes = append(tetrominoes, Tetrominoe{squares, rune(runeCounter), 0, 0})
			squares = make([]Square, 0)
			il += 5
			runeCounter++
		}
	}
	tetrominoes = CleanCoords(tetrominoes)
	return tetrominoes
}

func CleanCoords(t []Tetrominoe) []Tetrominoe {
	xc := 0
	xf := false
	yc := 0
	yf := false
	for i := range t {
		for n := 0; n < 4; n++ {
			for _, s := range t[i].squares {
				if s.x == n && !xf {
					xc = n
					xf = true
				}
				if s.y == n && !yf {
					yc = n
					yf = true
				}
			}
		}
		for j := range t[i].squares {
			t[i].squares[j].x = t[i].squares[j].x - xc
			t[i].squares[j].y = t[i].squares[j].y - yc
		}
		xc = 0
		yc = 0
		xf = false
		yf = false
	}
	return t
}

func AdjacentSquareCounter(s []string, i, j, il int) int {
	counter := 0
	length := len(s[i]) - 1
	if i-il < length {
		if s[i+1][j] == '#' {
			counter++
		}
	}
	if j < length {
		if s[i][j+1] == '#' {
			counter++
		}
	}
	if i%4 != 0 {
		if s[i-1][j] == '#' {
			counter++
		}
	}
	if j > 0 {
		if s[i][j-1] == '#' {
			counter++
		}
	}
	return counter
}

func FindHeightAndWidth(t []Tetrominoe) []Tetrominoe {
	for _, v := range t {
		for _, h := range v.squares {
			if h.y > v.height {
				v.height = h.y
			}
			if h.x > v.width {
				v.width = h.x
			}
		}
	}
	return t
}
