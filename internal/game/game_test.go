package game

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestNextTurn(t *testing.T) {
	is := is.New(t)
	is.Equal(nextTurn(0, 0), -1)
	is.Equal(nextTurn(0, 1), 0)
	is.Equal(nextTurn(1, 3), 2)
	is.Equal(nextTurn(2, 3), 0)
}

func TestNewGame(t *testing.T) {
	p := &Player{}

	_, err := NewGame([5]*Player{p})

	is := is.New(t)
	is.Equal(err, ErrInvalidPlayer)

	p = &Player{Hand: Hand{0: CardAmbassador}}
	_, err = NewGame([5]*Player{p})
	is.Equal(err, ErrInvalidPlayerAmount)

	g, err := NewGame([5]*Player{p, p})
	is.NoErr(err)
	is.Equal(g.max, 2)
}

func TestAddClaim(t *testing.T) {
	c := &Claim{}
	_, err := addClaim([]*Claim{}, c)

	is := is.New(t)
	is.Equal(err, c.IsValid())

	val, err := addClaim(nil, nil)
	is.Equal(val, nil)
	is.Equal(err, ErrInvalidParameters)
	player := &Player{Hand: Hand{CardAssassin, CardContessa}}
	c, err = NewClaim(player, CardAssassin)
	is.NoErr(err)

	newc, err := NewClaim(player, CardAssassin)
	is.NoErr(err)

	_, err = addClaim([]*Claim{c}, newc)
	is.Equal(err, ErrInvalidCounterClaim)

	arr, err := addClaim([]*Claim{}, c)
	is.NoErr(err)
	is.Equal([]*Claim{c}, arr)
}

func TestGameAction(t *testing.T) {
	g := &Game{}

	a1, a2 := &Action{}, &Action{}

	is := is.New(t)

	is.NoErr(g.Action(a1))
	is.Equal(a1, g.action[0])

	is.NoErr(g.Action(a2))
	is.Equal(a2, g.action[1])

	is.Equal(g.Action(a2), ErrInvalidAction)
}

func TestGameDoAction(t *testing.T) {
	g := &Game{}

	is := is.New(t)

	err := g.DoAction()
	is.Equal(err, ErrInvalidAction)

	act := &Action{Character: CardDuke}
	act.Kind = ActionCharacter
	act.author = &Player{Hand: Hand{CardAmbassador, CardDuke}}

	g.action = [2]*Action{act, nil}
	is.NoErr(g.DoAction())

	is.Equal(act.author.Coins, uint8(3))
	is.Equal(g.action[0], nil)

	g.action = [2]*Action{nil, act}
	is.NoErr(g.DoAction())

	is.Equal(act.author.Coins, uint8(6))
	is.Equal(g.action[0], nil)
	is.Equal(g.action[1], nil)
}

func TestGameNextTurn(t *testing.T) {
	g := &Game{}

	g.max = 2
	g.turn = NewNotifier()
	g.turn.Set(0)

	_, ch := g.TurnSubscribe()
	wait := make(chan struct{})

	go func() {
		val := <-ch
		wait <- val
	}()
	g.NextTurn()

	is := is.New(t)
	is.Equal(g.turn.Get().(int), 1)

	select {
	case <-wait:
	case <-time.After(time.Millisecond):
		t.Fatalf("NextTurn doesn't notify subscribers")
	}
}

func TestGameClaim(t *testing.T) {
	g := &Game{}
	g.claim = NewNotifier()

	_, ch := g.ClaimSubscribe()
	wait := make(chan struct{})

	go func() {
		val := <-ch
		wait <- val
	}()

	is := is.New(t)
	is.True(g.Claim(&Claim{}) != nil)

	c, err := NewClaim(&Player{Hand: Hand{0: CardContessa}}, CardContessa)
	is.NoErr(err)
	is.NoErr(g.Claim(c))

	select {
	case <-wait:
	case <-time.After(time.Millisecond):
		t.Fatalf("g.Claim doesn't notify subscribers")
	}
}

func TestGameClaimGet(t *testing.T) {
	g := &Game{}
	g.claim = NewNotifier()
	c := &Claim{}
	g.claim.Set(c)

	is := is.New(t)
	is.Equal(g.ClaimGet(), c)
}

func TestGameClaimUnsubscribe(t *testing.T) {
	g := &Game{}
	g.claim = NewNotifier()
	g.claim.chnls[time.Time{}] = nil
	g.ClaimUnsubscribe(time.Time{})
	g.claim.Announce()
}

func TestGameTurnGet(t *testing.T) {
	g := &Game{}
	g.turn = NewNotifier()
	g.turn.Set(3)

	is := is.New(t)
	is.Equal(g.TurnGet(), 3)
}

func TestGameTurnUnsubscribe(t *testing.T) {
	g := &Game{}
	g.turn = NewNotifier()
	g.turn.chnls[time.Time{}] = nil
	g.TurnUnsubscribe(time.Time{})
	g.turn.Announce()
}
