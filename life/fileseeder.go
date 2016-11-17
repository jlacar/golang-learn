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
		morelocs, lastRow := parseFieldConfig(l, row)
		row = lastRow
		if morelocs != nil {
			locs = append(locs, morelocs...)
		}
	}
	return &FileLocationProvider{w: w, h: h, locs: locs}
}

func parseFieldConfig(configstr string, lastRow int) (locs []FieldLocation, row int) {
	if strings.IndexRune(configstr, '#') == 0 {
		log.Println(configstr)
		return nil, lastRow
	}

	// check for valid format (row : cells config)
	parts := strings.Split(configstr, ":")
	if len(parts) != 2 {
		log.Printf("Line ignored: [%v]\n", configstr)
		return nil, lastRow
	}

	// check for NN:... or ++:...
	y, err := strconv.Atoi(parts[0])
	if parts[0] == "++" {
		y = lastRow + 1
	} else if err != nil {
		log.Printf("Invalid row number in field file: [%v]\n", parts[0])
		return nil, lastRow
	}

	cols := parseFieldConfigColumns(parts[1])
	if len(cols) == 0 {
		log.Printf("Empty line ignored")
		return nil, y
	}

	row = y
	locs = make([]FieldLocation, len(cols))
	for i, x := range cols {
		locs[i] = FieldLocation{Y: y, X: x}
	}
	return
}

func parseFieldConfigColumns(configstr string) []int {
	cols := []int{}
	chars := strings.Split(configstr, "")
	for x, ch := range chars {
		if ch != " " {
			cols = append(cols, x)
		}
	}
	if len(cols) != 0 {
		return cols
	}
	return nil
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
