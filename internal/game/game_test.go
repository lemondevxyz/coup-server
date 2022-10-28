package game

import (
	"testing"
	"github.com/matryer/is"
)

func TestNextTurn(t *testing.T) {
	is := is.New(t)
	is.Equal(nextTurn(0, 0), -1)
	is.Equal(nextTurn(0, 1), 0)
	is.Equal(nextTurn(1, 3), 2)
	is.Equal(nextTurn(2, 3), 0)
}

func TestAddClaim(t *testing.T) {
	c := &Claim{}
	_, err := addClaim([]*Claim{}, c)
	
	is := is.New(t)
	is.Equal(err, c.IsValid())
	
	player := &Player{Hand: Hand{CardAssassin, CardContessa}}
	c, err = NewClaim(player, CardAssassin)
	is.NoErr(err)
	
	newc, err := NewClaim(player, CardAssassin)
	is.NoErr(err)
	
	_, err = addClaim([]*Claim{c}, newc)
	is.Equal(err, ErrInvalidCounterClaim)
}

func TestGameNextTurn(t *testing.T) {
	g := &Game{}
	
	g.max = 2
	g.turn = 0
	
	g.NextTurn()

	is := is.New(t)
	is.Equal(g.turn, 1)
}

func TestGameClaim(t *testing.T) {
	g := &Game{}
	g.claims = []*Claim{}
	
	is := is.New(t)
	is.True(g.Claim(&Claim{}) != nil)
	
	c, err := NewClaim(&Player{Hand: Hand{0: CardContessa}}, CardContessa)
	is.NoErr(err)
	is.NoErr(g.Claim(c))
}