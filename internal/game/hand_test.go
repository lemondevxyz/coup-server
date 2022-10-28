package game

import (
	"testing"
	"github.com/matryer/is"
)

func TestHandIsEmpty(t *testing.T) {
	is := is.New(t)
	empty := Hand{CardEmpty, CardEmpty}
	
	is.True(empty.IsEmpty())
	
	nonEmpty := [3]Hand{
		Hand{CardEmpty, CardDuke},
		Hand{CardDuke, CardEmpty},
		Hand{CardDuke, CardDuke},
	}
	
	for _, v := range nonEmpty {
		is.True(!v.IsEmpty())
	}
}

func TestHandEqual(t *testing.T) {
	is := is.New(t)
	hand := Hand{CardEmpty, CardEmpty}
	notEqual := [3]Hand{
		Hand{CardEmpty, CardDuke},
		Hand{CardDuke, CardEmpty},
		Hand{CardAssassin, CardContessa},
	}
	equal := Hand{CardEmpty, CardEmpty}
	
	for _, v := range notEqual {
		is.True(!v.IsEqual(hand))
	}
	
	is.Equal(equal, hand)
}