package game

import (
	"fmt"
	"sync"
	"time"
)

var (
	ErrInvalidCounterClaim = fmt.Errorf("invalid counter claim")
	ErrInvalidPlayer       = fmt.Errorf("invalid or dead player")
	ErrInvalidCharacter    = fmt.Errorf("invalid character")
	ErrInvalidParameters   = fmt.Errorf("invalid parameters")
	ErrInvalidAction       = fmt.Errorf("invalid action")
	ErrInvalidPlayerAmount = fmt.Errorf("players must be 2 to 5")
)

type Game struct {
	players   [5]*Player
	turn      *Notifier
	max       int
	claim     *Notifier
	action    [2]*Action
	actionMtx sync.Mutex
}

func NewGame(pl [5]*Player) (*Game, error) {
	g := &Game{players: [5]*Player{}}

	for k, v := range pl {
		if v != nil {
			if v.IsDead() {
				return nil, ErrInvalidPlayer
			}
			g.max = k + 1
		}
	}

	if g.max < 2 {
		return nil, ErrInvalidPlayerAmount
	}

	g.claim, g.turn = NewNotifier(), NewNotifier()
	g.turn.Set(0)

	//a asd
	return g, nil
}

func nextTurn(i int, max int) int {
	if max == 0 {
		return -1
	}

	if i+1 >= max {
		return 0
	}

	return i + 1
}

func addClaim(arr []*Claim, claim *Claim) ([]*Claim, error) {
	if arr == nil || claim == nil {
		return nil, ErrInvalidParameters
	}

	if err := claim.IsValid(); err != nil {
		return nil, err
	}

	if len(arr) > 0 {
		if !IsValidCounterClaim(arr[len(arr)-1].character, claim.character) {
			return nil, ErrInvalidCounterClaim
		}
	}

	return append(arr, claim), nil
}

func (g *Game) Claim(c *Claim) error {
	if err := c.IsValid(); err != nil {
		return err
	}

	g.claim.Set(c)
	g.claim.Announce()
	return nil
}

func (g *Game) ClaimSubscribe() (time.Time, <-chan struct{}) { return g.claim.Subscribe() }
func (g *Game) ClaimUnsubscribe(val time.Time)               { g.claim.Unsubscribe(val) }
func (g *Game) ClaimGet() *Claim                             { return g.claim.Get().(*Claim) }

func (g *Game) Action(a *Action) error {
	g.actionMtx.Lock()
	defer g.actionMtx.Unlock()

	if g.action[0] == nil {
		g.action[0] = a
	} else if g.action[1] == nil {
		g.action[1] = a
	} else {
		return ErrInvalidAction
	}

	return nil
}

func (g *Game) DoAction() error {
	g.actionMtx.Lock()
	defer g.actionMtx.Unlock()
	if g.action[1] != nil {
		g.action[1].Do()
	} else if g.action[0] != nil {
		g.action[0].Do()
	} else {
		return ErrInvalidAction
	}

	g.action[0], g.action[1] = nil, nil

	return nil
}

func (g *Game) NextTurn() {
	g.turn.Set(nextTurn(g.turn.Get().(int), g.max))
	g.turn.Announce()
}

func (g *Game) TurnGet() int                                { return g.turn.Get().(int) }
func (g *Game) TurnSubscribe() (time.Time, <-chan struct{}) { return g.turn.Subscribe() }
func (g *Game) TurnUnsubscribe(val time.Time)               { g.turn.Unsubscribe(val) }
