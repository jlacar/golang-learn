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
	icon     map[string]string
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
	ding, ok := icon[livename]
	if !ok {
		ding = icon["whitedot"]
	}
	livecell = []byte(" " + ding)
}

func checkSkipping() {
	if skipto < 0 {
		skipto = 0
	}
}

func parseflags() (width, height, stepsPerSecond int) {

	flag.Int64Var(&seed, "seed", 0,
		"seed for initial population (default random)")

	flag.IntVar(&height, "y", 30, "height of simulation field")
	flag.IntVar(&width, "x", 30, "width of simulation field")
	flag.IntVar(&gens, "n", 20, "display up to `N` generations")
	flag.IntVar(&stepsPerSecond, "r", 5, "display `N` generations per second")
	flag.IntVar(&skipto, "s", 0, "start displaying from generation `N`")
	flag.StringVar(&livename, "icon", "", "`name` of icon to use for live cells (default whitedot)")

	flag.Parse()

	initSeed()
	initDisplay()
	checkSkipping()

	return
}

func usage() {

	icon = make(map[string]string)
	icon["aster-1"] = "\u2731"
	icon["aster-2"] = "\u2749"
	icon["bug"] = "\u2603"
	icon["circle-x"] = "\u2A02"
	icon["dot-star"] = "\u272A"
	icon["fat-x"] = "\u2716"
	icon["green-x"] = "\u274E"
	icon["man-dribble"] = "\u26F9"
	icon["man-yellow"] = "\u26B1"
	icon["no-entry"] = "\u26D4"
	icon["redhat"] = "\u26D1"
	icon["skull-x"] = "\u2620"
	icon["snowman"] = "\u26C4"
	icon["star"] = "\u2606"
	icon["whitedot"] = "\u26AA"

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-x] [-y] [-r] [-n] [-s] [-seed] [-icon]\n\nOptions:\n\n",
			os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr,
			"\nAvailable icons for live cells:\n\n"+
				"Icon\tName\t\tDescription\n"+
				"----\t--------\t-----------\n"+
				icon["aster-1"]+"\taster-1\t\tAsterisk 1\n"+
				icon["aster-2"]+"\taster-2\t\tAsterisk 2\n"+
				icon["bug"]+"\tbug\t\tBug\n"+
				icon["circle-x"]+"\tcircle-x\tCircle with an X\n"+
				icon["dot-star"]+"\tdot-star\tDot with star\n"+
				icon["fat-x"]+"\tfat-x\t\tFat white X\n"+
				icon["green-x"]+"\tgreen-x\t\tGreen square with white X\n"+
				icon["man-dribble"]+"\tman-dribble\tMan dribbling ball\n"+
				icon["man-yellow"]+"\tman-yellow\tLittle yellow man\n"+
				icon["no-entry"]+"\tno-entry\tNo entry sign\n"+
				icon["redhat"]+"\tredhat\t\tRed hardhat with white cross\n"+
				icon["skull-x"]+"\tskull-x\t\tSkull and crossbones\n"+
				icon["snowman"]+"\tsnowman\t\tSnowman\n"+
				icon["star"]+"\tstar\t\tStar\n"+
				icon["whitedot"]+"\twhitedot\tWhite dot (default)\n",
		)
	}
}

func main() {

	usage()

	width, height, stepsPerSecond := parseflags()

	NewLife(width, height).simulate(
		gens,
		time.Second/time.Duration(stepsPerSecond),
	)
}
