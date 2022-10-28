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
	c := Claim{author: &Player{Hand: Hand{CardAssassin, CardEmpty}}, character: CardAssassin}
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
}

func TestClaimFinished(t *testing.T) {
	is := is.New(t)
	c := &Claim{}
	
	is.Equal(c.Finished(), false)
	pntr := false
	
	is.Equal((&Claim{succeed: &pntr}).Finished(), true)
	is.Equal((&Claim{challenge: &pntr}).Finished(), true)
}

func TestClaimWait(t *testing.T) {
	c := &Claim{}
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
}

func TestClaimPassOrChallenge(t *testing.T) {
	is := is.New(t)
	
	c := &Claim{}
	c.waitCalled = true
	c.wg.Add(1)

	var val *bool
	go func() {
		time.Sleep(time.Millisecond)
		val = c.passOrChallenge()
	}()
	c.wg.Wait()

	is.True(val != nil)
	c.challenge = val
	is.True(c.passOrChallenge() == nil)
}

func TestClaimPassOrChallengeWait(t *testing.T) {
	c := &Claim{}
	go func() {
		time.Sleep(time.Millisecond)
		c.passOrChallenge()
	}()
	c.Wait()
}

func TestClaimChallenge(t *testing.T) {
	c := &Claim{}
	go func() {
		time.Sleep(time.Millisecond)
		c.Challenge()
	}()
	c.Wait()

	is := is.New(t)
	is.True(c.challenge != nil)

	c.Challenge()
	is.True(c.passOrChallenge() == nil)
	is.True(c.challenge != nil)
}

func TestClaimPass(t *testing.T) {
	c := &Claim{}
	go func() {
		time.Sleep(time.Millisecond)
		c.Pass()
	}()
	c.Wait()

	is := is.New(t)
	is.True(c.succeed != nil)

	c.Pass()
	is.True(c.passOrChallenge() == nil)
	is.True(c.succeed != nil)
}
