package game

type Player struct {
	dead bool
	Coins uint8
	Hand Hand
}

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