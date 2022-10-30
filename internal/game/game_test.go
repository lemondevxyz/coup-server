package game

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestFindPlayerByPntr(t *testing.T) {
	is := is.New(t)

	v := &Player{}
	pl := []*Player{v, nil}

	is.Equal(findPlayerByPntr(pl, v), 0)
	is.Equal(findPlayerByPntr(pl, nil), -1)
}

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

/*
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
*/

func TestGameAction(t *testing.T) {
	g, err := NewGame([5]*Player{{Hand: Hand{CardAmbassador, CardAssassin}}, {Hand: Hand{CardDuke, CardContessa}}})

	is := is.New(t)
	is.NoErr(err)

	a1, a2 := &Action{}, &Action{}
	{
		g := &Game{}
		is.Equal(g.Action(a1), ErrInvalidActionAuthor)
	}
	a1.AuthorID, a2.AuthorID = 0, 0
	a1.Kind, a2.Kind = ActionIncome, ActionIncome

	is.True(g.Action(nil) == ErrInvalidAction)

	a1.AuthorID = 4
	is.Equal(g.Action(a1), ErrInvalidActionAuthor)

	a1.AuthorID = 0

	v := uint8(4)
	a1.AgainstID = &v
	is.Equal(g.Action(a1), ErrInvalidActionAgainst)
	*a1.AgainstID = 1

	is.NoErr(g.Action(a1))
	is.Equal(a1.against, g.players[1])

	g.claim = nil
	is.Equal(g.Action(a1), ErrInvalidGame)
	g.claim = NewNotifier()

	g.action[0], g.history = nil, []Action{}
	a1.Kind = 54
	is.Equal(g.Action(a1), ErrInvalidActionKind)

	boolval := true
	g.claim.Set(&Claim{
		challenge: &boolval,
	})

	a1.Kind = ActionIncome
	is.Equal(g.Action(a1), ErrInvalidActionFrozen)

	g.claim.Set(&Claim{
		succeed: &boolval,
	})

	is.NoErr(g.Action(a1))
	c, _ := g.ClaimGet()
	is.Equal(c, nil)
	g.history, g.action[0] = []Action{}, nil

	g.claim.Set(&Claim{})
	is.Equal(g.Action(a1), ErrInvalidActionFrozen)

	g.claim.Set(nil)

	is.NoErr(g.Action(a1))
	is.Equal(a1, g.action[0])

	is.Equal(g.history[0], *g.action[0])

	is.Equal(g.Action(a2), ErrInvalidCounterClaim)

	g.action[0].Kind = ActionFinancialAid
	a2.Kind = ActionCharacter
	a2.Character = CardDuke

	is.NoErr(g.Action(a2))
	is.Equal(a2, g.action[1])
	is.Equal(g.history[1], *g.action[1])
	is.Equal(g.Action(a2), ErrInvalidAction)

	g.action[0].Kind = ActionCharacter
	g.action[0].Character = CardAssassin
	a2.Character = CardContessa

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
	is := is.New(t)

	_, _, err := g.TurnSubscribe()
	is.Equal(err, ErrInvalidGame)

	g.max = 2
	g.turn = NewNotifier()
	g.turn.Set(0)

	_, ch, err := g.TurnSubscribe()
	is.NoErr(err)
	wait := make(chan struct{})

	go func() {
		val := <-ch
		wait <- val
	}()
	g.NextTurn()

	is.Equal(g.turn.Get().(int), 1)

	select {
	case <-wait:
	case <-time.After(time.Millisecond):
		t.Fatalf("NextTurn doesn't notify subscribers")
	}
}

func TestGameClaim(t *testing.T) {
	g := &Game{}

	is := is.New(t)
	_, _, err := g.ClaimSubscribe()
	is.Equal(err, ErrInvalidGame)

	g.claim = NewNotifier()

	_, ch, err := g.ClaimSubscribe()
	is.NoErr(err)

	wait := make(chan struct{})

	go func() {
		val := <-ch
		wait <- val
	}()

	is.True(g.Claim(&Claim{}) != nil)

	pl := &Player{Hand: Hand{0: CardContessa}}

	c, err := NewClaim(pl, CardContessa)
	is.NoErr(err)
	is.True(g.Claim(c) == ErrInvalidPlayer)

	g.players = [5]*Player{pl}

	c.Pass()
	is.True(g.Claim(c) == ErrInvalidClaimFinished)

	c.succeed = nil
	is.NoErr(g.Claim(c))
	is.Equal(c.Action(0), g.history[0])

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

	newc, err := g.ClaimGet()
	is.NoErr(err)
	is.Equal(newc, c)
}

