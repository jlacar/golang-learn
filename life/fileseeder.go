package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type FileLocationSource struct {
	w, h, i int

	locs []FieldLocation
}

func (f *FileLocationSource) NextLocation() (loc *FieldLocation) {
	assertHasNext(f)
	loc = &FieldLocation{Y: f.locs[f.i].Y, X: f.locs[f.i].X}
	f.i++
	return
}

func (f *FileLocationSource) HasNext() bool {
	return f.i < len(f.locs)
}

func NewFileLocationSource(path string, w, h int) *FileLocationSource {
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
	for _, l := range lines {
		morelocs := parseFieldConfig(l)
		if morelocs != nil {
			locs = append(locs, morelocs...)
		}
	}
	return &FileLocationSource{w: w, h: h, locs: locs}
}

func parseFieldConfig(configstr string) (locs []FieldLocation) {
	parts := strings.Split(configstr, ":")
	if len(parts) != 2 {
		return nil
	}
	cols := parseFieldConfigColumns(parts[1])
	if len(cols) == 0 {
		return nil
	}
	locs = make([]FieldLocation, len(cols))
	y, _ := strconv.Atoi(parts[0])
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
