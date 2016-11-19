package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type FileLocationProvider struct {
	w, h, i int

	locs []FieldLocation
}

func (f *FileLocationProvider) NewLocation() (loc *FieldLocation) {
	assertMoreLocations(f)
	loc = &FieldLocation{Y: f.locs[f.i].Y, X: f.locs[f.i].X}
	f.i++
	return
}

func (f *FileLocationProvider) MoreLocations() bool {
	return f.i < len(f.locs)
}

func NewFileLocationProvider(path string, w, h int) *FileLocationProvider {
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
	return &FileLocationProvider{w: w, h: h, locs: locs}
}

var columnOffset int // added to X of new FieldLocations

func parseConfigLine(configline string, lastRow int) (locs []FieldLocation, row int) {
	if strings.IndexRune(configline, '#') == 0 {
		log.Println(configline)
		return nil, lastRow
	}

	// check for valid format (NN|++|>>: ...)
	parts := strings.Split(configline, ":")
	if len(parts) != 2 {
		log.Println(configline)
		return nil, lastRow
	}

	header, data := parts[0], parts[1]

	// >>:NN -- offset column to NN
	if header == ">>" {
		co, err := strconv.Atoi(data)
		if err == nil {
			log.Printf(">> [%v]", co)
			columnOffset = co
		} else {
			log.Println(configline)
		}
		return nil, lastRow
	}

	y, err := strconv.Atoi(header)

	// ++: -- use lastRow + 1
	if header == "++" {
		y = lastRow + 1
	} else if err != nil {
		log.Println(configline)
		return nil, lastRow
	}

	// NN: ...
	cols := parseConfigLineMarkings(data)
	if len(cols) == 0 {
		log.Println(configline)
		return nil, y
	}

	row = y
	locs = make([]FieldLocation, len(cols))
	for i, x := range cols {
		locs[i] = FieldLocation{Y: y, X: x + columnOffset}
	}
	return
}

// parseFieldConfigColumns parses one line read from the
// configuration file.
func parseConfigLineMarkings(markedLine string) []int {
	cols := []int{}
	markings := strings.Split(markedLine, "")
	for x, mark := range markings {
		if mark != " " {
			cols = append(cols, x)
		}
	}
	return cols
}

// readLines reads a whole file into memory
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
