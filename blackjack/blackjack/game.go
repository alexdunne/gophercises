package blackjack

import (
	"fmt"

	"github.com/alexdunne/gophercises/blackjack/deck"
)

type state int8

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

// New creates a new instance of a Game with default values
func New() Game {
	return Game{
		state: statePlayerTurn,
	}
}

type Game struct {
	deck    []deck.Card
	state   state
	player  []deck.Card
	dealer  []deck.Card
	balance int
}

func (g *Game) currentHand() *[]deck.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("it is neither player's turn")
	}
}

// deal supplies the Player and Dealer with two cards each
func deal(g *Game) {
	g.player = make([]deck.Card, 0, 5)
	g.dealer = make([]deck.Card, 0, 5)

	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.deck = drawCard(g.deck)
		g.player = append(g.player, card)

		card, g.deck = drawCard(g.deck)
		g.dealer = append(g.dealer, card)
	}

	g.state = statePlayerTurn
}

func dealerPlay(hand []deck.Card) Move {
	score := Score(hand...)
	if score <= 16 || (score == 17 && minScore(hand...) != 17) {
		return MoveHit
	} else {
		return MoveStand
	}
}

func (g *Game) Play(ai AI) int {
	g.deck = deck.New(deck.Deck(3), deck.Shuffle)

	for i := 0; i < 10; i++ {
		deal(g)

		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(g.player))
			copy(hand, g.player)
			move := ai.Play(hand, g.dealer[0])
			move(g)
		}

		for g.state == stateDealerTurn {
			move := dealerPlay(g.dealer)
			move(g)
		}

		endHand(g, ai)
	}

	return g.balance
}

type Move func(*Game)

func MoveHit(g *Game) {
	var card deck.Card
	card, g.deck = drawCard(g.deck)

	hand := g.currentHand()
	*hand = append(*hand, card)

	if Score(*hand...) > 21 {
		MoveStand(g)
	}
}

func MoveStand(g *Game) {
	switch g.state {
	case statePlayerTurn:
		g.state = stateDealerTurn
	case stateDealerTurn:
		g.state = stateHandOver
	default:
		panic("it is neither player's turn")
	}
}

func drawCard(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

func endHand(g *Game, ai AI) {
	playerScore := Score(g.player...)
	dealerScore := Score(g.dealer...)

	switch {
	case playerScore > 21:
		fmt.Println("Player busted")
		g.balance--

	case dealerScore > 21:
		fmt.Println("Dealer busted")
		g.balance++

	case playerScore > dealerScore:
		fmt.Println("Player won")
		g.balance++

	case dealerScore > playerScore:
		fmt.Println("Dealer won")
		g.balance--

	case dealerScore == playerScore:
		fmt.Println("Draw")
	}
	fmt.Println()

	ai.GameResults([][]deck.Card{g.player}, g.dealer)

	g.player = nil
	g.dealer = nil
}

// Score reports the highest valid blackjack score of the hand
func Score(hand ...deck.Card) int {
	minScore := minScore(hand...)
	// if the min score is greater than 11 then we can't convert the Ace to 11 as we'll bust
	if minScore > 11 {
		return minScore
	}

	for _, c := range hand {
		if c.Rank == deck.Ace {
			// One of the cards is an Ace and we have room to change it for an 11 rather than a 1
			return minScore + 10
		}
	}

	return minScore
}

// Soft reports if an Ace in the hand is being used as an 11
func Soft(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)

	return minScore != score
}

func minScore(hand ...deck.Card) int {
	score := 0
	for _, card := range hand {
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
