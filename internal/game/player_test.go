package game

import (
	"testing"
)

func TestPlayerIsDead(t *testing.T) {
	p := &Player{}
	if !p.IsDead() {
		t.Fatalf("player should be dead but isn't")
	}
	
	p.Hand[0] = DukeCard
	if !p.IsDead() {
		t.Fatalf("player can be revived if his hand changed")
	}
	
	p = &Player{}
	p.Hand[0] = AssassinCard
	
	if p.IsDead() {
		t.Fatalf("non empty hand makes a player dead...")
	}
}