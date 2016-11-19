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
	"strconv"
	"time"
)

type FieldLocation struct {
	X, Y int
}

func NewFieldLocation(x, y int) *FieldLocation {
	return &FieldLocation{X: x, Y: y}
}

func (f *FieldLocation) String() string {
	return fmt.Sprintf("[row:%v, col:%v]", f.Y, f.X)
}

type LocationProvider interface {
	NextLocation() (loc *FieldLocation)
	MoreLocations() bool
}

func AssertMoreLocations(l LocationProvider) {
	if !l.MoreLocations() {
		log.Fatal("Illegal state: no more locations available")
	}
}

var (
	// provides locations for initial population
	Seeder LocationProvider

	// option flag variables
	fieldWidth  int
	fieldHeight int
	seedflag    string
	seed        int64
	gens        int
	gensPerSec  int
	startGen    int
	initPath    string
	iconName    string
)

// RandomLocationProvider provides random FieldLocations that
// fall within bounds of a field with a given width and height
type RandomLocationProvider struct {
	i             int
	width, height int
}

func NewRandomLocationProvider(w, h int) *RandomLocationProvider {
	return &RandomLocationProvider{width: w, height: h}
}

func (r *RandomLocationProvider) NextLocation() (loc *FieldLocation) {
	AssertMoreLocations(r)
	r.i++
	return &FieldLocation{X: rand.Intn(r.width), Y: rand.Intn(r.height)}
}

func (r *RandomLocationProvider) MoreLocations() bool {
	return r.i < r.width*r.height/4
}

// Field represents a two-dimensional field of cells.
type Field struct {
	state         [][]bool
	width, height int
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, w)
	}
	return &Field{state: s, width: w, height: h}
}

// set sets the state of the specified cell to the given value.
func (f *Field) set(loc *FieldLocation, alive bool) {
	if !f.contains(loc) {
		log.Printf("Out of bounds: %v", loc)
		return
	}
	f.state[loc.Y][loc.X] = alive
}

// contains checks if a Field includes a FieldLocation.
// Returns true if the give FieldLocation is within the
// boundaries of the receiving Field
func (f *Field) contains(loc *FieldLocation) bool {
	return loc.X < f.width && loc.Y < f.height
}

// alive reports whether the specified cell is alive.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) alive(x, y int) bool {
	x += f.width
	x %= f.width
	y += f.height
	y %= f.height
	return f.state[y][x] // && !f.BlackHoled(y, x)
}

// next returns the state of the specified cell at the next time step.
func (f *Field) next(x, y int) bool {
	// Count the adjacent cells that are alive.
	neighbors := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.alive(x+i, y+j) {
				neighbors++
			}
		}
	}
	// Return next state according to the game rules:
	//   exactly 3 neighbors: on,
	//   exactly 2 neighbors: maintain current state,
	//   otherwise: off.
	return neighbors == 3 || neighbors == 2 && f.alive(x, y)
}

// Life stores the state of a round of Conway's Game of Life.
type Life struct {
	thisGen, nextGen        *Field
	width, height, genCount int
}

// NewLife returns a new Life game state with initial state provided by Seeder
func NewLife(w, h int) *Life {
	firstGen := NewField(w, h)
	for Seeder.MoreLocations() {
		firstGen.set(Seeder.NextLocation(), true)
	}
	return &Life{
		thisGen: firstGen, nextGen: NewField(w, h),
		width: w, height: h,
	}
}

func (l *Life) prepareNextGeneration() {
	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			l.nextGen.set(&FieldLocation{X: x, Y: y}, l.thisGen.next(x, y))
		}
	}
}

func (l *Life) instateNextGeneration() {
	l.thisGen, l.nextGen = l.nextGen, l.thisGen
	l.genCount++
}

// Step advances the game to the next generation
func (l *Life) step() {
	l.prepareNextGeneration()
	l.instateNextGeneration()
}

