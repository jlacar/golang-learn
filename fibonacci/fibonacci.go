package main

import (
	"flag"
	"fmt"
	"os"
)

var n int

// fib returns a closure that generates the fibonacci series
func fib() func() uint64 {
	var fib0, fib1 uint64 = 0, 1
	return func() (f uint64) {
		f, fib0, fib1 = fib0, fib1, fib0+fib1
		return
	}
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [-n]\n\n"+
			"Options:\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&n, "n", 10, "print first `N` numbers of the Fibonacci series")
	flag.Parse()
}

func printSeries(heading string, times int, fn func() uint64) {
	fmt.Printf("\n%v:\n", heading)
	for i := 0; i < times; i++ {
		fmt.Println(fn())
	}
}

func main() {
	f, g := fib(), fib()

	printSeries("First series", n, f)
	printSeries("Second series", n+1, g)
	printSeries("Continue first series", n, f)
	printSeries("Continue second series", n, g)
}
