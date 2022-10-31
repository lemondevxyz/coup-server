package game

import (
	"testing"

	"github.com/matryer/is"
	//	"fmt"
)

func TestCoinsPlus(t *testing.T) {
	is := is.New(t)

	is.Equal(coinsPlus(0, 1), uint8(1))
	is.Equal(coinsPlus(0, 3), uint8(3))
	is.Equal(coinsPlus(0, 5), uint8(5))

	is.Equal(coinsPlus(10, 0), uint8(10))
	is.Equal(coinsPlus(11, 5), uint8(11))
}

func TestIncomeAction(t *testing.T) {
	is := is.New(t)
	is.Equal(IncomeAction(0), uint8(1))
	is.Equal(IncomeAction(10), uint8(10))
}

func TestFinancialAidAction(t *testing.T) {
	is := is.New(t)
	is.Equal(FinancialAidAction(0), uint8(2))
	is.Equal(FinancialAidAction(10), uint8(10))
}

func TestDukeAction(t *testing.T) {
	is := is.New(t)
	is.Equal(DukeAction(0), uint8(3))
	is.Equal(DukeAction(10), uint8(10))
}

func onlyHand(coins uint8, hand Hand) Hand   { return hand }
func onlyCoins(coins uint8, hand Hand) uint8 { return coins }

func TestRemoveAt(t *testing.T) {
	is := is.New(t)

	is.Equal(removeFromHand(0, Hand{CardContessa, CardAmbassador}), Hand{CardEmpty, CardAmbassador})
	is.Equal(removeFromHand(1, Hand{CardAmbassador, CardContessa}), Hand{CardAmbassador, CardEmpty})
	is.Equal(removeFromHand(0, Hand{CardEmpty, CardAmbassador}), Hand{CardEmpty, CardEmpty})
	is.Equal(removeFromHand(1, Hand{CardAmbassador, CardEmpty}), Hand{CardEmpty, CardEmpty})
}

func TestMinusAndRemove(t *testing.T) {
	isEmpty := func(coins *uint8, hand *Hand) func(uint8, Hand) bool {
		return func(newCoins uint8, newHand Hand) bool {
			return *coins == newCoins && hand.IsEqual(newHand)
		}
	}

	coinsVal := uint8(5)
	coins, place, hand := &coinsVal, uint8(0), Hand{CardDuke, CardContessa}

	is := is.New(t)
	isEmptyPt2 := isEmpty(coins, &hand)

	*coins = 0

	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			*coins = uint8(0)
		case 1:
			*coins = uint8(5)
			place = uint8(2)
		case 2:
			place = uint8(0)
			hand[0], hand[1] = CardEmpty, CardEmpty
		}
		is.True(isEmptyPt2(minusCoinsRemoveFromHand(5, *coins, place, hand)))
	}

	hand = Hand{CardContessa, CardAmbassador}
	is.True(!isEmptyPt2(minusCoinsRemoveFromHand(5, *coins, place, hand)))
	is.Equal(onlyCoins(minusCoinsRemoveFromHand(5, *coins, place, hand)), uint8(0))
	is.Equal(onlyHand(minusCoinsRemoveFromHand(5, *coins, place, hand)), removeFromHand(place, hand))
}

func TestCoupAction(t *testing.T) {
	is := is.New(t)
	coins, place, hand := uint8(7), uint8(0), Hand{CardContessa, CardAmbassador}

	haveCoins, haveHand := CoupAction(coins, place, hand)
	wantCoins, wantHand := minusCoinsRemoveFromHand(coins, coins, place, hand)
	is.Equal(wantCoins, haveCoins)
	is.Equal(wantHand, haveHand)
}

func TestAssassinAction(t *testing.T) {
	is := is.New(t)
	coins, place, hand := uint8(3), uint8(0), Hand{CardContessa, CardAmbassador}

	haveCoins, haveHand := AssassinAction(coins, place, hand)
	wantCoins, wantHand := minusCoinsRemoveFromHand(coins, coins, place, hand)
	is.Equal(wantCoins, haveCoins)
	is.Equal(wantHand, haveHand)
}

func TestClaimPunishmentAction(t *testing.T) {
	is := is.New(t)
	place, hand := uint8(0), Hand{CardContessa, CardAmbassador}

	haveHand := ClaimPunishmentAction(place, hand)
	_, wantHand := minusCoinsRemoveFromHand(0, 0, place, hand)
	is.Equal(wantHand, haveHand)
}

