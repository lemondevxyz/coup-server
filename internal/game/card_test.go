package game

import (
	"testing"
	"github.com/matryer/is"
)

func TestIsValidCard(t *testing.T) {
	is := is.New(t)
	for i := uint8(0); i < ^uint8(0); i++ {
		if i >= 1 && i <= 5 {
			is.True(IsValidCard(i))
		} else {
			is.True(!IsValidCard(i))
		}
	}
}