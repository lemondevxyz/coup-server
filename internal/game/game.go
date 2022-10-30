package game

import (
	"fmt"
	"sync"
	"time"
)

var (
	ErrInvalidCounterClaim       = fmt.Errorf("invalid counter claim")
	ErrInvalidPlayer             = fmt.Errorf("invalid or dead player")
	ErrInvalidCharacter          = fmt.Errorf("invalid character")
	ErrInvalidParameters         = fmt.Errorf("invalid parameters")
	ErrInvalidAction             = fmt.Errorf("invalid action")
	ErrInvalidActionFrozen       = fmt.Errorf("cannot create action because it is frozen by a claim")
	ErrInvalidPlayerAmount       = fmt.Errorf("players must be 2 to 5")
	ErrInvalidClaim              = fmt.Errorf("invalid claim")
	ErrInvalidClaimFinished      = fmt.Errorf("claim has already finished")
	ErrInvalidClaimNotChallenged = fmt.Errorf("claim has not been challenge")
	ErrInvalidClaimProvenAlready = fmt.Errorf("claim has been proven already")
	ErrInvalidArr                = fmt.Errorf("array is either nil or is empty")
	ErrInvalidGame               = fmt.Errorf("game was not initiated properly")
)

// Game is a data structure that essentially connects all the loose data
// structures together. It uses, Action, Claim, Notifier and Player to
// essentially run the game.
//
// Whilst Game is meant to essentially run the Game from scratch; it
// doesn't implement features like Timers which are essential in the
// game. Game is meant to be a data structure than can be used both
// synchronously and asynchronously.
//
// An Empty Game is invalid. Use NewGame to generate a new game.
//
// Do note: You are meant to have up keep of the slice of players when
//          using Game. Since, much of Game's internal design relies
//          heavily on an external package to translate client commands
//          into game actions.
type Game struct {
	players    [5]*Player
	turn       *Notifier
	max        int
	claim      *Notifier
	action     [2]*Action
	actionMtx  sync.Mutex
	history    []Action
	historyMtx sync.Mutex
}

// NewGame creates a new game via providing it with a slice of players.
// The slice of players cannot contain less than 2 nil values, it must
// have at-least 2 or more.
func NewGame(pl [5]*Player) (*Game, error) {
	g := &Game{players: pl}

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
	g.claim.Set(nil)

	return g, nil
}

// nextTurn is a function that adds 1 to the i value. If i + 1 is
// bigger than the max value, it returns 0.
//
// If max = 0; it returns -1.
func nextTurn(i int, max int) int {
	if max == 0 {
		return -1
	}

	if i+1 >= max {
		return 0
	}

	return i + 1
}

// findPlayerByPntr is a function that tries to a locate a player in a
// slice by comparing pointers. If found, the player's index is returned.
//
// If not, -1 is returned.
func findPlayerByPntr(arr []*Player, player *Player) int {
	if player == nil {
		return -1
	}

	for index, target := range arr {
		if target == player {
			return index
		}
	}

	return -1
}

// validateClaimAndItsPlayer is a function that takes in a slice of
// players and a claim; it then checks if any of the parameters are
// invalid and returns error; If not, it tries to locate the author in
// claim and returns it index if found.
func validateClaimAndItsPlayer(arr []*Player, c *Claim) (int, error) {
	if len(arr) == 0 {
		return -1, ErrInvalidArr
	}
	if c == nil {
		return -1, ErrInvalidClaim
	}

	index := findPlayerByPntr(arr, c.author)
	if index < 0 {
		return -1, ErrInvalidPlayer
	}

	return index, nil
}

// validateClaimAndItsPlayer is essentially a wrapper over the non-game
// function validateClaimAndItsPlayer.
func (g *Game) validateClaimAndItsPlayer(c *Claim) (int, error) {
	return validateClaimAndItsPlayer(g.players[:], c)
}

// addActionToHistory is a function that adds an action to the history
// slice of the game. It also locks and unlocks the history's mutex so
// that operations are safe when used asynchronously.
func (g *Game) addActionToHistory(a Action) {
	g.historyMtx.Lock()
	g.history = append(g.history, a)
	g.historyMtx.Unlock()
}

// addClaimToHistory is a wrapper around addActionToHistory and
// Claim.Action.
func (g *Game) addClaimToHistory(c *Claim, authorId uint8) {
	g.addActionToHistory(c.Action(authorId))
}

