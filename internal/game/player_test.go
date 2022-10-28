package game

import (
	"testing"
	"github.com/matryer/is"
)

func TestPlayerIsDead(t *testing.T) {
	is := is.New(t)
	p := &Player{}
	is.True(p.IsDead())
	
	p.Hand[0] = CardDuke
	is.True(p.IsDead())
	
	p = &Player{}
	p.Hand[0] = CardAssassin
	
	is.True(!p.IsDead())
}