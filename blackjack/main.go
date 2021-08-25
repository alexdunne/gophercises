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

const blackjack = 21

type State int

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("it is neither player's turn")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}

	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Player)
	copy(ret.Dealer, gs.Dealer)

	return ret
}

// Shuffle create a new shuffled deck of cards
func ShuffleDeck(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

// Deal supplier the Player and Dealer with two cards each
func Deal(gs GameState) GameState {
	ret := clone(gs)
	// allocate more than the initial cards so we don't have to resize later
	ret.Player = Hand{}
	ret.Dealer = Hand{}

	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Deck = ret.Deck[0], ret.Deck[1:]
		ret.Player = append(ret.Player, card)

		card, ret.Deck = ret.Deck[0], ret.Deck[1:]
		ret.Dealer = append(ret.Dealer, card)
	}

	ret.State = StatePlayerTurn
	return ret
}

// Hit supplies the current player with a card
func Hit(gs GameState) GameState {
	ret := clone(gs)

	var card deck.Card
	card, ret.Deck = ret.Deck[0], ret.Deck[1:]

	hand := ret.CurrentPlayer()
	*hand = append(*hand, card)

	if hand.Score() > blackjack {
		return Stand(ret)
	}

	return ret
}

// Stand progresses the GameState to the next state
func Stand(gs GameState) GameState {
	ret := clone(gs)

	switch ret.State {
	case StatePlayerTurn:
		ret.State = StateDealerTurn
	case StateDealerTurn:
		ret.State = StateHandOver
	default:
		panic("it is neither player's turn")
	}

	return ret
}

// EndHand announces the winner and tidies up the game state
func EndHand(gs GameState) GameState {
	ret := clone(gs)

	playerScore := ret.Player.Score()
	dealerScore := ret.Dealer.Score()

	fmt.Println("---")
	fmt.Println("Player: ", ret.Player, "\nScore: ", playerScore)
	fmt.Println("Dealer: ", ret.Dealer, "\nScore: ", dealerScore)

	switch {
	case playerScore > blackjack:
		fmt.Println("Player busted")
	case dealerScore > blackjack:
		fmt.Println("Dealer busted")
	case playerScore > dealerScore:
		fmt.Println("Player won")
	case dealerScore > playerScore:
		fmt.Println("Dealer won")
	case dealerScore == playerScore:
		fmt.Println("Draw")
	}
	fmt.Println()

	ret.Player = nil
	ret.Dealer = nil

	return ret
}

func main() {
	var gs GameState
	gs = ShuffleDeck(gs)
	gs = Deal(gs)

	var input string
	for gs.State == StatePlayerTurn {
		fmt.Println("Player: ", gs.Player)
		fmt.Println("Dealer: ", gs.Dealer.DealerString())
		fmt.Println("(s)tand or (h)it?")
		fmt.Scanf("%s\n", &input)

		switch input {
		case "h":
			gs = Hit(gs)
		case "s":
			gs = Stand(gs)
		}
	}

	for gs.State == StateDealerTurn {
		if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.minScore() != 17) {
			gs = Hit(gs)
		} else {
			gs = Stand(gs)
		}
	}

	gs = EndHand(gs)
}
