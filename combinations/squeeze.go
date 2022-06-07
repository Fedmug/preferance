package combinations

import "log"

const (
	byteLen = 8
	byteCap = 256
)

type Endian int8

const (
	// squeeze to lowest card, i.e., 7
	little Endian = iota

	// squeeze to highest card, i.e., A
	big
)

type SuitHandCodeSqueezer struct {
	endian         Endian
	squeezeTable   [byteCap][byteCap]uint8
	unsqueezeTable [byteCap][byteCap]uint8
}

var squeezerLittle SuitHandCodeSqueezer
var squeezerBig SuitHandCodeSqueezer

func (sq *SuitHandCodeSqueezer) initTables(endian Endian) {
	sq.endian = endian
	for first := 0; first < byteCap; first++ {
		for second := 0; second < byteCap; second++ {
			if first&second == 0 {
				squeezedCode := squeezeSecondHandCode(uint8(first), uint8(second), endian)
				sq.squeezeTable[first][second] = squeezedCode
				sq.unsqueezeTable[first][squeezedCode] = uint8(second)
			}
		}
	}
}

// Given mask and code, retrieve squeezed code, e. g.
// 00110100, 01000010 -> 00001010 if endian = little
// 00110100, 01000010 -> 01010000 if endian = big
func (sq *SuitHandCodeSqueezer) squeeze(mask, code uint8) uint8 {
	return sq.squeezeTable[mask][code]
}

// Given mask and code, retrieve unsqueezed code, e. g.
// 00010100, 00001101 -> 00101001 if endian = little
// 00010100, 00110100 -> 00101001 if endian = big
func (sq *SuitHandCodeSqueezer) unsqueeze(mask, code uint8) uint8 {
	return sq.unsqueezeTable[mask][code]
}

func squeezeSecondHandCode(firstHandCode, secondHandCode uint8, endian Endian) uint8 {
	var i, result, shift uint8
	for i = 0; i < byteLen; i++ {
		index := i
		if endian == big {
			index = byteLen - 1 - i
		}
		if firstHandCode&(1<<index) == 0 {
			if endian == little {
				result |= (secondHandCode & (1 << index)) >> shift
			} else if endian == big {
				result |= (secondHandCode & (1 << index)) << shift
			} else {
				log.Fatalf("unknown endian: %v", endian)
			}
		} else {
			shift++
		}
	}
	return result
}

func squeeze(mask, code uint8, endian Endian) uint8 {
	if endian == little {
		return squeezerLittle.squeeze(mask, code)
	}
	return squeezerBig.squeeze(mask, code)
}

func unsqueeze(mask, code uint8, endian Endian) uint8 {
	if endian == little {
		return squeezerLittle.unsqueeze(mask, code)
	}
	return squeezerBig.unsqueeze(mask, code)
}

func init() {
	squeezerLittle.initTables(little)
	squeezerBig.initTables(big)
}
