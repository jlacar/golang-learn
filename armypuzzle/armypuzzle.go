// see https://coderanch.com/t/673317/java/Java-Regiment-Army-Class

package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Regiment struct {
	name             string
	number, strength int
}

type Army struct {
	regiments []*Regiment
}

func (a *Army) solve() {
	reportRegimentStatus(a.regiments)

	weekRegiment5goes := 0
	for week := 1; week <= 20; week++ {
		a.update()
		pos, biggest := a.biggestRegiment()
		a.shipout(pos)

		reportWeekStatus(week, biggest)
		reportRegimentStatus(a.regiments)

		if biggest.number == 5 {
			weekRegiment5goes = week
		}
	}
	fmt.Printf("\nAnswer: Regiment 5 waits %v weeks to ship out\n", weekRegiment5goes)
}

func (a *Army) shipout(r int) {
	a.regiments = append(a.regiments[:r], a.regiments[r+1:]...)
}

func reportWeekStatus(w int, shippedOut *Regiment) {
	fmt.Printf("\nWeek %d\n", w)
	fmt.Printf("Regiment %v (%v) with %v men shipped out\n", shippedOut.number,
		shippedOut.name, shippedOut.strength)
}

func reportRegimentStatus(regiments []*Regiment) {
	const format = "%3v  %-15s %5v\n"
	fmt.Printf("\nRegiment status (%v available)\n\n", len(regiments))
	fmt.Printf(format, "#", "Name", "Men")
	total := 0
	for _, r := range regiments {
		fmt.Printf(format, r.number, r.name, r.strength)
		total += r.strength
	}
	if len(regiments) == 0 {
		fmt.Printf(format, "-", "(none)", "-")
	} else {
		fmt.Printf(format, "", "TOTAL", total)
	}
}

func (a *Army) update() {
	for _, r := range a.regiments {
		if r.number == 5 {
			r.strength += 30
		} else {
			r.strength += 100
		}
	}
}

func (a *Army) biggestRegiment() (pos int, mostMen *Regiment) {
	pos, mostMen = 0, a.regiments[0]
	for i, r := range a.regiments {
		if r.strength > mostMen.strength {
			pos = i
			mostMen = r
		}
	}
	return
}

func NewArmy(regimentList []string) *Army {
	strength := 50 * len(regimentList)
	regs := make([]*Regiment, len(regimentList))
	for i, s := range regimentList {
		parts := strings.Split(s, " ")
		num, _ := strconv.Atoi(parts[0])
		regs[i] = &Regiment{number: num, name: parts[1], strength: strength}
		strength -= 50
	}
	return &Army{regiments: regs}
}

func main() {
	army := NewArmy([]string{
		"1 Aardvarks",
		"2 Begonias",
		"3 Chrysanthemums",
		"4 Dhalias",
		"5 Elephants",
		"6 Ferrets",
		"7 GilaMonsters",
		"8 Hyraxes",
		"9 Ibex",
		"10 Jackyls",
		"11 KimodoDragons",
		"12 Lemurs",
		"13 Marigolds",
		"14 Nonames",
		"15 Opossums",
		"16 Porcupines",
		"17 Quahogs",
		"18 Rhododendrons",
		"19 Swordfish",
		"20 Tapirs",
	})
	army.solve()
}
