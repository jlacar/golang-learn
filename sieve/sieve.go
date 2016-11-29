// A Go implementation of The Sieve of Eratosthenes
package main

import (
	"fmt"
	"os"
	"strconv"
)

var primes []bool

func findPrimes(max int) {
	primes = make([]bool, max+1)

	for i := 2; i < len(primes); i++ {
		primes[i] = true
	}

	for i := 2; i < len(primes); i++ {
		if primes[i] {
			for j := 2 * i; j < len(primes); j += i {
				primes[j] = false
			}
		}
	}
}

func listPrimes() {
	count := 0
	for i := 0; i < len(primes); i++ {
		if primes[i] {
			fmt.Printf("%4v, ", i)
			count++
			if count == 20 {
				fmt.Print("\n")
				count = 0
			}
		}
	}
	fmt.Print("\n")
}

func main() {
	max, _ := strconv.Atoi(os.Args[1])

	findPrimes(max)
	listPrimes()
}
