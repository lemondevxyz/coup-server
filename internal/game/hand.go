package game

import (
	"fmt"
	"strconv"
	"strings"
)

// Hand is a slice that could contain two cards.
type Hand [2]uint8

// String returns a string version of hand.
//
// This function is particularily useful when paired with Hand.Unmarshal
// because it allows database to store this value and unmarshal it
// when wanted.
func (h Hand) String() string { return fmt.Sprintf("%d:%d", h[0], h[1]) }

// Unmarshal reads the values in a string formatted by String and sets
// the Hand's values to those read from the string.
func (h *Hand) Unmarshal(str string) error {
	split := strings.Split(str, ":")
	if len(split) < 2 {
		return fmt.Errorf("bad format")
	}

	fir, err := strconv.ParseUint(split[0], 10, 8)
	if err != nil {
		return err
	}

	sec, err := strconv.ParseUint(split[1], 10, 8)
	if err != nil {
		return err
	}

	h[0], h[1] = uint8(fir), uint8(sec)

	return nil
}

// IsEmpty returns true whenever Hand[0] and Hand[1] == CardEmpty
func (h Hand) IsEmpty() bool {
	return h[0] == CardEmpty && h[1] == CardEmpty
}

// IsEqual tests if the Hand is equal with another Hand.
func (h Hand) IsEqual(v Hand) bool {
	return h[0] == v[0] && h[1] == v[1]
}