// Claim is a function that sets the current claim for the Game.
//
// To understand how Claim works; a Claim is essentially a saying like
// this: "Hey, I have a CardContessa". The saying could be true or false,
// but to really evaluate if it is true or not, another player *must*
// challenge the Claim through Game.ClaimChallenge.
//
// Once the Claim has been challenged; the original player who made the
// Claim must provide proof through Game.ClaimProof. The proof could match
// the original claim's saying or not. Meaning, that players could claim
// to have a card and when they are asked to prove it, they could choose
// to not prove it.
//
// If the Claim's proof was unsatisfactory, the underlying Action of the
// claim will not succeed. Instead, the player who challenge the original
// player who made the claim can choose which card to remove from the
// original player.
//
// If the Claim's proof was satisfactory, the original player can choose
// which card they want to remove from the player who made the challenge.
// Also, the original Claim's action will succeed.
//
// However, if a Claim has not been challenged for a period of time,
// ClaimPass should be called. Once ClaimPass has been called, it is also
// advisable to create the action and finally call DoAction.
//
// Lastly, Claims do not necessarily have to be challenged to prevent
// them. One could always counter-claim if the original Action is counterable.
// Like, for example, an Assassin and a Contessa or a Captain and another
// Captain.
func (g *Game) Claim(c *Claim) error {
	if err := c.IsValid(); err != nil {
		return err
	}

	if c.IsFinished() {
		return ErrInvalidClaimFinished
	}

	index, err := g.validateClaimAndItsPlayer(c)
	if err != nil {
		return err
	}

	g.addClaimToHistory(c, uint8(index))

	g.claim.Set(c)
	g.claim.Announce()

	return nil
}

// ClaimPass is a function that makes the underlying claim pass.
//
// This makes the action of the claim succeed unless it is countered
// by another claim.
//
// It is good practice to call ClaimPass after a duration has elapsed.
//
// It is also good practice to choose different durations at key points
// of the game. For example, giving a player a 20 second duration to
// validate a Claim at the end of a game, but giving them a 5 second
// duration at the start of a game.
func (g *Game) ClaimPass() error {
	c, err := g.ClaimGet()
	if err != nil {
		return err
	}
	if c.IsFinished() {
		return ErrInvalidClaimFinished
	}

	index, err := g.validateClaimAndItsPlayer(c)
	if err != nil {
		return err
	}

	c.Pass()
	g.addClaimToHistory(c, uint8(index))

	return nil
}

// ClaimChallenge challenges the current claim.
//
// After ClaimChallenge all subsequent calls to Action and DoAction are
// frozen. They can be later made unfrozen through calling ClaimProof.
func (g *Game) ClaimChallenge(challenger *Player) error {
	c, err := g.ClaimGet()
	if err != nil {
		return err
	}
	if c == nil {
		return ErrInvalidClaim
	}

	if c.IsFinished() {
		return ErrInvalidClaimFinished
	}

	againstId, err := g.validateClaimAndItsPlayer(c)
	if err != nil {
		return err
	}

	authorId := findPlayerByPntr(g.players[:], challenger)
	if authorId < 0 {
		return ErrInvalidPlayer
	}

	actionAgainstId := uint8(againstId)

	g.addActionToHistory(Action{
		AuthorID:  uint8(authorId),
		AgainstID: &actionAgainstId,
		Kind:      ActionClaimChallenge,
		Character: c.character,
	})

	c.Challenge()

	return nil
}

// ClaimProof adds the claim proof if the claim was challenged.
func (g *Game) ClaimProof(character uint8) error {
	c, err := g.ClaimGet()
	if err != nil {
		return err
	}
	if c == nil {
		return ErrInvalidClaim
	}

	if !IsValidCard(character) {
		return ErrInvalidCharacter
	}

	_, challenge := c.Results()
	if challenge == nil || !*challenge {
		return ErrInvalidClaimNotChallenged
	}

	g.historyMtx.Lock()
	defer g.historyMtx.Unlock()

	act := g.history[len(g.history)-1]
	if act.Kind == ActionClaimProof {
		return ErrInvalidClaimProvenAlready
	}

	g.history = append(g.history, Action{
		AuthorID:  act.AuthorID,
		AgainstID: act.AgainstID,
		Kind:      ActionClaimProof,
		Character: character,
	})

	return nil
}

// ClaimWasProven returns true, nil if the claim was proven. Else, it returns
// false, nil. Or, in-case there wasn't any claim to begin with, it returns
// false, InvalidClaim.
func (g *Game) ClaimWasProven() (bool, error) {
	c, err := g.ClaimGet()
	if err != nil {
		return false, err
	}
	if c == nil {
		return false, ErrInvalidClaim
	}

	g.historyMtx.Lock()
	defer g.historyMtx.Unlock()

	if len(g.history) < 2 {
		return false, ErrInvalidClaim
	}

	lastAction := g.history[len(g.history)-1]
	beforeLastAction := g.history[len(g.history)-2]
	if lastAction.Kind == ActionClaimProof &&
		beforeLastAction.Kind == ActionClaimChallenge {
		if beforeLastAction.Character == lastAction.Character {
			return true, nil
		} else {
			return false, nil
		}
	}

	return false, ErrInvalidClaim
}

