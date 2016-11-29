package main

import (
	"fmt"
)

// fib is a closure that generates the fibonacci series
func fib() func() int {
	var fib0, fib1 = 0, 1
	return func() (f int) {
		f, fib0, fib1 = fib0, fib1, fib0+fib1
		return
	}
}

func main() {
	f := fib()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