func TestCaptainAction(t *testing.T) {
	is := is.New(t)
	coins, other := CaptainAction(2, 0)
	is.True(coins == 2)
	is.True(other == 0)

	coins, other = CaptainAction(2, 1)
	is.True(coins == 3)
	is.True(other == 0)

	coins, other = CaptainAction(2, 2)
	is.True(coins == 4)
	is.True(other == 0)

	coins, other = CaptainAction(4, 3)
	is.True(coins == 6)
	is.True(other == 1)
}

func TestAmbassadorAction(t *testing.T) {
	currentHand := Hand{CardAmbassador, CardContessa}
	nextHand := Hand{CardDuke, CardAssassin}

	is := is.New(t)

	wantDeck := Hand{nextHand[0], nextHand[1]}
	wantHand := Hand{currentHand[0], currentHand[1]}
	isEqual := func(one, two uint8) {
		haveHand, haveDeck := AmbassadorAction([2]uint8{one, two}, currentHand, nextHand)

		is.Equal(wantHand, haveHand)
		is.Equal(wantDeck, haveDeck)
	}

	// 2, 1
	isEqual(2, 2)
	// must not mutate
	is.Equal(currentHand, Hand{CardAmbassador, CardContessa})

	// 0, 1
	wantDeck = Hand{currentHand[0], currentHand[1]}
	wantHand = Hand{nextHand[0], nextHand[1]}
	//is.Equal(want, AmbassadorAction([2]uint8{0, 1}, currentHand, nextHand))
	isEqual(0, 1)

	// 0, 2
	wantDeck = Hand{currentHand[0], nextHand[1]}
	wantHand = Hand{nextHand[0], currentHand[1]}
	isEqual(0, 2)

	// 1, 2
	wantDeck = Hand{nextHand[0], currentHand[0]}
	wantHand = Hand{nextHand[1], currentHand[1]}
	isEqual(1, 2)

	// 1, 1
	wantDeck = nextHand
	wantHand = currentHand
	isEqual(1, 1)

	// 1, 0
	wantDeck = Hand{currentHand[1], currentHand[0]}
	wantHand = Hand{nextHand[1], nextHand[0]}
	isEqual(1, 0)

	// 2, 1
	wantDeck = Hand{nextHand[0], currentHand[1]}
	wantHand = Hand{currentHand[0], nextHand[1]}
	isEqual(2, 1)

	// 2, 0
	wantDeck = Hand{currentHand[1], nextHand[1]}
	wantHand = Hand{currentHand[0], nextHand[0]}
	isEqual(2, 0)
}

func TestActionDo(t *testing.T) {
	player := &Player{}

	a := &Action{Kind: ActionIncome, author: player}
	a.do()

	is := is.New(t)
	is.Equal(player.Coins, uint8(1))

	a.Kind = ActionCoup

	target := &Player{Hand: Hand{0: CardAssassin}}

	player.Coins = 7

	a.Kind = ActionCoup
	a.against = target

	val := uint8(0)

	coins, hand := CoupAction(player.Coins, val, target.Hand)
	a.AssassinPlace = &val
	a.do()

	is.Equal(target.Hand, hand)
	is.Equal(player.Coins, coins)

	a.Kind = ActionFinancialAid
	a.do()
	is.Equal(player.Coins, uint8(2))

	// Character specific tests
	a.Kind = ActionCharacter
	a.author = player
	a.against = target
	// Duke
	{
		a.Character = CardDuke
		coins := uint8(0)

		player.Coins = coins
		a.do()
		is.Equal(player.Coins, DukeAction(coins))
	}
	// Assassin
	{
		target.Hand[0] = CardAssassin
		a.Character = CardAssassin
		player.Coins = 3
		coins, hand := AssassinAction(player.Coins, 0, target.Hand)

		val := uint8(0)

		a.AssassinPlace = &val
		a.do()

		is.Equal(target.Hand, hand)
		is.Equal(player.Coins, coins)
	}
	// Captain
	{
		player.Coins = 0
		target.Coins = 2
		pCoins, tCoins := CaptainAction(player.Coins, target.Coins)

		a.Character = CardCaptain
		a.do()

		is.Equal(pCoins, player.Coins)
		is.Equal(tCoins, target.Coins)
	}
	// Ambassador
	{
		hand := Hand{0: CardContessa, 1: CardDuke}
		places := Hand{0: 0, 1: 1}

		newHand, newDeck := AmbassadorAction(places, player.Hand, hand)
		a.Character = CardAmbassador
		a.AmbassadorPlace = places
		a.AmbassadorHand = hand
		a.do()

		is.Equal(player.Hand, newHand)
		is.Equal(a.AmbassadorHand, newDeck)
	}
	// ClaimPunishment
	{
		coins := uint8(2)
		player.Coins = coins
		a.Character, target.Hand[0] = CardAssassin, CardAssassin
		hand := ClaimPunishmentAction(0, a.against.Hand)

		val := uint8(0)

		a.Kind = ActionClaimPunishment
		a.AssassinPlace = &val
		a.do()

		is.Equal(target.Hand, hand)
		is.Equal(player.Coins, coins)
	}
}

