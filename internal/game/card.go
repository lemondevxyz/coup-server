package game

const (
	CardEmpty uint8 = iota
	CardAssassin
	CardDuke
	CardAmbassador
	CardCaptain
	CardContessa
)

func IsValidCard(v uint8) bool {
	return v >= CardAssassin && v <= CardContessa
}