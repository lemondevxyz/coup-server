package game

import (
	"fmt"
	"strconv"
	"strings"
)

type Hand [2]uint8

func (h Hand) String() string { return fmt.Sprintf("%d:%d", h[0], h[1]) }
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

func (h Hand) IsEmpty() bool {
	return h[0] == CardEmpty && h[1] == CardEmpty
}

func (h Hand) IsEqual(v Hand) bool {
	return h[0] == v[0] && h[1] == v[1]
}
