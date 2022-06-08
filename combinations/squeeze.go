package combinations

import "log"

const (
	byteLen = 8
	byteCap = 256
)

type Endian int8

const (
	// squeeze to lowest card, i.e., 7
	Little Endian = iota

	// squeeze to highest card, i.e., A
	Big
)

// Given mask and code, retrieve squeezed code, e. g.
// 00110100, 01000010 -> 00001010 if endian = little
// 00110100, 01000010 -> 01010000 if endian = big
var squeezerLittle [byteCap][byteCap]uint8
var squeezerBig [byteCap][byteCap]uint8

// Given mask and code, retrieve unsqueezed code, e. g.
// 00010100, 00001101 -> 00101001 if endian = little
// 00010100, 00110100 -> 00101001 if endian = big
var unsqueezerLittle [byteCap][byteCap]uint8
var unsqueezerBig [byteCap][byteCap]uint8

func init() {
	for first := 0; first < byteCap; first++ {
		for second := 0; second < byteCap; second++ {
			if first&second == 0 {
				squeezedCode := squeezeSecondHandCode(uint8(first), uint8(second), Little)
				squeezerLittle[first][second] = squeezedCode
				unsqueezerLittle[first][squeezedCode] = uint8(second)

				squeezedCode = squeezeSecondHandCode(uint8(first), uint8(second), Big)
				squeezerBig[first][second] = squeezedCode
				unsqueezerBig[first][squeezedCode] = uint8(second)
			}
		}
	}
}

func squeezeSecondHandCode(firstHandCode, secondHandCode uint8, endian Endian) uint8 {
	var i, result, shift uint8
	for i = 0; i < byteLen; i++ {
		index := i
		if endian == Big {
			index = byteLen - 1 - i
		}
		if firstHandCode&(1<<index) == 0 {
			if endian == Little {
				result |= (secondHandCode & (1 << index)) >> shift
			} else if endian == Big {
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

func Squeeze(mask, code uint8, endian Endian) uint8 {
	if endian == Little {
		return squeezerLittle[mask][code]
	}
	return squeezerBig[mask][code]
}

func Unsqueeze(mask, code uint8, endian Endian) uint8 {
	if endian == Little {
		return unsqueezerLittle[mask][code]
	}
	return unsqueezerBig[mask][code]
}
