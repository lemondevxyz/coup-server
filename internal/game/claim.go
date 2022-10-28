package game

import "sync"

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

func NewClaim(player *Player, character uint8) (*Claim, error) {
	c := &Claim{author: player, character: character}
	if err := c.IsValid(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c Claim) IsValid() error {
	if c.author == nil || c.author.IsDead() {
		return ErrInvalidPlayer
	}

	if !IsValidCard(c.character) {
		return ErrInvalidCharacter
	}

	return nil
}

func (c *Claim) Wait() {
	c.once.Do(func() {
		c.mtx.Lock()
		if c.waitCalled {
			c.mtx.Unlock()
			return
		}
	
		c.wg.Add(1)
		c.waitCalled = true
		c.mtx.Unlock()
	
		c.wg.Wait()
	})
}

func (c Claim) Finished() bool {
	if c.succeed != nil || c.challenge != nil {
		return true
	}

	return false
}

func (c *Claim) passOrChallenge() *bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.Finished() {
		return nil
	}
	b := true

	if c.waitCalled {
		c.wg.Done()
		c.waitCalled = false
	}

	return &b
}

func (c *Claim) Challenge() {
	if val := c.passOrChallenge(); val != nil {
		c.challenge = val
	}
}

func (c *Claim) Pass() {
	if val := c.passOrChallenge(); val != nil {
		c.succeed = val
	}
}
