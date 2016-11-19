package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// FileLocationProvider is a LocationProvider implementation that
// reads a text file and parses it for live cell locations. The
// text file format expected is as follows:
// - A line that starts with "#" or does not contain at least one ":"
//   is ignored as a comment
// - A line that starts with ">>:" is interpreted as a column offset
//   setting. A number greater than or equal to 0 must be specified
//   after the ":" and sets the column offset for converting relative
//   column numbers to absolute for all subsequent lines parsed
// - A line that starts with "NN:" is assigned that number as its row
// - A line that starts with "++:" is assigned the row number of the previous
//   line + 1

type FileLocationProvider struct {
	i    int
	locs []FieldLocation
}

// NextLocation returns the next FieldLocation available from the receiving provider.
func (f *FileLocationProvider) NextLocation() (loc *FieldLocation) {
	AssertMoreLocations(f)
	loc = &f.locs[f.i]
	f.i++
	return
}

// MoreLocations returns true if there are more FieldLocations available from the receiving provider.
func (f *FileLocationProvider) MoreLocations() bool {
	return f.i < len(f.locs)
}

// NewFileLocationProvider creates a FileLocationProvider with the given
// height and width and reads a field configuration from the file
// specified by path.
func NewFileLocationProvider(path string) *FileLocationProvider {
	lines, err := readLines(path)
	if err != nil {
		log.Printf("Could not read file [%v]: %v", path, err.Error())
		return nil
	}
	if len(lines) == 0 {
		log.Printf("File [%v] is empty", path)
		return nil
	}
	locs := []FieldLocation{}
	row := -1
	for _, l := range lines {
		morelocs, lastRow := parseConfigLine(l, row)
		row = lastRow
		if len(morelocs) != 0 {
			locs = append(locs, morelocs...)
		}
	}
	return &FileLocationProvider{locs: locs}
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
