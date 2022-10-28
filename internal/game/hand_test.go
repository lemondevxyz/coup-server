package game

import (
	"testing"
)

func TestHandIsEmpty(t *testing.T) {
	empty := Hand{CardEmpty, CardEmpty}
	
	if !empty.IsEmpty() {
		t.Fatalf("the empty hand is non-empty")
	}
	
	nonEmpty := [3]Hand{
		Hand{CardEmpty, CardDuke},
		Hand{CardDuke, CardEmpty},
		Hand{CardDuke, CardDuke},
	}
	
	for _, v := range nonEmpty {
		if v.IsEmpty() {
			t.Fatalf("non-empty hand is empty: %v", v)
		}
	}
}

func TestHandEqual(t *testing.T) {
	hand := Hand{CardEmpty, CardEmpty}
	notEqual := [3]Hand{
		Hand{CardEmpty, CardDuke},
		Hand{CardDuke, CardEmpty},
		Hand{CardAssassin, CardContessa},
	}
	equal := Hand{CardEmpty, CardEmpty}
	
	for _, v := range notEqual {
		if v.Equal(hand) {
			t.Fatalf("hand shouldn't equal other hand: %v, %v", v, hand)
		}
	}
	
	if !equal.Equal(hand) {
		t.Fatalf("hand should equal other hand: %v, %v", equal, hand)
	}
}