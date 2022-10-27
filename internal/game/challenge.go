package game

import (
	"fmt"
	"time"
)

type Challenge struct {
	creator   *Player
	kind      uint8
	versus    *Player
	timestamp time.Time
	response  *Challenge
	succeed   *bool
}

func (c *Challenge) IsValid() error {
	if c.creator == nil {
		return ErrCreatorNil
	}

	if c.creator.IsDead() {
		return ErrCreatorDead
	}

	if !IsChallengeKindValid(c.kind) {
		return ErrBadKind
	}

	if c.timestamp.IsZero() {
		return ErrZeroTimestamp
	}

	return nil
}

func (c *Challenge) Succeed() bool {
	if c.response != nil {
		ret := !c.response.Succeed()
		c.succeed = &ret
		return ret
	}

	for c.succeed == nil {
	}

	return *c.succeed
}

func (c *Challenge) Unwrap() *Challenge {
	return c.response
}

func (c *Challenge) RespondWith(v *Challenge) {
	c.response = v
}

func (c *Challenge) Last() *Challenge {
	last := c
	for last != nil {
		if last.Unwrap() != nil {
			last = last.Unwrap()
		} else {
			break
		}
	}

	return last
}

func (c *Challenge) Passed() bool {
	if time.Now().After(c.timestamp) {
		return true
	}

	return false
}

var (
	ErrCreatorNil              = fmt.Errorf("creator is nil")
	ErrCreatorDead             = fmt.Errorf("creator is dead")
	ErrBadKind                 = fmt.Errorf("bad kind")
	ErrZeroTimestamp           = fmt.Errorf("zero timestamp")
	ErrChallengePassed         = fmt.Errorf("challenge passed")
	ErrChallengeForbiddenChild = fmt.Errorf("parent challenge kind cannot have a child")
	ErrChallengeNeedsParent    = fmt.Errorf("challenge needs parent challenge")
	ErrChallengeInvalidCounter = fmt.Errorf("challenge invalid counter")
	ErrChallengeHasResponse    = fmt.Errorf("challenge already has resposnse")
)

const (
	// Works on any action besides coup; income or financial aid.
	// Unstoppable.
	ChallengeConfront uint8 = iota + 1
	// Take 1 coin from the bank. Unstoppable
	ChallengeIncome
	// Take 2 coins. Preventable by a Duke
	ChallengeFinancialAid
	// Kill a player's card with 7 coins. Unstoppable
	ChallengeCoup
	// Kill a player's cord with 3 coins. Preventable
	// by a Contessa.
	ChallengeActionAssassin
	// Take 3 coins from the bank. Unstoppable
	ChallengeActionDuke
	// Steal from someone. Preventable by a Captain
	// or an ambassador.
	ChallengeActionCaptain
	// Look at two cards from the deck;
	// and swap (optionally) one or two cards with those two.
	// Unstoppable.
	ChallengeActionAmbassador
	// Duke denies anyone to take foreign aid. Unstoppable.
	ChallengeCounterDuke
	// Captain & Ambassador denies anyone to steal from them.
	// Unstopable.
	ChallengeCounterCaptain
	ChallengeCounterAmbassador
	// Contessa prevents an assassination made against her.
	ChallengeCounterContessa
	// Reveal Card can only be used against ChallengeConfront
	ChallengeCounterReveal
)

func IsCounterChallengeValid(action, counter uint8) bool {
	switch action {
	case ChallengeFinancialAid:
		return counter == ChallengeCounterDuke
	case ChallengeActionAssassin:
		return counter == ChallengeCounterContessa ||
			counter == ChallengeConfront
	case ChallengeActionCaptain:
		return counter == ChallengeConfront ||
			counter == ChallengeCounterCaptain ||
			counter == ChallengeCounterAmbassador
	// anything that cannot have a child to work
	case ChallengeIncome:
		fallthrough
	case ChallengeCoup:
		return false
	// anything that needs a parent action to work
	// or to not work
	case 0:
		switch counter {
		case ChallengeCounterDuke:
			fallthrough
		case ChallengeCounterContessa:
			fallthrough
		case ChallengeCounterAmbassador:
			fallthrough
		case ChallengeCounterCaptain:
			fallthrough
		case ChallengeConfront:
			return false
		default:
			return true
		}
	// anything that *cannot* have a parent to work
	default:
		if counter == ChallengeCoup ||
			counter == ChallengeIncome { //asd
			return false
		}
	}

	// cannot confront a coup, a confront, or an income, or a finanicial aid
	if counter == ChallengeConfront {
		switch action {
		case 0:
			fallthrough
		case ChallengeConfront:
			fallthrough
		case ChallengeIncome:
			fallthrough
		case ChallengeFinancialAid:
			fallthrough
		case ChallengeCoup:
			return false
		default:
			return true
		}
	}

	return false
}

func IsChallengeKindValid(kind uint8) bool {
	return kind >= ChallengeConfront && kind <= ChallengeCounterContessa
}

func NewChallenge(parent *Challenge, creator *Player, kind uint8, versus *Player, duration time.Duration) (*Challenge, error) {
	challenge := &Challenge{
		creator:   creator,
		kind:      kind,
		versus:    versus,
		timestamp: time.Now().Add(duration),
	}

	if err := challenge.IsValid(); err != nil {
		return nil, err
	}

	if parent != nil {
		if parent.IsValid() != nil {
			return nil, fmt.Errorf("parent: %w", parent.IsValid())
		}

		if parent.Passed() {
			return nil, ErrChallengePassed
		}

		if !IsCounterChallengeValid(parent.kind, kind) {
			return nil, ErrChallengeInvalidCounter
		}

		if parent.response != nil {
			return nil, ErrChallengeHasResponse
		}

		parent.response = challenge
	} else {
		if !IsCounterChallengeValid(0, kind) {
			return nil, ErrChallengeNeedsParent
		}
	}

	return challenge, nil
}
