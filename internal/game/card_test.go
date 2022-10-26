package game

import (
	"testing"
)

func TestIncome(t *testing.T) {
	if Income(0) != 1 {
		t.Fatalf("Income doesn't give one coin")
	}

	if Income(10) == 11 {
		t.Fatalf("Income doesn't stop at 10")
	}
}

func TestFinancialAid(t *testing.T) {
	if FinancialAid(0) != 2 {
		t.Fatalf("FinancialAid doesn't give two coins")
	}

	if FinancialAid(10) != 10 {
		t.Fatalf("FinancialAid doesn't stop at 10")
	}
}

func TestDukeAction(t *testing.T) {
	if DukeAction(0) != 3 {
		t.Fatalf("DukeAction doesn't give three coins")
	}

	if DukeAction(10) != 10 {
		t.Fatalf("DukeAction doesn't stop at 10")
	}
}

func TestAssassinAction(t *testing.T) {
	coins, place, hand := uint8(3), uint8(0), Hand{DukeCard, ContessaCard}

	isEmpty := func() bool {
		newCoins, newHand := AssassinAction(coins, place, hand)

		if coins == newCoins && hand[0] == newHand[0] && hand[1] == newHand[1] {
			return true
		}

		return false
	}

	coins = 0
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			coins = uint8(0)
		case 1:
			coins = uint8(3)
			place = uint8(2)
		case 2:
			place = uint8(0)
			hand = Hand{EmptyCard, EmptyCard}
		}
		if !isEmpty() {
			t.Fatalf("bad parameters still yield the action; %d %d %v", coins, place, hand)
		}
	}

	hand = Hand{DukeCard, ContessaCard}
	if isEmpty() {
		t.Log("for reference: coins must be 3 or more, place needs to be more than 1, hand must not be empty")
		t.Fatalf("good parameters but no action; %d %d %v", coins, place, hand)
	}

	if !hand.Equal(Hand{DukeCard, ContessaCard}) {
		t.Fatalf("AssassinAction mutates the underlying hand")
	}
}

func TestCaptainAction(t *testing.T) {
	coins, other := CaptainAction(2, 0)
	if coins != 2 || other != 0 {
		t.Fatalf("CaptainAction doesn't stop stealing at 0 coins: %v, %v", coins, other)
	}

	coins, other = CaptainAction(2, 1)
	if coins != 3 || other != 0 {
		t.Fatalf("CaptainAction doesn't steal only one coin because the player doesn't have more: %v, %v", coins, other)
	}

	coins, other = CaptainAction(2, 2)
	if coins != 4 || other != 0 {
		t.Fatalf("CaptainAction doesn't steal 2 coins when the player has 2 coins: %v, %v", coins, other)
	}

	coins, other = CaptainAction(4, 3)
	if coins != 6 || other != 1 {
		t.Fatalf("CaptainAction doesn't steal 2 coins when the player has more than 2 coins: %v, %v", coins, other)
	}
}

func TestAmbassadorAction(t *testing.T) {
	currentHand := Hand{AmbassadorCard, ContessaCard}
	nextHand := Hand{DukeCard, AssassinCard}

	// 2, 1
	want := Hand{currentHand[0], nextHand[1]}
	if !want.Equal(AmbassadorAction([2]uint8{2, 1}, currentHand, nextHand)) {
		t.Fatalf("AmbassadorAction doesn't swap out the second card from the current hand to the second card from the next hand. %v, %v, %v", [2]uint8{2, 1}, currentHand, nextHand)
	}

	if !currentHand.Equal(Hand{AmbassadorCard, ContessaCard}) {
		t.Fatalf("AmbassadorAction modifies the underlying hand")
	}

	// 0, 1
	want = Hand{nextHand[0], nextHand[1]}

	if !want.Equal(AmbassadorAction([2]uint8{0, 1}, currentHand, nextHand)) {
		t.Fatalf("first and second card from the first hand should be swapped with their counterparts in the second hand: %v, %v, %v", [2]uint8{0, 1}, currentHand, nextHand)
	}

	// 0, 2
	want = Hand{nextHand[0], currentHand[1]}
	if !want.Equal(AmbassadorAction([2]uint8{0, 2}, currentHand, nextHand)) {
		t.Fatalf("first card from the first hand should be swapped with the first card from the second hand")
	}

	// 1, 2
	want = Hand{nextHand[1], currentHand[1]}
	if !want.Equal(AmbassadorAction([2]uint8{1, 2}, currentHand, nextHand)) {
		t.Fatalf("first card from the first hand should be swapped with the second card from the second hand")
	}

	// 1, 1
	want = currentHand
	if !want.Equal(AmbassadorAction([2]uint8{1, 1}, currentHand, nextHand)) {
		t.Fatalf("nothing should happen but it does; because swapping two cards with the same card induces duplicate cards")
	}

	// 2, 1
	want = Hand{currentHand[0], nextHand[1]}
	if !want.Equal(AmbassadorAction([2]uint8{2, 1}, currentHand, nextHand)) {
		t.Fatalf("second card from the first hand should be swapped with its countepart from the second hand")
	}

	// 2, 0
	want = Hand{currentHand[0], nextHand[0]}
	if !want.Equal(AmbassadorAction([2]uint8{2, 0}, currentHand, nextHand)) {
		t.Fatalf("second card from the first hand should be swapped with the first card in the second hand")
	}

}
