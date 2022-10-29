package game

import (
	"github.com/matryer/is"
	"testing"
)

func TestHandString(t *testing.T) {
	h := Hand{0: 5, 1: 4}
	is := is.New(t)
	is.Equal(h.String(), "5:4")
}

func TestHandUnmarshal(t *testing.T) {
	h := &Hand{}

	is := is.New(t)
	is.True(h.Unmarshal("") != nil)
	is.True(h.Unmarshal("a:4") != nil)
	is.True(h.Unmarshal("4:a") != nil)
	is.NoErr(h.Unmarshal("4:4"))

	is.Equal(*h, Hand{4, 4})
}

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
