// Source: https://github.com/alexdunne/gophercises/tree/main/deck-of-cards
//go:generate stringer -type=Suit,Rank
package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker
)

var suits = []Suit{Spade, Diamond, Club, Heart}

type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace
	maxRank = King
)

type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}

	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

type OptsFunc func([]Card) []Card

// New creates a deck of cards represented a slice. A deck is made up of a card of each Rank for each Suit
func New(opts ...OptsFunc) []Card {
	var cards []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}

	for _, opt := range opts {
		cards = opt(cards)
	}

	return cards
}

func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(c []Card) []Card {
		sort.Slice(c, less(c))
		return c
	}
}

func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

// absRank calculates the absolute rank of a given card based on its Suit and Rank
func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for i, j := range r.Perm(len(cards)) {
		ret[i] = cards[j]
	}

	return ret
}

// Jokers adds an arbitary number of Joker cards to the deck
func Jokers(n int) func([]Card) []Card {
	return func(c []Card) []Card {
		for i := 0; i < n; i++ {
			c = append(c, Card{
				Rank: Rank(i),
				Suit: Joker,
			})
		}

		return c
	}
}

// Filter calls a filter function on each card and returns a []Card of all Cards that pass the test
func Filter(f func(c Card) bool) func([]Card) []Card {
	return func(c []Card) []Card {
		var ret []Card

		for _, card := range c {
			if f(card) {
				ret = append(ret, card)
			}
		}

		return ret
	}
}

// Deck appends a duplicate set of given cards to the given cards
func Deck(n int) func([]Card) []Card {
	return func(c []Card) []Card {
		var ret []Card

		for i := 0; i < n; i++ {
			ret = append(ret, c...)
		}

		return ret
	}
}
