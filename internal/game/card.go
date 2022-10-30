package game

const (
	CardEmpty uint8 = iota
	// CardAssassin is a card that kills opponent's card when the player
	// has 3 or more coins. This action is counterable by CardContessa.
	CardAssassin
	// CardDuke is a card that adds 3 coins to the player's bank.
	// CardDuke also has a counter action; which is to prevent any player
	// from taking financial aid.
	CardDuke
	// CardAmbassador is a card that takes 3 coins from the deck, looks
	// at them, and swaps them if wanted.
	// CardAmbassador also has a counter action; which is to prevent any
	// Captain from stealing from them.
	CardAmbassador
	// CardCaptain is a card that steals from any player except fellow
	// Captains and Ambassadors.
	// CardCaptain also has a counter action; which is to counter act
	// other Captains from stealing from themselves.
	CardCaptain
	// CardContessa is a card that has a counter of action; which is to
	// prevent any assassination attempts on itself by other Assassins.
	//
	// Do note: An assassination attempt is not the same as a coup.
	CardContessa
)

// IsValidCard returns true if the value is in between CardAssassin
// && CardContessa.
func IsValidCard(v uint8) bool {
	return v >= CardAssassin && v <= CardContessa
}
