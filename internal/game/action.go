package game

import "fmt"

const (
	// ActionIncome adds one to your income. See IncomeAction
	ActionIncome uint8 = iota + 1
	// ActionFinancialAid adds two to your income. See FinancialAidAction
	ActionFinancialAid
	// ActionCoup takes away 7 coins from you; lets you remove a card
	// from the opponent. See CoupAction
	ActionCoup
	// ActionCharacter executes the character's special action. See
	// DukeAction, CaptainAction, AmbassadorAction, AssassinAction.
	ActionCharacter
)

const (
	// ActionClaim, ActionClaimPassed, ActionClaimChallenge and
	// ActionClaimProof are all used for history. None of these have
	// any effect on Action.do
	ActionClaim uint8 = (^uint8(0)) - (iota + 1)
	ActionClaimPassed
	ActionClaimChallenge
	ActionClaimProof
	ActionClaimTakeCard
	ActionClaimPunishment
)

// coinsPlus is a function that adds an amount (plus) to the original
// amount of coins (coins). If coins equals 10 or more; then the function
// returns coins as is.
func coinsPlus(coins uint8, plus uint8) uint8 {
	if coins >= 10 {
		return coins
	}

	return coins + plus
}

// IncomeAction is a function that adds 1 to the amount of coins that
// was supplied to it.
//
// If the player has 10 or more coins, then IncomeAction returns the
// amount of coins it was supplied without modification.
func IncomeAction(coins uint8) uint8 { return coinsPlus(coins, 1) }

// FinancialAidAction is virtually the same as IncomeAction with two
// exceptions:
// - It gives the player 2 coins instead of 1
// - This action is stoppable by the Duke.
func FinancialAidAction(coins uint8) uint8 { return coinsPlus(coins, 2) }

// DukeAction is virtually the same as IncomeAction with two exceptions:
// - It gives the player 3 coins instead of 1
// - This action is stoppable by Challenging the player's claim.
func DukeAction(coins uint8) uint8 { return coinsPlus(coins, 3) }

// removeFromHand essentially sets hand[place] to EmptyCard. If
// hand[place] is already EmptyCard, set the other index to EmptyCard.
func removeFromHand(place uint8, hand Hand) Hand {
	newHand := Hand{hand[0], hand[1]}
	if newHand[place] == CardEmpty {
		switch place {
		case 0:
			newHand[1] = CardEmpty
		case 1:
			newHand[0] = CardEmpty
		}
	} else {
		newHand[place] = CardEmpty
	}

	return newHand
}

// miunsCoinsRemoveFromHand basically takes in a minimum amount of
// coins (minus); checks if coins equals that or more; if it does
// returns coins - minus and removeFromHand(place, hand).
//
// If it doesn't, returns the values as they were without modification.
func minusCoinsRemoveFromHand(minus, coins uint8, place uint8, hand Hand) (uint8, Hand) {
	if coins < minus || place > 1 || hand.IsEmpty() {
		return coins, hand
	}

	return coins - minus, removeFromHand(place, hand)
}

// CoupAction is a function that implements the Coup Action. Essentially,
// what it does is check if the player has 7 or more coins. If it does,
// it substracts 7, and removes a card from Hand at place.
//
// For example, if place was set to 0 and Hand was
// {CardContessa, CardAmbassador} it becomes
// {CardEmpty, CardAmbassador}.
//
// This function does not affect the data if one of the conditions apply:
// - Coins is less than 7
// - Place is more than 1
// - Hand is already empty
func CoupAction(coins uint8, place uint8, hand Hand) (uint8, Hand) {
	return minusCoinsRemoveFromHand(7, coins, place, hand)
}

// AssassinAction is the same as CoupAction but with three exceptions:
// - Instead of paying 7 coins, the player pays 3
// - The player must claim to have an Assassin card
// - This action is counterable by challenging the player's claim or
//   through claiming to have a Contessa card.
func AssassinAction(coins uint8, place uint8, hand Hand) (uint8, Hand) {
	return minusCoinsRemoveFromHand(3, coins, place, hand)
}

// CaptainAction is a function that steals coins from another player.
// If the target has 2 or more coins, the Captain steals 2. If they have
// 1, the Captain steals 1. If they have zero, the Captain doesn't steal
// anything.
func CaptainAction(captainCoins uint8, targetCoins uint8) (uint8, uint8) {
	diff := int(targetCoins) - 2
	if diff < 0 {
		diff = 0
	}

	return captainCoins + (targetCoins - uint8(diff)), uint8(diff)
}

// AmbassadorAction is a function that swaps or doesn't, depending on
// what the player chooses, cards from the deck to the player's hand.
//
// Essentially, the value of places controls which cards are swapped and
// which stay. If places[0] isn't 2 or more, then it takes places[0] from
// the deck and swaps it with the player's hand. The same goes for
// places[1].
//
// Do note: AmbassadorAction does not mutate the underlying hand and deck
//          values.
func AmbassadorAction(places [2]uint8, hand Hand, deck Hand) (copyHand Hand, copyDeck Hand) {
	copyHand = Hand{hand[0], hand[1]}
	copyDeck = Hand{deck[0], deck[1]}
	// avoid duplicate card
	if places[0] == places[1] && places[0] <= 1 {
		return
	}

	if places[0] < 2 {
		copyHand[0], copyDeck[places[0]] = copyDeck[places[0]], copyHand[0]
	}

	if places[1] < 2 {
		copyHand[1], copyDeck[places[1]] = copyDeck[places[1]], copyHand[1]
	}

	return
}

