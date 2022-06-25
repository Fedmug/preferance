package comb

import (
	"log"
)

const (
	byteLen   = 8
	byteCap   = 256
	totalRows = 9840
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

var rowsLittle []uint32
var rowToIndexLittle map[uint32]int

var rowsBig []uint32
var rowToIndexBig map[uint32]int

func initSqueezers() {
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

func getSqueezedCode(first, second, third uint8, endian Endian) uint32 {
	mask := ^(first + second + third)
	firstSqueezed := uint32(squeeze(mask, first, endian))
	secondSqueezed := uint32(squeeze(mask, second, endian))
	thirdSqueezed := uint32(squeeze(mask, third, endian))
	return firstSqueezed + secondSqueezed<<byteLen + thirdSqueezed<<(2*byteLen)
}

func initRows() {
	var code uint32
	rowsLittle = make([]uint32, 0, totalRows+1)
	rowsBig = make([]uint32, 0, totalRows+1)
	rowToIndexLittle = make(map[uint32]int, 1<<16)
	rowToIndexBig = make(map[uint32]int, 1<<16)
	rowsIndexLittle := make(map[uint32]int, totalRows+1)
	rowsIndexBig := make(map[uint32]int, totalRows+1)
	var indexLittle, indexBig int
	for ; code < 1<<24; code++ {
		first := uint8(code & 0xFF)
		second := uint8((code >> byteLen) & 0xFF)
		third := uint8((code >> (2 * byteLen)) & 0xFF)
		if first&second == 0 && second&third == 0 && first&third == 0 {
			codeLittle := getSqueezedCode(first, second, third, Little)
			if _, ok := rowsIndexLittle[codeLittle]; !ok {
				rowsLittle = append(rowsLittle, codeLittle)
				rowsIndexLittle[codeLittle] = indexLittle
				indexLittle++
			}
			rowToIndexLittle[code] = rowsIndexLittle[codeLittle]

			codeBig := getSqueezedCode(first, second, third, Big)
			if _, ok := rowsIndexBig[codeBig]; !ok {
				rowsBig = append(rowsBig, codeBig)
				rowsIndexBig[codeBig] = indexBig
				indexBig++
			}
			rowToIndexBig[code] = rowsIndexBig[codeBig]
		}
	}
}

func init() {
	initSqueezers()
	initRows()
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

func squeeze(mask, code uint8, endian Endian) uint8 {
	if endian == Little {
		return squeezerLittle[mask][code]
	}
	return squeezerBig[mask][code]
}

func unsqueeze(mask, code uint8, endian Endian) uint8 {
	if endian == Little {
		return unsqueezerLittle[mask][code]
	}
	return unsqueezerBig[mask][code]
}
