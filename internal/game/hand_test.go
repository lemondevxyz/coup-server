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