// ClaimSubscribe is a function that returns the underlying Claim
// Notifier Subscribe method.
func (g *Game) ClaimSubscribe() (t time.Time, c <-chan struct{}, err error) {
	defer func() {
		if recover() != nil {
			err = ErrInvalidGame
		}
	}()

	t, c = g.claim.Subscribe()

	return
}

// ClaimUnsubscribe is a function that returns the underlying Claim
// Notifier Unsubscribe method.
func (g *Game) ClaimUnsubscribe(val time.Time) (err error) {
	defer func() {
		if recover() != nil {
			err = ErrInvalidGame
		}
	}()

	g.claim.Unsubscribe(val)

	return
}

// ClaimGet is a function that returns the underlying Claim Notifier Get
// method but with one difference: The value return will always be a
// pointer to a claim.
func (g *Game) ClaimGet() (c *Claim, err error) {
	defer func() {
		if recover() != nil {
			err = ErrInvalidGame
		}
	}()

	if g.claim.Get() == nil {
		c, err = nil, nil
	} else {
		c, err = g.claim.Get().(*Claim), nil
	}

	return
}

// Action is a function that sets a Game's underlying Action. Once an action
// or two have been set, Actions are executed with Game.DoAction.
//
// If two actions have been set, in the same breath, before calling DoAction,
// then DoAction only executes the last Action.
func (g *Game) Action(a *Action) error {
	if a == nil {
		return ErrInvalidAction
	}

	index := a.AuthorID
	if int(index) >= len(g.players) || g.players[index] == nil {
		return ErrInvalidActionAuthor
	}
	a.author = g.players[index]

	if a.AgainstID != nil {
		index := *a.AgainstID
		if int(index) >= len(g.players) || g.players[index] == nil {
			return ErrInvalidActionAgainst
		}

		a.against = g.players[index]
	}

	if err := a.IsValid(); err != nil {
		return err
	}

	g.actionMtx.Lock()
	defer g.actionMtx.Unlock()

	c, err := g.ClaimGet()
	if err != nil {
		return err
	}

	if c != nil {
		if c.IsFinished() {
			_, challenge := c.Results()

			if challenge != nil {
				return ErrInvalidActionFrozen
			} else {
				g.claim.Set(nil)
				g.claim.Announce()
			}
		} else {
			return ErrInvalidActionFrozen
		}
	}

	if g.action[0] == nil {
		g.action[0] = a
	} else if g.action[1] == nil {
		if !IsValidCounterAction(*g.action[0], *a) {
			return ErrInvalidCounterClaim
		}

		g.action[1] = a
	} else {
		return ErrInvalidAction
	}

	g.historyMtx.Lock()
	g.history = append(g.history, *a)
	g.historyMtx.Unlock()

	return nil
}

// DoAction is a function that executes only the last Action and clears
// the "stack" of actions.
func (g *Game) DoAction() error {
	g.actionMtx.Lock()
	defer g.actionMtx.Unlock()
	if g.action[1] != nil {
		g.action[1].do()
	} else if g.action[0] != nil {
		g.action[0].do()
	} else {
		return ErrInvalidAction
	}

	g.action[0], g.action[1] = nil, nil

	return nil
}

// NextTurn changes the turn and announce it.
func (g *Game) NextTurn() {
	g.turn.Set(nextTurn(g.turn.Get().(int), g.max))
	g.turn.Announce()
}

// TurnGet is a function that returns the underlying Claim Notifier Get
// method, but with one difference: the return variable will always be
// an int.
func (g *Game) TurnGet() (i int, err error) {
	defer func() {
		if recover() != nil {
			i, err = -1, ErrInvalidGame
		}
	}()

	return g.turn.Get().(int), nil
}

// TurnSubscribe is a function that returns the underlying Claim Subscribe
// method.
func (g *Game) TurnSubscribe() (t time.Time, c <-chan struct{}, err error) {
	defer func() {
		if recover() != nil {
			err = ErrInvalidGame
		}
	}()

	t, c = g.turn.Subscribe()

	return
}

// TurnUnsubscribe is a function that returns the underlying Claim Unsubscribe
// method.
func (g *Game) TurnUnsubscribe(val time.Time) (err error) {
	defer func() {
		if recover() != nil {
			err = ErrInvalidGame
		}
	}()

	g.turn.Unsubscribe(val)

	return
}