func ClaimPunishmentAction(place uint8, hand Hand) Hand {
	_, copyHand := minusCoinsRemoveFromHand(0, 0, place, hand)
	return copyHand
}

// Action is a structure of an action that's used to both manipulate
// the game programmatically, and to store actual Actions in databases.
//
// Action does not provide a specific method for saving, deleting, or
// even creating Databases for itself. Instead, it is left for outside
// packages to fit it with the Database.
//
// Action is essentially a mutable version of functions like CoupAction,
// IncomeAction and so on. It mutates anything important that relates
// to the Action's functionality.
//
// There are 4 types of actions:
// - ActionIncome
// - ActionFinancialAid
// - ActionCoup
// - ActionCharacter
//
// The reason this is seperated into 4 types instead of 3 + 5
// (5 = character amount) is because characters have to be claimed
// first. Only if that claim pass does the actual Action gets executed.
//
// Besides adding for claim functionality within Action would only
// complicate the package more than it has to be. This is a design choice
// that will never change.
type Action struct {
	AuthorID  uint8 `json:"author_id"`
	author    *Player
	Character uint8 `json:"character"`
	Kind      uint8 `json:"kind"`
	against   *Player
	AgainstID *uint8 `json:"against_id"`
	// Action specific fields; nullable
	// Used for assassin's action and coup
	AssassinPlace *uint8 `json:"assassin_place"`
	// Used for ambassador. Place denotes what to swap
	// and Hand denotes the two cards drawn from the deck.
	AmbassadorPlace Hand `json:"ambassador_place"`
	AmbassadorHand  Hand `json:"ambassador_hand"`
}

var (
	ErrInvalidActionAuthor     = fmt.Errorf("author: %w", ErrInvalidPlayer)
	ErrInvalidActionAgainst    = fmt.Errorf("against: %w", ErrInvalidPlayer)
	ErrInvalidActionSamePlayer = fmt.Errorf("author and against are the same player")
	ErrInvalidActionPlace      = fmt.Errorf("place must be [0, 1]")
	ErrInvalidActionKind       = fmt.Errorf("kind cannot be zero or bigger than ActionCharacter unless it is ActionClaimPunishment")
)

func (a Action) IsValid() error {
	if a.author == nil || a.author.IsDead() {
		return ErrInvalidActionAuthor
	}

	if a.AgainstID != nil {
		if a.against == nil || a.against.IsDead() {
			return ErrInvalidActionAgainst
		}
	}

	if a.AssassinPlace != nil {
		if *a.AssassinPlace > 1 {
			return ErrInvalidActionPlace
		}
	}

	if (a.Kind == 0 || a.Kind > ActionCharacter) && a.Kind != ActionClaimPunishment {
		return ErrInvalidActionKind
	}

	return nil
}

// do executes the underlying action if it matches; or does nothing silently.
// Essentially, it connects parameters & functionas together to mutate
// underlying player data.
//
// For example, say a player wanted to execute the IncomeAction, do would
// execute the function and mutate the player's coin amount to that of
// Income.
func (a *Action) do() {
	switch a.Kind {
	case ActionIncome:
		a.author.Coins = IncomeAction(a.author.Coins)
	case ActionCoup:
		a.author.Coins, a.against.Hand = CoupAction(a.author.Coins, *a.AssassinPlace, a.against.Hand)
	case ActionFinancialAid:
		a.author.Coins = FinancialAidAction(a.author.Coins)
	case ActionCharacter:
		switch a.Character {
		case CardAssassin:
			a.author.Coins, a.against.Hand = AssassinAction(a.author.Coins, *a.AssassinPlace, a.against.Hand)
		case CardDuke:
			a.author.Coins = DukeAction(a.author.Coins)
		case CardCaptain:
			a.author.Coins, a.against.Coins = CaptainAction(a.author.Coins, a.against.Coins)
		case CardAmbassador:
			a.author.Hand, a.AmbassadorHand = AmbassadorAction(a.AmbassadorPlace, a.author.Hand, a.AmbassadorHand)
		}
	case ActionClaimPunishment:
		a.against.Hand = ClaimPunishmentAction(*a.AssassinPlace, a.against.Hand)
	}
}

// setPlayer tries to find the players via AuthorID and AgainstID and sets
// action and against to the players it found respectively.
func (a *Action) setPlayer(players []*Player) error {
	index := a.AuthorID
	if int(index) >= len(players) || players[index] == nil {
		return ErrInvalidActionAuthor
	}
	a.author = players[index]

	if a.AgainstID != nil {
		index := *a.AgainstID
		if int(index) >= len(players) || players[index] == nil {
			return ErrInvalidActionAgainst
		}

		if a.author == players[index] {
			return ErrInvalidActionSamePlayer
		}

		a.against = players[index]
	}

	return nil
}

func (a Action) validClaim(c *claim) error {
	if c == nil {
		return ErrInvalidClaim
	} else if !c.IsFinished() {
		return ErrInvalidClaimHasNotFinished
	}

	succeed, challenge := c.Results()
	if succeed == nil && challenge != nil {
		return ErrInvalidActionFrozen
	} else if a.Character != c.character {
		return ErrInvalidCharacter
	}

	return nil
}

func IsValidCounterAction(a Action, b Action) bool {
	if b.Kind == ActionCharacter {
		return (a.Kind == ActionFinancialAid && b.Character == CardDuke) ||
			(a.Kind == ActionCharacter && IsValidCounterClaim(a.Character, b.Character))
	}

	return false
}