func TestGameClaimUnsubscribe(t *testing.T) {
	g := &Game{}

	is := is.New(t)
	is.Equal(g.ClaimUnsubscribe(time.Time{}), ErrInvalidGame)

	g.claim = NewNotifier()
	g.claim.chnls[time.Time{}] = nil
	is.NoErr(g.ClaimUnsubscribe(time.Time{}))
}

func TestGameTurnGet(t *testing.T) {
	is := is.New(t)

	g := &Game{}
	_, err := g.TurnGet()
	is.Equal(err, ErrInvalidGame)

	g.turn = NewNotifier()
	g.turn.Set(3)

	newt, err := g.TurnGet()
	is.NoErr(err)
	is.Equal(newt, 3)
}

func TestGameTurnUnsubscribe(t *testing.T) {
	g := &Game{}

	is := is.New(t)
	is.Equal(g.TurnUnsubscribe(time.Time{}), ErrInvalidGame)

	g.turn = NewNotifier()
	g.turn.chnls[time.Time{}] = nil
	is.NoErr(g.TurnUnsubscribe(time.Time{}))
}

func TestValidateClaimAndItsPlayer(t *testing.T) {
	_, err := validateClaimAndItsPlayer(nil, &Claim{})

	is := is.New(t)
	is.Equal(err, ErrInvalidArr)
	_, err = validateClaimAndItsPlayer([]*Player{}, nil)
	is.Equal(err, ErrInvalidArr)
	_, err = validateClaimAndItsPlayer([]*Player{{}}, nil)
	is.Equal(err, ErrInvalidClaim)

	arr := []*Player{}

	player := &Player{}
	other := &Player{}

	arr = append(arr, player)
	_, err = validateClaimAndItsPlayer(arr, &Claim{
		author: other,
	})
	is.Equal(err, ErrInvalidPlayer)

	index, err := validateClaimAndItsPlayer(arr, &Claim{
		author: player,
	})
	is.Equal(index, 0)
	is.Equal(err, nil)
}

func TestGameClaimPass(t *testing.T) {
	is := is.New(t)

	g := &Game{}

	is.Equal(g.ClaimPass(), ErrInvalidGame)
	g.claim = NewNotifier()

	g.claim.Set(&Claim{
		succeed: new(bool),
	})

	is.Equal(g.ClaimPass(), ErrInvalidClaimFinished)

	g.players = [5]*Player{
		0: {Hand: Hand{0: CardAmbassador}},
	}

	g.claim.Set(&Claim{})

	newc, err := g.ClaimGet()
	is.NoErr(err)
	index, err := g.validateClaimAndItsPlayer(newc)

	is.Equal(g.ClaimPass(), err)
	g.claim.Set(&Claim{
		author: g.players[0],
	})

	is.NoErr(g.ClaimPass())
	t.Log(index, g.history[0])
	is.Equal(g.history[0].AuthorID, uint8(0))
}

func TestGameClaimChallenge(t *testing.T) {
	is := is.New(t)

	{
		g := &Game{}

		_, err := g.ClaimGet()
		is.Equal(g.ClaimChallenge(nil), err)
	}

	pl1, pl2 := &Player{Hand: Hand{0: CardAmbassador}}, &Player{Hand: Hand{1: CardCaptain}}

	g, err := NewGame([5]*Player{pl1, pl2})

	is.NoErr(err)
	is.Equal(g.ClaimChallenge(nil), ErrInvalidClaim)

	is.NoErr(g.Claim(&Claim{
		author:    pl1,
		character: CardAmbassador,
	}))

	c := &Claim{
		author:    &Player{},
		character: CardAmbassador,
	}
	g.claim.Set(c)

	_, err = g.validateClaimAndItsPlayer(c)
	is.Equal(g.ClaimChallenge(&Player{
		Hand: Hand{0: CardAmbassador},
	}), err)

	c, err = g.ClaimGet()
	is.NoErr(err)

	c.Pass()
	is.Equal(g.ClaimChallenge(&Player{}), ErrInvalidClaimFinished)

	is.NoErr(g.Claim(&Claim{
		author:    pl1,
		character: CardAmbassador,
	}))

	is.Equal(g.ClaimChallenge(&Player{}), ErrInvalidPlayer)
	is.Equal(g.ClaimChallenge(pl2), nil)
}
