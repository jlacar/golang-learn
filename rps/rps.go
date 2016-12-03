package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Move int

const (
	ROCK Move = iota
	SPOCK
	PAPER
	LIZARD
	SCISSORS
	LEN_Move int = iota
)

type MatchUp struct {
	p1, p2 Move
	w, l   string
}

/*
http://bigbangtheory.wikia.com/wiki/Rock_Paper_Scissors_Lizard_Spock

Scissors cut Paper
Paper covers Rock
Rock crushes Lizard
Lizard poisons Spock
Spock smashes Scissors
Scissors decapitates Lizard
Lizard eats Paper
Paper disproves Spock
Spock vaporizes Rock
(and as it always has) Rock crushes Scissors

*/

var pairings = []*MatchUp{
	&MatchUp{SCISSORS, PAPER, "cuts", "cut"},
	&MatchUp{PAPER, ROCK, "covers", "covered"},
	&MatchUp{ROCK, LIZARD, "crushes", "crushed"},
	&MatchUp{LIZARD, SPOCK, "poisons", "poisoned"},
	&MatchUp{SPOCK, SCISSORS, "smashes", "smashed"},
	&MatchUp{SCISSORS, LIZARD, "decapitates", "decapitated"},
	&MatchUp{LIZARD, PAPER, "eats", "eaten"},
	&MatchUp{PAPER, SPOCK, "disproves", "disproved"},
	&MatchUp{SPOCK, ROCK, "vaporizes", "vaporized"},
	&MatchUp{ROCK, SCISSORS, "crushes", "crushed"},
}

func (m *MatchUp) WinResult() string {
	return fmt.Sprintf("%v %v %v", m.p1, m.w, m.p2)
}

func (m *MatchUp) LoseResult() string {
	tobe := "is"
	if m.p2 == SCISSORS {
		tobe = "are"
	}
	return fmt.Sprintf("%v %v %v by %v", m.p2, tobe, m.l, m.p1)
}

func (m Move) String() string {
	var names = []string{
		"Rock",
		"Spock",
		"Paper",
		"Lizard",
		"Scissors",
	}
	return names[m]
}

func (m1 Move) Versus(m2 Move) string {
	matchUp, err := findMatchUp(m1, m2)
	if err != nil {
		log.Fatal(err)
	}
	if m1 == m2 || m1.Beats(m2) {
		return matchUp.WinResult()
	}
	return matchUp.LoseResult()
}

func (m1 Move) Beats(m2 Move) bool {
	return m1 != m2 && (int(m1-m2)+LEN_Move)%LEN_Move <= 2
}

func findMatchUp(p1, p2 Move) (*MatchUp, error) {
	if p1 == p2 {
		return &MatchUp{p1, p2, "ties", ""}, nil
	}
	for _, m := range pairings {
		if m.p1 == p1 && m.p2 == p2 || m.p1 == p2 && m.p2 == p1 {
			return m, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No pairing found for %v vs %v", p1, p2))
}

func randomMove() Move {
	return Move(rand.Intn(LEN_Move))
}

func random10matches() {
	for i := 0; i < 10; i++ {
		p1, p2 := randomMove(), randomMove()
		fmt.Println(p1.Versus(p2))
	}
}

func showAllMatchUps() {
	for p1 := ROCK; int(p1) < LEN_Move; p1++ {
		for p2 := ROCK; int(p2) < LEN_Move; p2++ {
			fmt.Println(p1.Versus(p2))
		}
	}
}

func showWinningMatchUps() {
	for p1 := ROCK; int(p1) < LEN_Move; p1++ {
		for p2 := p1 + 1; int(p2) < LEN_Move; p2++ {
			if p1.Beats(p2) {
				fmt.Println(p1.Versus(p2))
			} else if p2.Beats(p1) {
				fmt.Println(p2.Versus(p1))
			}
		}
	}
}

func SheldonExplains() {
	for i, m := range pairings {
		if i == len(pairings)-1 {
			fmt.Print("...and as it always has, ")
		}
		fmt.Println(m.WinResult())
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	fmt.Println("All matchups:")
	showAllMatchUps()

	fmt.Println("\nWinning matchups:")
	showWinningMatchUps()

	fmt.Println("\n10 random matchups:")
	random10matches()

	fmt.Println("\nSheldon explains Rock-Paper-Scissors-Lizard-Spock:")
	SheldonExplains()
}
