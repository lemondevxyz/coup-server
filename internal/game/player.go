package game

// Player is a structure representing a player. A player consists of
// two things: a Hand(collection of cards) and a Coin balance.
type Player struct {
	dead  bool
	Coins uint8
	Hand  Hand
}

// IsDead returns true if the player has an empty hand.
//
// Do note: Once IsDead() returns true, it will return true for all
//          subsequent calls despite the Hand's actual value.
func (p *Player) IsDead() bool {
	if p.dead {
		return true
	}

	if p.Hand.IsEmpty() {
		p.dead = true
		return true
	}

	return false
}