// String returns the game board as a string.
func (l *Life) String() string {
	const deadcell = "  "
	var buf bytes.Buffer
	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			cell := []byte(deadcell)
			if l.thisGen.alive(x, y) {
				cell = livecell
			}
			buf.Write(cell)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (l *Life) showCurrentGeneration(nth int) {
	fmt.Printf("\n\nGeneration %v (%v of %v):\n%v", l.genCount+1,
		nth-startGen+1, gens, l)
}

func (l *Life) showRunInfo() {
	fmt.Printf("%v generations calculated.\n\n", l.genCount)
	fmt.Printf("To continue: %v -y %v -x %v %v -icon %v -s %v -n %v\n", os.Args[0],
		l.height, l.width, seedflag, iconName, l.genCount, gens,
	)
}

func (l *Life) stepThroughAll(gens int) {
	delay := time.Second / time.Duration(gensPerSec)
	maxgen := gens + startGen
	for i := 0; i < maxgen; i++ {
		if startGen <= i {
			l.showCurrentGeneration(i)
			time.Sleep(delay)
		}
		l.step()
	}
}

// simulate calculates the specified number of generations
func (l *Life) simulate(gens int) {
	fmt.Printf("\nConway's Game of Life\n")
	l.stepThroughAll(gens)
	l.showRunInfo()
}

func initStartGen() {
	if startGen > 1 {
		fmt.Printf("\nStarting from generation %v...", startGen)
		startGen--
	} else {
		startGen = 0
	}
}

// initSeed initializes the Seeder and seed-related vars
func initSeed() {
	// -f option
	if initPath != "" {
		Seeder = NewFileLocationProvider(initPath)
		seedflag = "-f " + initPath
	}

	// default / fallback
	if Seeder == nil {
		if seed == 0 {
			seed = time.Now().UnixNano()
		}
		rand.Seed(seed)
		Seeder = NewRandomLocationProvider(fieldWidth, fieldHeight)
		seedflag = "-seed " + strconv.FormatInt(seed, 10)
	}
}

var livecell []byte

func initDisplay() {
	s, ok := icon[iconName]
	if !ok {
		iconName = "blue-circle"
		s = icon[iconName]
	}
	livecell = []byte(" " + s)
}

var icon = map[string]string{
	"aster-1":      "\u2731",
	"aster-2":      "\u2749",
	"blue-circle":  "\u23FA",
	"blue-square":  "\u23F9",
	"bug":          "\u2603",
	"circle-plus":  "\u2A01",
	"circle-x":     "\u2A02",
	"dot-star":     "\u272A",
	"fat-x":        "\u2716",
	"flower":       "\u273F",
	"green-x":      "\u274E",
	"man-dribble":  "\u26F9",
	"man-yellow":   "\u26B1",
	"no-entry":     "\u26D4",
	"redhat":       "\u26D1",
	"skull-x":      "\u2620",
	"snowflake":    "\u274A",
	"snowman":      "\u26C4",
	"square-big":   "\u2B1C",
	"square-small": "\u25A9",
	"star-yellow":  "\u2B50",
	"star-white":   "\u2605",
	"star-6pt":     "\u2736",
	"star-8pt":     "\u2738",
	"whitedot":     "\u26AA",
}

func init() {
	flag.Usage = usage

	flag.Int64Var(&seed, "seed", 0,
		"seed for initial population (default random)\n\tignored if -f option specified and valid")

	flag.StringVar(&initPath, "f", "", "read initial population from `filename`\n\tif valid, -seed option is ignored")
	flag.IntVar(&fieldHeight, "y", 30, "height of simulation field")
	flag.IntVar(&fieldWidth, "x", 30, "width of simulation field")
	flag.IntVar(&gens, "n", 20, "display up to `N` generations")
	flag.IntVar(&gensPerSec, "r", 5, "display `N` generations per second")
	flag.IntVar(&startGen, "s", 0, "start displaying from generation `N`")
	flag.StringVar(&iconName, "icon", "", "`name` of icon to use for live cells (default blue-circle)")
}

func usage() {

	fmt.Fprintf(os.Stderr, "Usage: %s [-x] [-y] [-r] [-n] [-s] [-f] [-seed] [-icon]\n\n"+
		"Options:\n\n", os.Args[0])
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

// processArgs processes command line arguments
func processArgs() {
	flag.Parse()

	initSeed()
	initStartGen()
	initDisplay()
}

func main() {
	processArgs()
	NewLife(fieldWidth, fieldHeight).simulate(gens)
}
