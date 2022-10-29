package game

import (
	"testing"

	"github.com/matryer/is"
	//	"fmt"
)

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

func isEmpty(coins *uint8, hand *Hand) func(uint8, Hand) bool {
	return func(newCoins uint8, newHand Hand) bool {
		return *coins == newCoins && hand.IsEqual(newHand)
	}
}

func onlyHand(coins uint8, hand Hand) Hand { return hand }

func TestCoupAction(t *testing.T) {
	coinsVal := uint8(3)
	coins, place, hand := &coinsVal, uint8(0), Hand{CardDuke, CardContessa}

	is := is.New(t)
	isEmpty := isEmpty(coins, &hand)

	*coins = 0

	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			*coins = uint8(0)
		case 1:
			*coins = uint8(7)
			place = uint8(2)
		case 2:
			place = uint8(0)
			hand[0], hand[1] = CardEmpty, CardEmpty
		}
		is.True(isEmpty(CoupAction(*coins, place, hand)))
	}
	//asd
	hand = Hand{CardDuke, CardContessa}
	is.True(!isEmpty(CoupAction(*coins, place, hand)))

	is.True(hand.IsEqual(Hand{CardDuke, CardContessa}))
	is.Equal(onlyHand(CoupAction(*coins, 0, hand))[0], CardEmpty)
	is.Equal(onlyHand(CoupAction(*coins, 1, hand))[1], CardEmpty)
}

func TestAssassinAction(t *testing.T) {
	coinsVal := uint8(3)
	coins, place, hand := &coinsVal, uint8(0), Hand{CardDuke, CardContessa}

	isEmpty := isEmpty(coins, &hand)

	is := is.New(t)
	*coins = 0
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			*coins = uint8(0)
		case 1:
			*coins = uint8(3)
			place = uint8(2)
		case 2:
			place = uint8(0)
			hand[0], hand[1] = CardEmpty, CardEmpty
		}
		is.True(isEmpty(AssassinAction(*coins, place, hand)))
	}

	hand = Hand{CardDuke, CardContessa}
	is.True(!isEmpty(AssassinAction(*coins, place, hand)))

	is.True(hand.IsEqual(Hand{CardDuke, CardContessa}))
	_, newHand := AssassinAction(*coins, place, hand)
	is.Equal(newHand[0], CardEmpty)

	_, newHand = AssassinAction(*coins, 1, hand)
	is.Equal(newHand[1], CardEmpty)
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

	want := Hand{currentHand[0], currentHand[1]}
	isEqual := func(one, two uint8) {
		is.Equal(want, AmbassadorAction([2]uint8{one, two}, currentHand, nextHand))
	}

	// 2, 1
	isEqual(2, 2)
	// must not mutate
	is.Equal(currentHand, Hand{CardAmbassador, CardContessa})

	// 0, 1
	want = Hand{nextHand[0], nextHand[1]}
	//is.Equal(want, AmbassadorAction([2]uint8{0, 1}, currentHand, nextHand))
	isEqual(0, 1)

	// 0, 2
	want = Hand{nextHand[0], currentHand[1]}
	isEqual(0, 2)

	// 1, 2
	want = Hand{nextHand[1], currentHand[1]}
	isEqual(1, 2)

	// 1, 1
	want = currentHand
	isEqual(1, 1)

	// 1, 0
	want = Hand{nextHand[1], nextHand[0]}
	isEqual(1, 0)

	// 2, 1
	want = Hand{currentHand[0], nextHand[1]}
	isEqual(2, 1)

	// 2, 0
	want = Hand{currentHand[0], nextHand[0]}
	isEqual(2, 0)
}

func TestActionDo(t *testing.T) {
	player := &Player{}

	a := &Action{Kind: ActionIncome, author: player}
	a.Do()

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
	a.Do()

	is.Equal(target.Hand, hand)
	is.Equal(player.Coins, coins)

	a.Kind = ActionFinancialAid
	a.Do()
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
		a.Do()
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
		a.Do()

		is.Equal(target.Hand, hand)
		is.Equal(player.Coins, coins)
	}
	// Captain
	{
		player.Coins = 0
		target.Coins = 2
		pCoins, tCoins := CaptainAction(player.Coins, target.Coins)

		a.Character = CardCaptain
		a.Do()

		is.Equal(pCoins, player.Coins)
		is.Equal(tCoins, target.Coins)
	}
	// Ambassador
	{
		hand := Hand{0: CardContessa, 1: CardDuke}
		places := Hand{0: 0, 1: 1}

		newHand := AmbassadorAction(places, player.Hand, hand)
		a.Character = CardAmbassador
		a.AmbassadorPlace = places
		a.AmbassadorHand = hand
		a.Do()

		is.Equal(player.Hand, newHand)
	}

}
