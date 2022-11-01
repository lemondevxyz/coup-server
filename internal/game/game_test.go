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
	is := is.New(t)

	p := &Player{Hand: Hand{0: CardAmbassador}}
	_, err := NewGame([5]*Player{p})
	is.Equal(err, ErrInvalidPlayerAmount)

	g, err := NewGame([5]*Player{p, p})
	is.NoErr(err)
	is.Equal(g.max, 2)
}

func TestGameAction(t *testing.T) {
	g, err := NewGame([5]*Player{{Hand: Hand{CardAmbassador, CardAssassin}}, {Hand: Hand{CardDuke, CardContessa}}})

	is := is.New(t)
	is.NoErr(err)

	a1, a2 := Action{}, Action{}
	a1.AuthorID = 255

	is.Equal(g.Action(a1), a1.setPlayer(g.players[:]))

	a1.AuthorID = 0
	is.NoErr(a1.setPlayer(g.players[:]))
	is.Equal(g.Action(a1), a1.IsValid())

	a1.Kind = ActionCharacter
	is.Equal(g.Action(a1), ErrInvalidAction)
	g.claim = &claim{}

	is.Equal(g.Action(a1), a1.validClaim(g.claim))

	a1.Kind = ActionCharacter
	g.claim = &claim{character: CardContessa}
	g.claim.succeed, g.claim.challenge = new(bool), new(bool)
	a1.Character = CardContessa

	is.Equal(g.Action(a1), ErrInvalidActionKind)

	g.claim = nil

	a1.Kind = ActionFinancialAid
	is.NoErr(g.Action(a1))

	is.Equal(g.Action(Action{AuthorID: 1, Kind: ActionIncome}), ErrInvalidCounterClaim)

	g.claim = &claim{}
	g.claim.succeed = new(bool)
	*g.claim.succeed = true
	g.claim.character = CardDuke

	a2.Kind = ActionCharacter
	a2.Character = g.claim.character

	is.NoErr(g.Action(a2))
	is.Equal(g.Action(a2), ErrInvalidAction)

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
	_, err := validateClaimAndItsPlayer(nil, &claim{})

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
	_, err = validateClaimAndItsPlayer(arr, &claim{
		author: other,
	})
	is.Equal(err, ErrInvalidPlayer)

	index, err := validateClaimAndItsPlayer(arr, &claim{
		author: player,
	})
	is.Equal(index, 0)
	is.Equal(err, nil)
}

func TestGameClaim(t *testing.T) {
	g, err := NewGame([5]*Player{{Hand: Hand{CardAmbassador, CardAssassin}}, {Hand: Hand{CardDuke, CardContessa}}})
	is := is.New(t)
	is.NoErr(err)

	g.claim = &claim{}
	is.Equal(g.Claim(nil, 0), ErrInvalidClaimOngoing)

	g.claim = nil
	is.Equal(g.Claim(nil, CardContessa), (&claim{}).IsValid())
	is.Equal(g.Claim(&Player{Hand: Hand{0: CardContessa}}, CardContessa), ErrInvalidPlayer)

	is.Equal(g.Claim(g.players[0], CardAmbassador), nil)
	is.True(g.claim != nil)
	is.Equal(g.history[len(g.history)-1], g.claim.Action(0))
}

func TestGameClaimPass(t *testing.T) {
	g, err := NewGame([5]*Player{{Hand: Hand{CardAmbassador, CardAssassin}}, {Hand: Hand{CardDuke, CardContessa}}})

	is := is.New(t)
	is.NoErr(err)

	is.Equal(g.ClaimPass(), ErrInvalidClaim)
	g.claim = &claim{succeed: new(bool)}

	is.Equal(g.ClaimPass(), ErrInvalidClaimFinished)

	g.claim.succeed = nil
	g.claim.author = g.players[0]

	is.NoErr(g.ClaimPass())
	is.Equal(g.history[len(g.history)-1], g.claim.Action(0))
	is.True(g.claim.succeed != nil)
}

