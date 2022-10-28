package game

import (
	"fmt"
	"sync"
)

var (
	ErrInvalidCounterClaim = fmt.Errorf("invalid counter claim")
	ErrInvalidPlayer       = fmt.Errorf("invalid or dead player")
	ErrInvalidCharacter    = fmt.Errorf("invalid character")
	ErrInvalidParameters   = fmt.Errorf("invalid parameters")
)

type Game struct {
	players [5]*Player
	turn    int
	max     int
	claims  []*Claim
	mtx     sync.Mutex
}

func NewGame(pl [5]*Player) (*Game, error) {
	g := &Game{players: [5]*Player{}, turn: 0}

	for k, v := range pl {
		if v != nil {
			if v.IsDead() {
				return nil, ErrInvalidPlayer
			}
			g.max = k + 1
		}
	}

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
	g.mtx.Lock()
	defer g.mtx.Unlock()

	arr, err := addClaim(g.claims, c)
	if err != nil {
		return err
	}

	g.claims = arr

	return nil
}

func (g *Game) NextTurn() {
	g.turn = nextTurn(g.turn, g.max)
}
