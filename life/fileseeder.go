package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// FileLocationProvider is a LocationProvider implementation that
// uses a field definition file as the source for live cell locations.
type FileLocationProvider struct {
	path             string
	i, width, height int
	locs             []FieldLocation
}

// NextLocation returns the next FieldLocation read from the file
func (f *FileLocationProvider) NextLocation() (loc *FieldLocation) {
	loc = &f.locs[f.i]
	f.i++
	return
}

// MoreLocations returns true if there are more FieldLocations available
func (f FileLocationProvider) MoreLocations() bool {
	return f.i < len(f.locs)
}

// MinimumBounds reports the minumum width and height of a field that
// can accomodate all the FieldLocations that will be provided.
func (f FileLocationProvider) MinimumBounds() (width, height int) {
	return f.width, f.height
}

func (f FileLocationProvider) String() string {
	return fmt.Sprintf("FileLocationProvider: file: %v minX: %v, minY: %v", f.path, f.width, f.height)
}

// NewFileLocationProvider creates a FileLocationProvider that gets its
// its FieldLocations from the field definition file specified by path.
func NewFileLocationProvider(path string) (*FileLocationProvider, error) {
	lines, err := readLines(path)

	if err != nil {
		log.Println(err.Error())
		return nil, fmt.Errorf("Could not read file [%v]", path)
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("File [%v] is empty", path)
	}

	locs := []FieldLocation{}
	var minX, minY int
	row := 0
	for _, l := range lines {
		morelocs, lastrow := parseConfigLine(l, row)
		row = lastrow
		if len(morelocs) != 0 {
			locs = append(locs, morelocs...)
		}
		minY = max(minY, row)
		minX = maxCol(minX, locs)
	}

	return &FileLocationProvider{path: path, locs: locs, width: minX + 1, height: minY + 1}, nil
}

func maxCol(x int, locs []FieldLocation) (max int) {
	max = x
	for _, l := range locs {
		if l.X > max {
			max = l.X
		}
	}
	return
}

// ignorable checks if the given configuration line can be ignored for parsing
// and returns true if it starts with "#" or does not contain ":".
func ignorable(configline string) bool {
	if strings.HasPrefix(configline, "#") || !strings.Contains(configline, ":") {
		log.Println(configline)
		return true
	}
	return false
}

var columnOffset int // added to relative column #s to get absolute #s

// parseConfigLine parses a line from a field configuration file
// and returns a slice of FieldLocations and the field row that these
// FieldLocations are on. The returned row number is used to update the
// baseline used with the relative row number header, "++".
func parseConfigLine(configline string, lastRow int) (locs []FieldLocation, row int) {

	if ignorable(configline) {
		return nil, lastRow
	}

	// separate line header from settings
	parts := strings.Split(configline, ":")
	header, settings := parts[0], strings.TrimRightFunc(parts[1], unicode.IsSpace)

	// >>:NN -- set column offset to NN
	if header == ">>" {
		co, err := strconv.Atoi(settings)
		if err == nil && co >= 0 {
			log.Printf(">> [%v]", co)
			columnOffset = co
		} else {
			log.Println(configline)
		}
		return nil, lastRow
	}

	y, err := strconv.Atoi(header)

	// ++: -- use relative row number
	if header == "++" {
		y = lastRow + 1
	} else if err != nil {
		log.Println(configline)
		return nil, lastRow
	}

	// NN: ...
	cols := parseConfigLineSettings(settings)
	if len(cols) == 0 {
		log.Println(configline)
		return nil, y
	}

	return toFieldLocations(cols, y), y
}

// toFieldLocations maps the row # and relative column #s
// and returns a slice of FieldLocations
func toFieldLocations(cols []int, y int) (locs []FieldLocation) {
	locs = make([]FieldLocation, len(cols))
	for i, x := range cols {
		locs[i] = *NewFieldLocation(x+columnOffset, y)
	}
	return
}

// parseConfigLineSettings parses the marks in a line
// and returns a slice of column #s of live cells.
// The #s are relative to the start of the given string.
// A space marks a dead cell; anything else marks a live cell.
func parseConfigLineSettings(settings string) []int {
	cols := []int{}
	markings := strings.Split(settings, "")
	for x, mark := range markings {
		if mark != " " {
			cols = append(cols, x)
		}
	}
	return cols
}

// readLines reads a field configuration file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