func TestGameClaimChallenge(t *testing.T) {
	g, err := NewGame([5]*Player{{Hand: Hand{CardAmbassador, CardAssassin}}, {Hand: Hand{CardDuke, CardContessa}}})

	is := is.New(t)
	is.NoErr(err)

	is.Equal(g.ClaimChallenge(g.players[1]), ErrInvalidClaim)
	g.claim = &claim{succeed: new(bool)}

	is.Equal(g.ClaimChallenge(g.players[1]), ErrInvalidClaimFinished)

	g.claim.succeed = nil
	g.claim.author = g.players[1]

	is.Equal(g.ClaimChallenge(g.players[1]), ErrInvalidActionSamePlayer)

	g.claim.author = g.players[0]
	is.Equal(g.ClaimChallenge(&Player{}), ErrInvalidPlayer)

	is.NoErr(g.ClaimChallenge(g.players[1]))

	is.Equal(g.history[len(g.history)-1].AuthorID, uint8(1))
	is.Equal(*g.history[len(g.history)-1].AgainstID, uint8(0))
	is.True(g.claim.challenge != nil)
}

func TestGameClaimProve(t *testing.T) {
	g, err := NewGame([5]*Player{{Hand: Hand{CardAmbassador, CardAssassin}}, {Hand: Hand{CardDuke, CardContessa}}})

	is := is.New(t)
	is.NoErr(err)

	_, err = g.ClaimProve(CardContessa)
	is.Equal(err, ErrInvalidClaim)
	g.claim = &claim{}

	_, err = g.ClaimProve(CardContessa)
	is.Equal(err, ErrInvalidClaimHasNotFinished)

	g.claim.succeed = new(bool)
	_, err = g.ClaimProve(CardContessa)
	is.Equal(err, ErrInvalidClaimNotChallenged)

	g.claim = nil

	is.NoErr(g.Claim(g.players[0], CardContessa))
	is.NoErr(g.ClaimChallenge(g.players[1]))

	result, err := g.ClaimProve(CardContessa)

	is.NoErr(err)
	is.True(result)

	is.Equal(g.history[len(g.history)-1].AuthorID, uint8(0))
	is.Equal(*g.history[len(g.history)-1].AgainstID, uint8(1))

	g.history = g.history[:len(g.history)-1]
	g.claim.succeed = nil

	result, err = g.ClaimProve(CardAssassin)
	is.NoErr(err)
	is.True(!result)
}

func TestGameDoAction(t *testing.T) {

	is := is.New(t)

	g := &Game{}

	pl1, pl2 := &Player{}, &Player{}

	g.action[1] = &Action{author: pl2, Kind: ActionIncome}
	g.action[0] = &Action{author: pl1, Kind: ActionIncome}

	is.NoErr(g.DoAction())

	is.Equal(pl1.Coins, uint8(0))
	is.Equal(pl2.Coins, uint8(1))
	is.Equal(g.action[0], nil)
	is.Equal(g.action[1], nil)

	g.action[0] = &Action{author: pl1, Kind: ActionIncome}
	is.NoErr(g.DoAction())
	is.Equal(pl1.Coins, uint8(1))

	is.Equal(g.DoAction(), ErrInvalidAction)

	g.action[0] = &Action{author: pl1, Kind: ActionCharacter, Character: CardAmbassador, AmbassadorHand: Hand{CardContessa, CardContessa}}

	is.NoErr(g.DoAction())
	is.Equal(len(g.deck), 2)

}

func TestGenerateDeck(t *testing.T) {
	is := is.New(t)

	want, have := [15]uint8{}, [15]uint8{}
	is.Equal(copy(want[:], shuffleCards(normalDeck[:])), 15)
	is.Equal(copy(have[:], shuffleCards(normalDeck[:])), 15)

	is.True(want != have)
}

func TestGameShuffle(t *testing.T) {
	g := &Game{}
	g.deck = shuffleCards(normalDeck[:])

	is := is.New(t)

	oldDeck, newDeck := [15]uint8{}, [15]uint8{}
	is.Equal(copy(oldDeck[:], g.deck), 15)

	g.Shuffle()

	is.Equal(15, copy(newDeck[:], g.deck))

	is.True(oldDeck != newDeck)
}

func TestGameDrawCards(t *testing.T) {
	g := &Game{}
	g.deck = shuffleCards(normalDeck[:])

	is := is.New(t)

	drawn := g.DrawCards(2)
	is.Equal(len(drawn), 2)
	last := g.deck[:2]
	is.Equal(drawn, last)
	is.Equal(len(g.deck), 13)
}
