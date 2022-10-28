package game

type Hand [2]uint8

func (h Hand) IsEmpty() bool {
	return h[0] == CardEmpty && h[1] == CardEmpty
}

func (h Hand) Equal(v Hand) bool {
	return h[0] == v[0] && h[1] == v[1]
}