package game

import (
	"testing"
	"time"
)

func TestIsValidCounterClaim(t *testing.T) {
	if !IsValidCounterClaim(CardAssassin, CardContessa) {
		t.Fatalf("Assassin should be countered by Contessa")
	}
	
	if !IsValidCounterClaim(CardCaptain, CardAmbassador) ||
				!IsValidCounterClaim(CardCaptain, CardCaptain) {
		t.Fatalf("Captain should be countered by an Ambassador or a Captain")
	}
	
	if IsValidCounterClaim(CardContessa, CardAssassin) {
		t.Fatalf("A Contessa cannot be countered by an assassin")
	}
}

func TestClaimValid(t *testing.T) {
	c := Claim{author: &Player{Hand: Hand{CardAssassin, CardEmpty}}, character: CardAssassin}
	if err := c.IsValid(); err != nil {
		t.Fatalf("new claim should work")
	}
	
	c.author = &Player{}
	if err := c.IsValid(); err != ErrInvalidPlayer {
		t.Fatalf("Claim.IsValid() doesn't check if player is dead or not: %v", err)
	}
	
	c.author = &Player{Hand: Hand{CardEmpty, CardAssassin}}
	c.character = 0
	if err := c.IsValid(); err != ErrInvalidCharacter {
		t.Fatalf("Claim.IsValid() doesn't check if character is invalid or not: %v", err)
	}
}

func TestNewClaim(t *testing.T) {
	if _, err := NewClaim(nil, 0); err == nil {
		t.Fatalf("NewClaim should return err because Claim is invalid")
	}
}

func TestClaimFinished(t *testing.T) {
	c := &Claim{}
	if c.Finished() {
		t.Fatalf("claim shouldn't be finished")
	}
	
	pntr := false
	if (Claim{succeed: &pntr}).Finished() == false || (Claim{challenge: &pntr}).Finished() == false {
		t.Fatalf("Claim.Finished doesn't account for either challenge field or succeed field")
	}
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
}

func TestClaimPassOrChallenge(t *testing.T) {
	c := &Claim{}
	c.waitCalled = true
	c.wg.Add(1)
	
	var val *bool
	go func() {
		time.Sleep(time.Millisecond)
		val = c.passOrChallenge()
	}()
	c.wg.Wait()
	
	if val == nil {
		t.Fatalf("challenge is not set..")
	}
	
	b := false
	c.challenge = &b
	if val := c.passOrChallenge(); val != nil {
		t.Fatalf("Claim has already finished but passOrChallenge() returns non-nil: %v", val)
	}
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
	
	if c.challenge == nil {
		t.Fatalf("Challenge() doesn't set field challenge")
	}
	
	c.Challenge()
	if c.passOrChallenge() == nil && c.challenge == nil {
		t.Fatalf("Challenge() sets challenge to nil once passOrChallenge is nil")
	}
}

func TestClaimPass(t *testing.T) {
	c := &Claim{}
		go func() {
		time.Sleep(time.Millisecond)
		c.Pass()
	}()
	c.Wait()
	
	if c.succeed == nil {
		t.Fatalf("Challenge() doesn't set field challenge")
	}
	
	c.Pass()
	if c.passOrChallenge() == nil && c.succeed == nil {
		t.Fatalf("Challenge() sets challenge to nil once passOrChallenge is nil")
	}
}