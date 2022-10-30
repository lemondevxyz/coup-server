package game

import "sync"

// IsValidCounterClaim is a function that determines whether a Claim was
// a valid counter claim. Put simply, there are cards that counter other
// cards. For example, a Contessa can counter an assassin. However,
// some cards cannot counter others, like how a Duke cannot counter an
// Ambassador.
//
// IsValidCounterClaim ensures that the data inputted by the user for a
// counter claim is valid.
func IsValidCounterClaim(character, counter uint8) bool {
	switch character {
	case CardAssassin:
		return counter == CardContessa
	case CardCaptain:
		return counter == CardAmbassador ||
			counter == CardCaptain
	}

	return false
}

// Claim is a data structure that tells the internal packages; "Hey I
// am this character". A claim does not have to be actually true; it
// could be a lie.
//
// For example, a player could not have an ambassador
// card and yet they could still claim it.
//
// If the claim passes, the action could proceed. If not and the player
// is challenged, the player must provide a card that either matches
// their claim or *choose* not to.
//
// This means that a player could technically have a card that they
// claim to have but instead choose not to show it.
//
// Claim will *never* have a timer implementation under it. It is meant
// to be manipulated functionally.
//
// Note: Claims should not be modified but instead used only once.
//       Any counter claim, used to defend the player, should be done
//       through creating another claim.
type Claim struct {
	author     *Player
	character  uint8
	succeed    *bool
	challenge  *bool
	wg         sync.WaitGroup
	mtx        sync.Mutex
	once       sync.Once
	waitCalled bool
}

// NewClaim is a function that creates a valid Claim or return an error.
//
// NewClaim only returns an error if one of the parameters is invalid.
// NewClaim uses Claim.IsValid() to check the validity of the claim.
//
// An empty claim is always invalid.
func NewClaim(player *Player, character uint8) (*Claim, error) {
	c := &Claim{author: player, character: character}
	if err := c.IsValid(); err != nil {
		return nil, err
	}

	return c, nil
}

// IsValid is a function that checks for the validity of the claim.
// Do note: This function doesn't check whether or not a player has a
// card. Instead, it checks if the player is nil or not nil but dead,
// and if the character of the claim is invalid.
//
// The underlying functions are Player.IsDead() and IsValidCard() for
// the player and the card/character respectively.
func (c *Claim) IsValid() error {
	if c.author == nil || c.author.IsDead() {
		return ErrInvalidPlayer
	}

	if !IsValidCard(c.character) {
		return ErrInvalidCharacter
	}

	return nil
}

// Wait is a function that doesn't return until Claim.Pass() or
// Claim.Challenge() has been called.
//
// If Claim.Challenge() or Claim.Pass() were called before Wait was called
// then Wait returns immediately.
//
// Do note: Wait is only meant to be called once; Subsequent calls will
//          return immediately.
func (c *Claim) Wait() {
	if c.IsFinished() {
		return
	}

	c.once.Do(func() {
		c.mtx.Lock()

		c.wg.Add(1)
		c.waitCalled = true
		c.mtx.Unlock()

		c.wg.Wait()
	})
}

// Chan is a function that returns a channel that gets closed if
// Claim.Pass() or Claim.Challenge() has been called.
//
// Do note: Chan() uses the function Wait() under the hood, so all of
//          Wait()'s side-effects apply to Chan()
func (c *Claim) Chan() chan struct{} {
	if c.IsFinished() {
		return nil
	}

	ch := make(chan struct{})
	go func(c *Claim, ch chan struct{}) {
		c.Wait()
		close(ch)
	}(c, ch)

	return ch
}

// Results is a function that returns two values: passed & challenge.
// Passed is only set if Claim.Pass() was called. Challenge is only set
// if Claim.Challenge() was called.
//
// Do note: modification of these values won't impact the underlying
// structure's value.
func (c *Claim) Results() (passed *bool, challenge *bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.succeed != nil {
		val := *c.succeed
		passed = &val
	}

	if c.challenge != nil {
		val := *c.challenge
		challenge = &val
	}

	return
}

// IsFinished is a function that returns if Pass() or Challenge()
// have been called.
func (c *Claim) IsFinished() bool {
	if c.succeed != nil || c.challenge != nil {
		return true
	}

	return false
}

// passOrChallenge is a helper function that returns a pointer with
// a value of true if the Claim was not finished. If it was finished,
// it returns nil.
//
// passOrChallenge is mainly used to reduce code duplication between
// Claim.Pass() and Claim.Challenge().
func (c *Claim) passOrChallenge() *bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.IsFinished() {
		return nil
	}
	b := true

	if c.waitCalled {
		c.wg.Done()
		c.waitCalled = false
	}

	return &b
}

// Challenge sets the Claim's challenge value to true unless Claim.Pass()
// was called first.
func (c *Claim) Challenge() {
	if val := c.passOrChallenge(); val != nil {
		c.challenge = val
	}
}

// Pass sets the Claim's pass value to true unless Claim.Pass() was
// called first.
func (c *Claim) Pass() {
	if val := c.passOrChallenge(); val != nil {
		c.succeed = val
	}
}

// Action returns an Action that's used to store a Claim in a history
// array.
func (c *Claim) Action(authorid uint8) Action {
	a := Action{
		AuthorID:  authorid,
		Character: c.character,
	}

	if !c.IsFinished() {
		a.Kind = ActionClaim
	} else {
		if c.challenge != nil {
			a.Kind = ActionClaimChallenge
		} else {
			a.Kind = ActionClaimPassed
		}
	}

	return a
}
