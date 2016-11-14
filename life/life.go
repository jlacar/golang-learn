// My mods to the Go implementation of Conway's Game of Life.
//
// based on https://golang.org/doc/play/life.go
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	deadcell = "  "
)

var (
	seed     int64
	gens     int
	skipto   int
	livename string
	livecell []byte
	dingbats map[string]string
)

// Field represents a two-dimensional field of cells.
type Field struct {
	s    [][]bool
	w, h int
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, w)
	}
	return &Field{s: s, w: w, h: h}
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, b bool) {
	f.s[y][x] = b
}

// Alive reports whether the specified cell is alive.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) Alive(x, y int) bool {
	x += f.w
	x %= f.w
	y += f.h
	y %= f.h
	return f.s[y][x]
}

// Next returns the state of the specified cell at the next time step.
func (f *Field) Next(x, y int) bool {
	// Count the adjacent cells that are alive.
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Alive(x+i, y+j) {
				alive++
			}
		}
	}
	// Return next state according to the game rules:
	//   exactly 3 neighbors: on,
	//   exactly 2 neighbors: maintain current state,
	//   otherwise: off.
	return alive == 3 || alive == 2 && f.Alive(x, y)
}

// Life stores the state of a round of Conway's Game of Life.
type Life struct {
	a, b    *Field
	w, h, g int
}

// NewLife returns a new Life game state with a random initial state.
func NewLife(w, h int) *Life {
	a := NewField(w, h)
	for i := 0; i < (w * h / 4); i++ {
		a.Set(rand.Intn(w), rand.Intn(h), true)
	}
	return &Life{
		a: a, b: NewField(w, h),
		w: w, h: h,
	}
}

// Step advances the game by one instant, recomputing and updating all cells.
func (l *Life) step() {
	// Update the state of the next field (b) from the current field (a).
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			l.b.Set(x, y, l.a.Next(x, y))
		}
	}
	// Swap fields a and b.
	l.a, l.b = l.b, l.a

	// increment generation count
	l.g++
}

// String returns the game board as a string.
func (l *Life) String() string {
	var buf bytes.Buffer
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			cell := []byte(deadcell)
			if l.a.Alive(x, y) {
				cell = livecell
			}
			buf.Write(cell)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (l *Life) showGeneration(nth int) {
	fmt.Printf("\n\nGeneration %v (%v of %v):\n\n%v", l.g, nth-skipto, gens, l)
}

func (l *Life) simulate(gens int, delay time.Duration) {

	fmt.Printf("\nConway's Game of Life\n")

	if skipto != 0 {
		fmt.Printf("\nStarting from generation %v...", skipto)
	}

	maxgen := gens + skipto
	skipto--
	for i := 0; i < maxgen; i++ {
		l.step()
		if skipto <= i {
			l.showGeneration(i)
			time.Sleep(delay)
		}
	}

	fmt.Printf("%v generations, %v x %v grid, seed=%v\n\n", l.g, l.h, l.w, seed)
}

func initSeed() {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rand.Seed(seed)
}

func initDisplay() {
	ding, ok := dingbats[livename]
	if !ok {
		ding = dingbats["whitedot"]
	}
	livecell = []byte(" " + ding)
}

func checkSkipping() {
	if skipto < 0 {
		skipto = 0
	}
}

func parseflags() (width, height, perSec int) {

	flag.Int64Var(&seed, "seed", 0,
		"seed for initial population (default random)")

	flag.IntVar(&perSec, "d", 5, "delay 1/`N` seconds between generations")
	flag.IntVar(&height, "h", 30, "height of simulation field")
	flag.IntVar(&width, "w", 30, "width of simulation field")
	flag.IntVar(&gens, "n", 20, "display up to `N` generations")
	flag.IntVar(&skipto, "from", 0, "display from generation `N`")
	flag.StringVar(&livename, "live", "", "`name` of dingbat used to depict a live cell (default whitedot)")

	flag.Parse()

	initSeed()
	initDisplay()
	checkSkipping()

	return
}

func addUsageInfo() {
	// \u263A smile-white
	// \u263B smile-black
	// \u26AA dot-white
	// \u26B1 little-man
	// \u26F9 man-dribble

	dingbats = make(map[string]string)
	dingbats["aster-1"] = "\u2731"
	dingbats["aster-2"] = "\u2749"
	dingbats["bug"] = "\u2603"
	dingbats["circle-x"] = "\u2A02"
	dingbats["dot-star"] = "\u272A"
	dingbats["fat-x"] = "\u2716"
	dingbats["green-x"] = "\u274E"
	dingbats["little-man"] = "\u26B1"
	dingbats["no-entry"] = "\u26D4"
	dingbats["redhat"] = "\u26D1"
	dingbats["skull-x"] = "\u2620"
	dingbats["snowman"] = "\u26C4"
	dingbats["star"] = "\u2606"
	dingbats["whitedot"] = "\u26AA"

	defaultUsage := flag.Usage
	flag.Usage = func() {
		defaultUsage()
		fmt.Fprintf(os.Stderr,
			"\nAvailable dingbats for live cells:\n\n"+
				"Name    \tDescription\n"+
				"--------\t-----------\n"+
				"aster-1 \tAsterisk 1\n"+
				"aster-2 \tAsterisk 2\n"+
				"bug     \tBug\n"+
				"circle-x\tWhite circle with an X\n"+
				"dot-star\tDot with star\n"+
				"fat-x   \tFat white X\n"+
				"green-x \tGreen square with white X\n"+
				"little-man\tLittle yellow man\n"+
				"no-entry\tRed no entry sign\n"+
				"redhat  \tRed hardhat with white cross\n"+
				"skull-x \tSkull and crossbones\n"+
				"snowman \tSnowman\n"+
				"star    \tWhite star\n"+
				"whitedot\tWhite dot (default)\n",
		)
	}
}

func main() {

	addUsageInfo()

	width, height, perSec := parseflags()

	NewLife(width, height).simulate(
		gens,
		time.Second/time.Duration(perSec),
	)
}
