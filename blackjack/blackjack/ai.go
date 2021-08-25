package blackjack

import (
	"fmt"

	"github.com/alexdunne/gophercises/blackjack/deck"
)

type AI interface {
	Bet() int
	Play(hand []deck.Card, dealer deck.Card) Move
	GameResults(hand [][]deck.Card, dealer []deck.Card)
}

func HumanAI() humanAI {
	return humanAI{}
}

type humanAI struct{}

func (ai humanAI) Bet() int {
	return 1
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("Player: ", hand)
		fmt.Println("Dealer: ", dealer)
		fmt.Println("(s)tand or (h)it?")

		var input string
		fmt.Scanf("%s\n", &input)

		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		default:
			fmt.Println("Invalid option")
		}
	}
}

func (ai humanAI) GameResults(hand [][]deck.Card, dealer []deck.Card) {
	fmt.Println("Results")
	fmt.Println("Player: ", hand)
	fmt.Println("Dealer: ", dealer)
	fmt.Println("")
}
