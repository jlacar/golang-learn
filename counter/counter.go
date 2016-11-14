package main

import "fmt"

type Counter int

func (c *Counter) PostDecr() (v int) {
  v = int(*c)
  *c--
  return
}

func (c *Counter) PreDecr() (v int) {
  *c--
  return int(*c)
}

func (c *Counter) PostIncr() (v int) {
  v = int(*c)
  *c++
  return
}

func (c *Counter) PreIncr() (int) {
  *c++
  return int(*c)
}

func main() {
	var hits Counter
	fmt.Println(hits)
	fmt.Printf("++hits: %v hits++: %v\n", hits.PreIncr(), hits.PostIncr())
	fmt.Println(hits)

	var c Counter

	for c < 5 {
	   fmt.Println(c.PreIncr())
	}
}
