package main

import (
	"fmt"
	"strings"

	"github.com/alexdunne/gophercises/blackjack/deck"
)

type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}

	return strings.Join(strs, ", ")
}

func (h Hand) DealerString() string {
	return h[0].String() + ", HIDDEN"
}

func (h Hand) Score() int {
	minScore := h.minScore()
	// if the min score is greater than 11 then we can't convert the Ace to 11 as we'll bust
	if minScore > 11 {
		return minScore
	}

	for _, c := range h {
		if c.Rank == deck.Ace {
			// One of the cards is an Ace and we have room to change it for an 11 rather than a 1
			return minScore + 10
		}
	}

	return minScore
}

func (h Hand) minScore() int {
	score := 0
	for _, card := range h {
		score += min(int(card.Rank), 10)
	}

	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func main() {
	cards := deck.New(deck.Deck(3), deck.Shuffle)
	var card deck.Card
	var player, dealer Hand

	for i := 0; i < 2; i++ {
		for _, hand := range []*Hand{&player, &dealer} {
			card, cards = drawCard(cards)
			*hand = append(*hand, card)
		}
	}

	var input string
	for input != "s" {
		fmt.Println("Player: ", player)
		fmt.Println("Dealer: ", dealer.DealerString())
		fmt.Println("(s)tand or (h)it?")
		fmt.Scanf("%s\n", &input)

		switch input {
		case "h":
			card, cards = drawCard(cards)
			player = append(player, card)
		}
	}

	// if the dealer's current score is 16 or a "soft 17" (scoring 17 but 11 of the points are made up by an ace)
	for dealer.Score() <= 16 || (dealer.Score() == 17 && dealer.minScore() != 17) {
		card, cards = drawCard(cards)
		dealer = append(dealer, card)
	}

	playerScore := player.Score()
	dealerScore := dealer.Score()

	fmt.Println("---")
	fmt.Println("Player: ", player, "\nScore: ", playerScore)
	fmt.Println("Dealer: ", dealer, "\nScore: ", dealerScore)

	switch {
	case playerScore > 21:
		fmt.Println("Player busted")
	case dealerScore > 21:
		fmt.Println("Dealer busted")
	case playerScore > dealerScore:
		fmt.Println("Player won")
	case dealerScore > playerScore:
		fmt.Println("Dealer won")
	case dealerScore == playerScore:
		fmt.Println("Draw")
	}
}

// drawCard returns the first card from the given cards and a new slice with that first element removed
func drawCard(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}
