package game

import (
	"time"
	"testing"
)

// a bit ugly not gonna lie
var kinds = []uint8{ChallengeConfront, ChallengeIncome, ChallengeFinancialAid, ChallengeCoup, ChallengeActionAssassin, ChallengeActionDuke, ChallengeActionCaptain, ChallengeActionAmbassador, ChallengeCounterDuke, ChallengeCounterCaptain, ChallengeCounterAmbassador, ChallengeCounterContessa}

func TestIsChallengeKindValid(t *testing.T) {
	if IsChallengeKindValid(0) {
		t.Fatalf("0 is not a valid challenge kind")
	}//
	
	if IsChallengeKindValid(ChallengeCounterContessa+1) {
		t.Fatalf("%d is not a valid challenge kind", ChallengeCounterContessa+1)
	}
	
	for _, v := range kinds {
		if !IsChallengeKindValid(v) {
			t.Fatalf("%d is a valid challenge kind", v)
		}
	}
}

func TestChallengeIsValid(t *testing.T) {
	chl := &Challenge{}
	
	valid := func() {
		if chl.IsValid() == nil {
			t.Fatalf("challenge should be invalid: %v", chl)
		}
	}
	
	chl.creator = &Player{}
	valid()
	
	chl.kind = ChallengeIncome
	valid()
	
	chl.timestamp = time.Now()
	valid()
	
	chl.creator = &Player{false, 2, Hand{DukeCard, AmbassadorCard}}
	if chl.IsValid() != nil {
		t.Fatalf("challenge should be valid...")
	}
}

func TestChallengeSucceed(t *testing.T) {
	chl := &Challenge{}
	
	v := make(chan struct{})
	go func() {
		time.Sleep(time.Millisecond)
		b := true
		chl.succeed = &b
		v <- struct{}{}
	}()
	
	if !chl.Succeed() {
		t.Fatalf("challenge hasn't Succeeded")
	}
	<-v
	
	parent := &Challenge{}
	child := &Challenge{}
	parent.response = child
	v = make(chan struct{})
	go func() {
		time.Sleep(time.Millisecond)
		b := true
		child.succeed = &b
		v <- struct{}{}
	}()
	// asd
	if parent.Succeed() || !child.Succeed() {
		t.Fatalf("challenge hasn't Succeeded")
	}
	<-v	
	
}

func TestChallengeUnwrap(t *testing.T) {
	chl := &Challenge{}
	chl.response = &Challenge{}
	
	if chl.response != chl.Unwrap() {
		t.Fatalf("Unwrap() != Challenge.response")
	}
	
	chl.response = nil
	if chl.response != chl.Unwrap() {
		t.Fatalf("Unwrap() != Challenge.response")
	}
}

func TestChallengeRespondWith(t *testing.T) {
	chl := &Challenge{}
	
	newchl := &Challenge{}
	chl.RespondWith(newchl)
	
	if chl.response != newchl {
		t.Fatalf("chl.response != newchl: %v, %v", chl.response, newchl)
	}
}

func TestChallengeLast(t *testing.T) {
	child := []*Challenge{&Challenge{}, &Challenge{}, &Challenge{}, &Challenge{}}
	
	for k := range child {
		if k+1 < len(child) {
			child[k].response = child[k+1]
		}
		child[k].kind = uint8(k) + 1
	}
	
	if	child[0].Last() != child[3] {
		t.Fatalf("chl.Last() != child[3]: %v, %v", child[0].Last(), child[3])
	}
}

func TestChallengePassed(t *testing.T) {
	chl := &Challenge{timestamp: time.Now().Add(-1 * time.Microsecond)}
	
	if !chl.Passed() {
		t.Fatalf("challenge should've passed")
	}
	
	chl.timestamp = time.Now().Add(time.Microsecond)
	if chl.Passed() {
		t.Fatalf("challenge shouldn't have passed")
	}
	
	time.Sleep(time.Microsecond * 2)
	if !chl.Passed() {
		t.Fatalf("challenge should've passed")
	}
}

func TestIsCounterChallengeValid(t *testing.T) {
	val := [][2]uint8{{ChallengeFinancialAid, ChallengeCounterDuke}, {ChallengeActionAssassin, ChallengeCounterContessa}, {ChallengeActionCaptain, ChallengeCounterCaptain}, {ChallengeActionCaptain, ChallengeCounterCaptain}}
	
	for _, slice := range val {
		if !IsCounterChallengeValid(slice[0], slice[1]) {
			t.Fatalf("%d %d should be a valid counter challenge", slice[0], slice[1])
		}
	}
	
	for _, v := range kinds {
		if v == ChallengeCoup || v == ChallengeIncome || v == ChallengeFinancialAid || v == ChallengeConfront {
		// asd
			continue
		}
		
		if !IsCounterChallengeValid(v, ChallengeConfront) {
			t.Fatalf("ChallengeConfront should work with everything besides itself and ChallengeCoup")
		}
	}
	
	countered := []uint8{ChallengeCounterDuke, ChallengeCounterContessa, ChallengeCounterAmbassador, ChallengeCounterCaptain, ChallengeConfront}
	for _, v := range countered {
		if IsCounterChallengeValid(0, v) {
			t.Fatalf("%v is a counter action not an action", v)
		}
	}
	
	if IsCounterChallengeValid(ChallengeCoup, ChallengeCounterDuke) {
		t.Fatalf("Coup cannot have a counter action")
	}
	
	if IsCounterChallengeValid(ChallengeIncome, ChallengeCoup) {
		t.Fatalf("Coup is not a counter action")
	}
	
	for _, v := range kinds {
		if IsCounterChallengeValid(ChallengeIncome, v) {
			t.Fatalf("cannot counter act an income action")
		} else if IsCounterChallengeValid(v, ChallengeIncome) {
			t.Fatalf("challenge income is a counter action")
		}
	}
}

func TestNewChallenge(t *testing.T) {
	_, err := NewChallenge(nil, nil, 0, nil, 0)
	if err != ErrCreatorNil {
		t.Fatalf("Err should be ErrCreatorNil: %v, %v", err, ErrCreatorNil)
	}
	
	_, err = NewChallenge(nil, &Player{}, 0, nil, 0)
	if err != ErrCreatorDead {
		t.Fatalf("Err should be ErrCreatorDead: %v, %v", err, ErrCreatorDead)
	}
	
	player := &Player{Hand: Hand{DukeCard, AmbassadorCard}}
	_, err = NewChallenge(nil, player, 0, nil, 0)
	if err != ErrBadKind {
		t.Fatalf("Err should be ErrBadKind: %v, %v", err, ErrBadKind)
	}
	
	_, err = NewChallenge(nil, player, ChallengeCoup, nil, 1000)
	if err != nil {
		t.Fatalf("NewChallenge != nil: %v", err.Error())
	}
	
	_, err = NewChallenge(nil, player, ChallengeConfront, nil, 1000)
	if err != ErrChallengeNeedsParent {
		t.Fatalf("Err should be ErrChallengeNeedsParent: %v, %v", err, ErrChallengeNeedsParent)
	}
	
	chl, err := NewChallenge(nil, player, ChallengeIncome, nil, 1000*1000)
	t.Log(chl, err)
	
	_, err = NewChallenge(chl, player, ChallengeCoup, nil, 1000*1000)
	if err != ErrChallengeInvalidCounter {
		t.Fatalf("Err should be ErrChallengeForbiddenChild: %v, %v", err, ErrChallengeForbiddenChild)
	}
		
}