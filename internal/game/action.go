package game

const (
	ActionIncome uint8 = iota + 1
	ActionFinancialAid
	ActionCoup
	ActionClaim
)

func IncomeAction(coins uint8) uint8 {
	if coins == 10 {
		return coins
	}

	return coins + 1
}

func FinancialAidAction(coins uint8) uint8 {
	if coins == 10 {
		return coins
	}

	return coins + 2
}

func CoupAction(coins uint8, place uint8, hand Hand) (uint8, Hand) {
	if coins < 7 || place > 1 || hand.IsEmpty() {
		return coins, hand
	}
	
	coins = coins - 7
	newHand := Hand{hand[0], hand[1]}
	newHand[place] = CardEmpty

	return coins, newHand
}

func AssassinAction(coins uint8, place uint8, hand Hand) (uint8, Hand) {
	if coins < 3 || place > 1 || hand.IsEmpty() {
		return coins, hand
	}

	coins = coins - 3
	newHand := Hand{hand[0], hand[1]}
	newHand[place] = CardEmpty

	return coins, newHand
}

func DukeAction(coins uint8) uint8 {
	if coins == 10 {
		return coins
	}

	return coins + 3
}

func CaptainAction(captainCoins uint8, otherCoins uint8) (uint8, uint8) {
	modifiedCoins := int(otherCoins) - 2
	if modifiedCoins < 0 {
		modifiedCoins = 0
	}

	return captainCoins + (otherCoins - uint8(modifiedCoins)), uint8(modifiedCoins)
}

func AmbassadorAction(places [2]uint8, currentHand Hand, nextHand Hand) Hand {
	newHand := Hand{currentHand[0], currentHand[1]}
	// avoid duplicate card
	if places[0] == places[1] && places[0] <= 1 {
		return newHand
	}
	
	if places[0] < 2 {
		newHand[0] = nextHand[places[0]]
	}
	
	if places[1] < 2 {
		newHand[1] = nextHand[places[1]]
	}
	
	return newHand
}

type Action struct {
	author *Player
	action uint8
}