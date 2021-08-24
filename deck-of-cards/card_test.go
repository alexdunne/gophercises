package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Queen, Suit: Spade})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ace of Hearts
	// Queen of Spades
	// Joker
}

func TestCreatesNewSuit(t *testing.T) {
	cards := New()

	if len(cards) != 52 {
		t.Errorf("expected deck to have 52 cards, got %d", len(cards))
	}
}

func TestNewSuitIsSortedInDefaultOrder(t *testing.T) {
	suits := []Suit{Spade, Diamond, Club, Heart}
	cards := New()

	// 4 suits, 13 cards each
	for cardIndex, card := range cards {
		expectedSuit := int(cardIndex / 13)
		expectedRank := cardIndex%13 + 1

		if card.Suit != suits[expectedSuit] {
			t.Errorf("expected card (%s) at index %d to be of suit type %s, got %s", card, cardIndex, suits[expectedSuit], card.Suit)
		}

		if int(card.Rank) != expectedRank {
			t.Errorf("expected card (%s) at index %d to have rank %d, got %d", card, cardIndex, expectedRank, card.Rank)
		}
	}
}

func TestSortingWithACustomSortFunction(t *testing.T) {
	suits := []Suit{Spade, Diamond, Club, Heart}
	cards := New(Sort(Less))

	// 4 suits, 13 cards each
	for cardIndex, card := range cards {
		expectedSuit := int(cardIndex / 13)
		expectedRank := cardIndex%13 + 1

		if card.Suit != suits[expectedSuit] {
			t.Errorf("expected card (%s) at index %d to be of suit type %s, got %s", card, cardIndex, suits[expectedSuit], card.Suit)
		}

		if int(card.Rank) != expectedRank {
			t.Errorf("expected card (%s) at index %d to have rank %d, got %d", card, cardIndex, expectedRank, card.Rank)
		}
	}
}

func TestAddingJokersToADeck(t *testing.T) {
	type testcase struct {
		jokers int
	}

	tests := []testcase{
		{jokers: 5},
		{jokers: 5},
		{jokers: 1},
		{jokers: 0},
	}

	for _, tc := range tests {
		cards := New(Jokers(tc.jokers))
		count := 0

		for _, card := range cards {
			if card.Suit == Joker {
				count++
			}
		}

		if count != tc.jokers {
			t.Errorf("expected %d jokers, got %d", tc.jokers, count)
		}
	}
}

func TestFilteringADeckOfCards(t *testing.T) {
	filter := func(c Card) bool {
		return c.Rank != Ace
	}

	cards := New(Filter(filter))

	for _, c := range cards {
		if c.Rank == Ace {
			t.Errorf("expected all cards with rank %s to be filtered, got %s", c.Rank.String(), c.String())
		}
	}
}

func TestDesk(t *testing.T) {
	cards := New(Deck(3))

	if len(cards) != 52*3 {
		t.Errorf("expected %d cards, got %d", 52*3, len(cards))
	}
}
