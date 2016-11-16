// My mods to the Go implementation of Conway's Game of Life.
//
// based on https://golang.org/doc/play/life.go
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	deadcell = "  "
)

type FieldLocation struct {
	X, Y int
}

type LocationSource interface {
	NextLocation() (loc *FieldLocation)
	HasNext() bool
}

var (
	Seeder LocationSource

	seed     int64
	gens     int
	startGen int
	initPath string
	iconName string
	livecell []byte
	icon     map[string]string
)

type RandomLocationProvider struct {
	i    int
	w, h int
}

func assertHasNext(l LocationSource) {
	if !l.HasNext() {
		log.Fatal("Illegal state: no more locations available")
	}
}

func (r *RandomLocationProvider) NextLocation() (loc *FieldLocation) {
	assertHasNext(r)
	r.i++
	return &FieldLocation{X: rand.Intn(r.w), Y: rand.Intn(r.h)}
}

func (r *RandomLocationProvider) HasNext() bool {
	return r.i < r.w*r.h/4
}

func NewRandomLocationProvider(w, h int) *RandomLocationProvider {
	return &RandomLocationProvider{w: w, h: h}
}

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
func (f *Field) Set(loc *FieldLocation, alive bool) {
	f.s[loc.Y][loc.X] = alive
}

// Alive reports whether the specified cell is alive.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) Alive(x, y int) bool {
	x += f.w
	x %= f.w
	y += f.h
	y %= f.h
	return f.s[y][x] // && !f.BlackHoled(y, x)
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

// NewLife returns a new Life game state with initial state provided by Seeder
func NewLife(w, h int) *Life {
	a := NewField(w, h)
	for Seeder.HasNext() {
		a.Set(Seeder.NextLocation(), true)
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
			l.b.Set(&FieldLocation{X: x, Y: y}, l.a.Next(x, y))
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
	fmt.Printf("\n\nGeneration %v (%v of %v):\n\n%v", l.g+1, nth-startGen+1, gens, l)
}

func (l *Life) showSummary() {
	fmt.Printf("%v generations calculated.\n\n", l.g)
	fmt.Printf("To continue: %v -y %v -x %v -seed %v -icon %v -s %v -n %v\n", os.Args[0],
		l.h, l.w, seed, iconName, l.g, gens,
	)
}

func (l *Life) simulate(gens int, delay time.Duration) {

	fmt.Printf("\nConway's Game of Life\n")

	if startGen > 1 {
		fmt.Printf("\nStarting from generation %v...", startGen)
		startGen--
	} else {
		startGen = 0
	}

	maxgen := gens + startGen
	for i := 0; i < maxgen; i++ {
		if startGen <= i {
			l.showGeneration(i)
			time.Sleep(delay)
		}
		l.step()
	}

	l.showSummary()
}

func initSeed(w, h int) {
	// check for initPath option
	if initPath != "" {
		Seeder = NewFileLocationSource(initPath, w, h)
	}

	// default to random location seeder
	if Seeder == nil {
		if seed == 0 {
			seed = time.Now().UnixNano()
		}
		rand.Seed(seed)
		Seeder = NewRandomLocationProvider(w, h)
	}
}

func initDisplay() {
	ding, ok := icon[iconName]
	if !ok {
		iconName = "blue-circle"
		ding = icon[iconName]
	}
	livecell = []byte(" " + ding)
}

func checkStartGeneration() {
	if startGen < 0 {
		startGen = 0
	}
}

func parseflags() (width, height, stepsPerSecond int) {

	flag.Int64Var(&seed, "seed", 0,
		"seed for initial population (default random)")

	flag.StringVar(&initPath, "f", "", "seed population from `filename`")
	flag.IntVar(&height, "y", 30, "height of simulation field")
	flag.IntVar(&width, "x", 30, "width of simulation field")
	flag.IntVar(&gens, "n", 20, "display up to `N` generations")
	flag.IntVar(&stepsPerSecond, "r", 5, "display `N` generations per second")
	flag.IntVar(&startGen, "s", 0, "start displaying from generation `N`")
	flag.StringVar(&iconName, "icon", "", "`name` of icon to use for live cells (default blue-circle)")

	flag.Parse()

	initSeed(width, height)
	initDisplay()
	checkStartGeneration()

	return
}

func usage() {

	icon = make(map[string]string)
	icon["aster-1"] = "\u2731"
	icon["aster-2"] = "\u2749"
	icon["blue-circle"] = "\u23FA"
	icon["blue-square"] = "\u23F9"
	icon["bug"] = "\u2603"
	icon["circle-plus"] = "\u2A01"
	icon["circle-x"] = "\u2A02"
	icon["dot-star"] = "\u272A"
	icon["fat-x"] = "\u2716"
	icon["flower"] = "\u273F"
	icon["green-x"] = "\u274E"
	icon["man-dribble"] = "\u26F9"
	icon["man-yellow"] = "\u26B1"
	icon["no-entry"] = "\u26D4"
	icon["redhat"] = "\u26D1"
	icon["skull-x"] = "\u2620"
	icon["snowflake"] = "\u274A"
	icon["snowman"] = "\u26C4"
	icon["square-big"] = "\u2B1C"
	icon["square-small"] = "\u25A9"
	icon["star-yellow"] = "\u2B50"
	icon["star-white"] = "\u2605"
	icon["star-6pt"] = "\u2736"
	icon["star-8pt"] = "\u2738"
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
				icon["blue-circle"]+"\tblue-circle\tBlue tile, white circle (default)\n"+
				icon["blue-square"]+"\tblue-square\tBlue tile, white square\n"+
				icon["bug"]+"\tbug\t\tBug\n"+
				icon["circle-plus"]+"\tcircle-plus\tCircle with a '+'\n"+
				icon["circle-x"]+"\tcircle-x\tCircle with an 'x'\n"+
				icon["dot-star"]+"\tdot-star\tDot with star\n"+
				icon["fat-x"]+"\tfat-x\t\tFat white X\n"+
				icon["flower"]+"\tflower\t\tFlower\n"+
				icon["green-x"]+"\tgreen-x\t\tGreen tile with white X\n"+
				icon["man-dribble"]+"\tman-dribble\tMan dribbling ball\n"+
				icon["man-yellow"]+"\tman-yellow\tLittle yellow man\n"+
				icon["no-entry"]+"\tno-entry\tNo entry sign\n"+
				icon["redhat"]+"\tredhat\t\tRed hardhat with white cross\n"+
				icon["skull-x"]+"\tskull-x\t\tSkull and crossbones\n"+
				icon["snowflake"]+"\tsnowflake\tSnowflake\n"+
				icon["snowman"]+"\tsnowman\t\tSnowman\n"+
				icon["square-big"]+"\tsquare-big\tBig square\n"+
				icon["square-small"]+"\tsquare-small\tSmall square\n"+
				icon["star-yellow"]+"\tstar-yellow\tYellow 5-point star\n"+
				icon["star-white"]+"\tstar-white\tWhite 5-point star\n"+
				icon["star-6pt"]+"\tstar-6pt\t6-point star\n"+
				icon["star-8pt"]+"\tstar-8pt\t8-point star\n"+
				icon["whitedot"]+"\twhitedot\tWhite dot\n",
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
