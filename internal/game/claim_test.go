package game

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestIsValidCounterClaim(t *testing.T) {
	is := is.New(t)
	is.Equal(IsValidCounterClaim(CardAssassin, CardContessa), true)
	is.Equal(IsValidCounterClaim(CardCaptain, CardCaptain), true)
	is.Equal(IsValidCounterClaim(CardCaptain, CardAmbassador), true)
	is.Equal(IsValidCounterClaim(CardContessa, CardAssassin), false)
}

func TestClaimValid(t *testing.T) {
	is := is.New(t)
	c := claim{author: &Player{Hand: Hand{CardAssassin, CardEmpty}}, character: CardAssassin}
	is.NoErr(c.IsValid())
	c.author = &Player{}
	is.Equal(c.IsValid(), ErrInvalidPlayer)
	c.author = &Player{Hand: Hand{CardEmpty, CardAssassin}}
	c.character = 0
	is.Equal(c.IsValid(), ErrInvalidCharacter)
}

func TestNewClaim(t *testing.T) {
	is := is.New(t)
	_, err := NewClaim(nil, 0)
	is.True(err != nil)

	_, err = NewClaim(&Player{Hand: Hand{0: CardCaptain}}, CardContessa)
	is.NoErr(err)
}

func TestClaimIsFinished(t *testing.T) {
	is := is.New(t)
	c := &claim{}

	is.Equal(c.IsFinished(), false)
	pntr := false

	is.Equal((&claim{succeed: &pntr}).IsFinished(), true)
	is.Equal((&claim{challenge: &pntr}).IsFinished(), true)
}

/*
func TestClaimWait(t *testing.T) {
	c := &claim{}
	v := make(chan struct{})
	go func() {
		time.Sleep(time.Millisecond)
		c.wg.Done()
		close(v)
	}()

	c.Wait()
	<-v

	is := is.New(t)
	is.True(c.waitCalled)

	c.Wait()

	c = &claim{}
	c.Challenge()

	c.Wait()
}
*/

func TestClaimPassOrChallenge(t *testing.T) {
	is := is.New(t)

	c := &claim{}
	//c.wg.Add(1)

	val := c.passOrChallenge()

	is.True(val != nil)
	c.challenge = val
	is.True(c.passOrChallenge() == nil)
}

func TestClaimPassOrChallengeWait(t *testing.T) {
	c := &claim{}
	go func() {
		time.Sleep(time.Millisecond)
		c.passOrChallenge()
	}()
	//c.Wait()
}

func TestClaimChallenge(t *testing.T) {
	c := &claim{}

	c.Challenge()
	is := is.New(t)
	is.True(c.challenge != nil)

	c.Challenge()
	is.True(c.passOrChallenge() == nil)
	is.True(c.challenge != nil)
}

func TestClaimPass(t *testing.T) {
	c := &claim{}
	c.Pass()

	is := is.New(t)
	is.True(c.succeed != nil)

	c.Pass()
	is.True(c.passOrChallenge() == nil)
	is.True(c.succeed != nil)
}

func TestClaimProve(t *testing.T) {
	c := &claim{}

	is := is.New(t)
	c.Challenge()
	c.Prove(false)
	is.True(c.succeed != nil)
	is.Equal(*c.succeed, false)
}

func TestClaimResults(t *testing.T) {
	c := &claim{}
	is := is.New(t)

	passed, challenge := c.Results()

	is.Equal(passed, nil)
	is.Equal(challenge, nil)

	b := false
	c.succeed = &b

	passed, _ = c.Results()
	is.True(passed != c.succeed)
	is.True(*passed == false)

	b2 := false
	c.challenge = &b2

	passed, challenge = c.Results()
	is.True(challenge != passed)
	is.True(challenge != c.challenge)
	is.True(*challenge == false)
}

/*
func TestClaimChan(t *testing.T) {
	c := &claim{}

	go func() {
		time.Sleep(time.Millisecond)
		c.Challenge()
	}()

	<-c.Chan()
	//asd

	c.challenge = nil
	<-c.Chan()

	v := false
	c.challenge = &v

	is := is.New(t)
	is.Equal(c.Chan(), nil)
}
*/

func TestClaimAction(t *testing.T) {
	c := &claim{character: CardAmbassador}

	newAction := func(id uint8, kind uint8, character uint8) Action {
		return Action{
			AuthorID:  id,
			Kind:      kind,
			Character: character,
		}
	}

	is := is.New(t)
	is.Equal(c.Action(0), newAction(0, ActionClaim, CardAmbassador))
	c.Pass()
	is.Equal(c.Action(0), newAction(0, ActionClaimPassed, CardAmbassador))

	c = &claim{character: CardAmbassador}
	c.Challenge()
	is.Equal(c.Action(0), newAction(0, ActionClaimChallenge, CardAmbassador))
	c.succeed = c.challenge
	is.Equal(c.Action(0), newAction(0, ActionClaimProof, CardAmbassador))
}
