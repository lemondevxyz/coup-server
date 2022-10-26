package game

import (
	"testing"
)

func TestHandIsEmpty(t *testing.T) {
	empty := Hand{EmptyCard, EmptyCard}
	
	if !empty.IsEmpty() {
		t.Fatalf("the empty hand is non-empty")
	}
	
	nonEmpty := [3]Hand{
		Hand{EmptyCard, DukeCard},
		Hand{DukeCard, EmptyCard},
		Hand{DukeCard, DukeCard},
	}
	
	for _, v := range nonEmpty {
		if v.IsEmpty() {
			t.Fatalf("non-empty hand is empty: %v", v)
		}
	}
}

func TestHandEqual(t *testing.T) {
	hand := Hand{EmptyCard, EmptyCard}
	notEqual := [3]Hand{
		Hand{EmptyCard, DukeCard},
		Hand{DukeCard, EmptyCard},
		Hand{AssassinCard, ContessaCard},
	}
	equal := Hand{EmptyCard, EmptyCard}
	
	for _, v := range notEqual {
		if v.Equal(hand) {
			t.Fatalf("hand shouldn't equal other hand: %v, %v", v, hand)
		}
	}
	
	if !equal.Equal(hand) {
		t.Fatalf("hand should equal other hand: %v, %v", equal, hand)
	}
}