func TestActionValid(t *testing.T) {
	a := &Action{}

	is := is.New(t)
	is.Equal(a.IsValid(), ErrInvalidActionAuthor)

	a.author = &Player{}
	is.Equal(a.IsValid(), ErrInvalidActionAuthor)

	a.author = &Player{Hand: Hand{0: CardAmbassador}}

	val := uint8(1)
	a.AgainstID = &val
	is.Equal(a.IsValid(), ErrInvalidActionAgainst)

	a.against = &Player{}
	is.Equal(a.IsValid(), ErrInvalidActionAgainst)

	a.against = &Player{Hand: Hand{0: CardAmbassador}}

	val = 2
	a.AssassinPlace = &val

	is.Equal(a.IsValid(), ErrInvalidActionPlace)
	val = 0
	a.AssassinPlace = &val

	is.Equal(a.IsValid(), ErrInvalidActionKind)

	a.Kind = ActionCharacter + 1
	is.Equal(a.IsValid(), ErrInvalidActionKind)

	a.Kind = ActionCharacter
	is.NoErr(a.IsValid())

	a.Kind = ActionClaimPunishment
	is.NoErr(a.IsValid())
}

func TestActionSetPlayers(t *testing.T) {
	a := &Action{}
	a.AuthorID = 255
	pl1, pl2 := &Player{}, &Player{}

	is := is.New(t)
	is.Equal(a.setPlayer([]*Player{pl1}), ErrInvalidActionAuthor)

	a.AuthorID = 0
	is.Equal(a.setPlayer([]*Player{nil}), ErrInvalidActionAuthor)
	is.NoErr(a.setPlayer([]*Player{pl1}))
	is.Equal(a.author, pl1)

	val := uint8(255)
	a.AgainstID = &val

	is.Equal(a.setPlayer([]*Player{pl1, pl2}), ErrInvalidActionAgainst)
	val = 1
	a.AgainstID = &val

	is.Equal(a.setPlayer([]*Player{pl1, nil}), ErrInvalidActionAgainst)
	is.Equal(a.setPlayer([]*Player{pl1, pl1}), ErrInvalidActionSamePlayer)
	is.NoErr(a.setPlayer([]*Player{pl1, pl2}))

	is.Equal(a.against, pl2)
}

func TestActionValidClaim(t *testing.T) {
	is := is.New(t)

	a := Action{}

	is.Equal(a.validClaim(nil), ErrInvalidClaim)
	is.Equal(a.validClaim(&claim{}), ErrInvalidClaimHasNotFinished)
	is.Equal(a.validClaim(&claim{challenge: new(bool)}), ErrInvalidActionFrozen)
	is.Equal(a.validClaim(&claim{succeed: new(bool), character: 1}), ErrInvalidCharacter)

	is.NoErr(a.validClaim(&claim{succeed: new(bool), character: 0}))
}

func TestIsValidCounterAction(t *testing.T) {
	is := is.New(t)

	is.True(IsValidCounterAction(Action{Kind: ActionFinancialAid}, Action{
		Kind:      ActionCharacter,
		Character: CardDuke,
	}))

	is.Equal(IsValidCounterAction(Action{
		Kind:      ActionCharacter,
		Character: CardAssassin,
	}, Action{
		Kind:      ActionCharacter,
		Character: CardContessa,
	}), IsValidCounterClaim(CardAssassin, CardContessa))

	is.True(!IsValidCounterAction(Action{Kind: ActionIncome}, Action{Kind: ActionIncome}))
}
