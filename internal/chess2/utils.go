package chess2

import (
	"fmt"
	"math/bits"
)

// DumpMask takes a bitmask and returns a string representation of the selected
// bits, as 0 and 1 on 8 lines.
func DumpMask(mask uint64) string {
	// When rendering, square 0 is the top left
	mask = bits.Reverse64(mask)
	return fmt.Sprintf(
		"%08b\n%08b\n%08b\n%08b\n%08b\n%08b\n%08b\n%08b",
		(mask&0xff00000000000000)>>56,
		(mask&0x00ff000000000000)>>48,
		(mask&0x0000ff0000000000)>>40,
		(mask&0x000000ff00000000)>>32,
		(mask&0x00000000ff000000)>>24,
		(mask&0x0000000000ff0000)>>16,
		(mask&0x000000000000ff00)>>8,
		(mask&0x00000000000000ff)>>0,
	)
}

func eachBitSubset64(mask uint64, f func(subset uint64)) {
	subset := uint64(0)
	for {
		f(subset)
		subset = (subset - mask) & mask
		if subset == 0 {
			break
		}
	}
}
