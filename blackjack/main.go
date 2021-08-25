package main

import (
	"fmt"

	"github.com/alexdunne/gophercises/blackjack/blackjack"
)

func main() {
	game := blackjack.New()
	winnings := game.Play(blackjack.HumanAI())

	fmt.Println(winnings)
}
