package game

import "testing"

func TestNextTurn(t *testing.T) {
	if nextTurn(0, 0) != -1 {
		t.Fatalf("nextTurn(0, 0) != -1")
	}
	
	if nextTurn(0, 1) != 0 {
		t.Fatalf("nextTurn(0, 1) = %d", nextTurn(0, 1))
	}
	
	if nextTurn(1, 3) != 2 {
		t.Fatalf("nextTurn(1, 3) != 2")
	}
	
	if nextTurn(2, 3) != 0 {
		t.Fatalf("nextTurn(2, 3) != 0")
	}
}

func TestAddClaim(t *testing.T) {
	c := &Claim{}
	_, err := addClaim([]*Claim{}, c)
	if err != c.IsValid() {
		t.Fatalf("addClaim doesn't check for underlying claim validaity %v", err)
	}
	
	player := &Player{Hand: Hand{CardAssassin, CardContessa}}
	c, err = NewClaim(player, CardAssassin)
	if err != nil {
		t.Fatalf("NewClaim: %v", err)
	}
	
	newc, err := NewClaim(player, CardAssassin)
	if err != nil {
		t.Fatalf("NewClaim: %v", err)
	}
	
	_, err = addClaim([]*Claim{c}, newc)
	if err != ErrInvalidCounterClaim {
		t.Fatalf("addClaim: %v", err)
	}
}

func TestGameNextTurn(t *testing.T) {
	g := &Game{}
	
	g.max = 2
	g.turn = 0
	
	g.NextTurn()
	if g.turn != 1 {
		t.Fatalf("Game.NextTurn doesn't change the turn field")
	}
}

func TestGameClaim(t *testing.T) {
	g := &Game{}
	g.claims = []*Claim{}
	
	if err := g.Claim(&Claim{}); err == nil {
		t.Fatalf("Game.Claim should return nil")
	}
	
	c, err := NewClaim(&Player{Hand: Hand{0: CardContessa}}, CardContessa)
	if err != nil {
		t.Fatalf("bad test: NewClaim: %s", err.Error())
	}
	
	if err = g.Claim(c); err != nil {
		t.Fatalf("Game.Claim: %s", err.Error())
	}
}