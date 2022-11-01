package game

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrInvalidCounterClaim        = fmt.Errorf("invalid counter claim")
	ErrInvalidPlayer              = fmt.Errorf("invalid or dead player")
	ErrInvalidCharacter           = fmt.Errorf("invalid character")
	ErrInvalidParameters          = fmt.Errorf("invalid parameters")
	ErrInvalidAction              = fmt.Errorf("invalid action")
	ErrInvalidActionFrozen        = fmt.Errorf("cannot create action because it is frozen by a claim")
	ErrInvalidPlayerAmount        = fmt.Errorf("players must be 2 to 5")
	ErrInvalidClaim               = fmt.Errorf("invalid claim")
	ErrInvalidClaimOngoing        = fmt.Errorf("claim is still on going")
	ErrInvalidClaimHasNotFinished = fmt.Errorf("claim hasn't finished")
	ErrInvalidClaimFinished       = fmt.Errorf("claim has already finished")
	ErrInvalidClaimNotChallenged  = fmt.Errorf("claim has not been challenge")
	ErrInvalidClaimProvenAlready  = fmt.Errorf("claim has been proven already")
	ErrInvalidArr                 = fmt.Errorf("array is either nil or is empty")
	ErrInvalidGame                = fmt.Errorf("game was not initiated properly")
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
//          using Game. Since much of Game's internal design relies
//          heavily on an external package to translate client commands
//          into game actions.
type Game struct {
	deck       []uint8
	deckMtx    sync.Mutex
	players    [5]*Player
	turn       *Notifier
	max        int
	claim      *claim
	claimMtx   sync.Mutex
	action     [2]*Action
	actionMtx  sync.Mutex
	history    []Action
	historyMtx sync.Mutex
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var normalDeck = [15]uint8{CardDuke, CardDuke, CardDuke,
	CardContessa, CardContessa, CardContessa,
	CardAssassin, CardAssassin, CardAssassin,
	CardAmbassador, CardAmbassador, CardAmbassador,
	CardCaptain, CardCaptain, CardCaptain}

// Durstenfeld's version of the fisher-yates algorithm
func shuffleCards(givenDeck []uint8) []uint8 {
	// copy the normal deck
	// deck := append([]uint8{}, normalDeck[:]...)
	deck := append([]uint8{}, givenDeck...)

	i := len(givenDeck) - 1
	for i > 0 {
		shuffledIndex := rand.Intn(15)

		deck[shuffledIndex], deck[i] = deck[i], deck[shuffledIndex]
		i--
	}

	return deck
}

// NewGame creates a new game via providing it with a slice of players.
// The slice of players cannot contain less than 2 nil values, it must
// have at-least 2 or more.
func NewGame(pl [5]*Player) (*Game, error) {
	g := &Game{players: pl}

	g.deck = shuffleCards(normalDeck[:])

	for k, v := range pl {
		if v != nil {
			g.max = k + 1
		} else {
			first, second := g.deck[0], g.deck[1]
			g.deck = g.deck[2:]

			// simple way to tell the history, hey two cards were given..
			// maybe, someday, this should be its own action instead of piggy
			// backing off the Ambassador's only good use but that's for
			// later discussion
			g.history = append(g.history, Action{
				Kind:            ActionCharacter,
				Character:       CardAmbassador,
				AmbassadorHand:  Hand{first, second},
				AmbassadorPlace: Hand{0, 1},
			})
		}
	}

	if g.max < 2 {
		return nil, ErrInvalidPlayerAmount
	}

	g.turn = NewNotifier()
	g.turn.Set(0)

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
func validateClaimAndItsPlayer(arr []*Player, c *claim) (int, error) {
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
func (g *Game) validateClaimAndItsPlayer(c *claim) (int, error) {
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
func (g *Game) addClaimToHistory(c *claim, authorId uint8) {
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
func (g *Game) Claim(author *Player, character uint8) error {
	g.claimMtx.Lock()
	defer g.claimMtx.Unlock()
	if g.claim != nil {
		return ErrInvalidClaimOngoing
	}

	c := &claim{author: author, character: character}
	if err := c.IsValid(); err != nil {
		return err
	}

	index, err := g.validateClaimAndItsPlayer(c)
	if err != nil {
		return err
	}

	g.addClaimToHistory(c, uint8(index))
	g.claim = c

	return nil
}

// ClaimPass makes the underlying claim pass, allowing future character
// actions to succeed.
func (g *Game) ClaimPass() error {
	g.claimMtx.Lock()
	defer g.claimMtx.Unlock()

	if g.claim == nil {
		return ErrInvalidClaim
	} else if g.claim.IsFinished() {
		return ErrInvalidClaimFinished
	}

	index, _ := g.validateClaimAndItsPlayer(g.claim)

	g.claim.Pass()
	g.addClaimToHistory(g.claim, uint8(index))

	return nil
}

// ClaimChallenge challenges the underlying claim. This function will freeze
// any calls to Action until the Claim has been proven and the appropriate
// player was punished.
func (g *Game) ClaimChallenge(challenger *Player) error {
	g.claimMtx.Lock()
	defer g.claimMtx.Unlock()

	if g.claim == nil {
		return ErrInvalidClaim
	} else if g.claim.IsFinished() {
		return ErrInvalidClaimFinished
	} else if g.claim.author == challenger {
		return ErrInvalidActionSamePlayer
	}

	challengerIndex := findPlayerByPntr(g.players[:], challenger)
	if challengerIndex < 0 {
		return ErrInvalidPlayer
	}

	index, _ := g.validateClaimAndItsPlayer(g.claim)

	historyItem := g.claim.Action(uint8(index))

	val := historyItem.AuthorID
	historyItem.AgainstID, historyItem.against = &val, historyItem.author
	historyItem.AuthorID, historyItem.author = uint8(challengerIndex), challenger

	g.claim.Challenge()
	g.addActionToHistory(historyItem)

	return nil
}

// ClaimProve provides the challenged claim with proof. The proof could
// be correct or incorrect, the result of this evaluation is stored in
// the first parameter.
//
// If the proof matched the claim; the challenger gets punished; if not;
// the claimant gets punished.
func (g *Game) ClaimProve(character uint8) (bool, error) {
	g.claimMtx.Lock()
	defer g.claimMtx.Unlock()

	if g.claim == nil {
		return false, ErrInvalidClaim
	} else if !g.claim.IsFinished() {
		return false, ErrInvalidClaimHasNotFinished
	} else if _, challenge := g.claim.Results(); challenge == nil {
		return false, ErrInvalidClaimNotChallenged
	}

	originalCharacter := uint8(0)

	g.historyMtx.Lock()

	lastAction := g.history[len(g.history)-1]
	lastAction.Kind = ActionClaimProof
	originalCharacter = lastAction.Character
	lastAction.Character = character
	lastAction.author, lastAction.against = lastAction.against, lastAction.author
	againstId := lastAction.AuthorID
	lastAction.AuthorID, lastAction.AgainstID = *lastAction.AgainstID, &againstId
	g.history = append(g.history, lastAction)
	g.historyMtx.Unlock()

	succeed := originalCharacter == character
	g.claim.Prove(originalCharacter == character)

	return succeed, nil
}

// Action is a function that sets a Game's underlying Action. Once an action
// or two have been set, Actions are executed with Game.DoAction.
//
// If two actions have been set, in the same breath, before calling DoAction,
// then DoAction only executes the last Action.
func (g *Game) Action(a Action) error {
	if err := a.setPlayer(g.players[:]); err != nil {
		return err
	}

	if err := a.IsValid(); err != nil {
		return err
	}

	g.claimMtx.Lock()
	c := g.claim
	g.claimMtx.Unlock()
	if c != nil {
		if err := a.validClaim(c); err != nil {
			return err
		}

		succeed, challenge := c.Results()
		if succeed != nil && challenge != nil && a.Kind != ActionClaimPunishment {
			return ErrInvalidActionKind
		}
	}

	g.actionMtx.Lock()
	defer g.actionMtx.Unlock()

	if g.action[0] == nil {
		g.action[0] = &Action{}
		*g.action[0] = a
	} else if g.action[1] == nil {
		if !IsValidCounterAction(*g.action[0], a) {
			return ErrInvalidCounterClaim
		}

		g.action[1] = &Action{}
		*g.action[1] = a
	} else {
		return ErrInvalidAction
	}

	return nil
}

// DoAction is a function that executes only the last Action and clears
// the "stack" of actions.
func (g *Game) DoAction() error {
	g.actionMtx.Lock()
	defer g.actionMtx.Unlock()

	var act *Action
	if g.action[1] != nil {
		act = g.action[1]
	} else if g.action[0] != nil {
		act = g.action[0]
	} else {
		return ErrInvalidAction
	}

	g.historyMtx.Lock()
	g.history = append(g.history, *act)
	g.historyMtx.Unlock()

	act.do()

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

// Shuffle shuffles the game deck
func (g *Game) Shuffle() {
	g.deckMtx.Lock()
	g.deck = shuffleCards(g.deck)
	g.deckMtx.Unlock()
}

// DrawCards draws cards from the Game's deck.
func (g *Game) DrawCards(n uint8) []uint8 {
	g.deckMtx.Lock()
	defer g.deckMtx.Unlock()

	return g.deck[:n]
}
