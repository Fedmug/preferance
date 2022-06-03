package deal

import (
	"fmt"
	"log"
	"math/bits"
)

const (
	byteLen = 8
	byteCap = 256
)

type Endian int8

const (
	little Endian = iota
	big
)

var BinomialCoefficients [byteLen + 1][byteLen + 1]int

func multinomial(handSizes [NumberOfHands]int8) int {
	n := handSizes[0] + handSizes[1] + handSizes[2]
	return BinomialCoefficients[n][handSizes[0]] * BinomialCoefficients[n-handSizes[0]][handSizes[1]]
}

type SuitHandCodeSqueezer struct {
	endian         Endian
	squeezeTable   [byteCap][byteCap]SuitHandCode
	unsqueezeTable [byteCap][byteCap]SuitHandCode
}

var squeezerLittle SuitHandCodeSqueezer
var squeezerBig SuitHandCodeSqueezer

func (sq *SuitHandCodeSqueezer) initTables(endian Endian) {
	sq.endian = endian
	for first := 0; first < byteCap; first++ {
		for second := 0; second < byteCap; second++ {
			if first&second == 0 {
				squeezedCode := squeezeSecondHandCode(SuitHandCode(first), SuitHandCode(second), endian)
				sq.squeezeTable[first][second] = squeezedCode
				sq.unsqueezeTable[first][squeezedCode] = SuitHandCode(second)
			}
		}
	}
}

// Given mask and code, retrieve squeezed code, e. g.
// 00110100, 01000010 -> 00001010 if endian = little
// 00110100, 01000010 -> 01010000 if endian = big
func (sq *SuitHandCodeSqueezer) squeeze(mask, code SuitHandCode) SuitHandCode {
	return sq.squeezeTable[mask][code]
}

// Given mask and code, retrieve unsqueezed code, e. g.
// 00010100, 00001101 -> 00101001 if endian = little
// 00010100, 00110100 -> 00101001 if endian = big
func (sq *SuitHandCodeSqueezer) unsqueeze(mask, code SuitHandCode) SuitHandCode {
	return sq.unsqueezeTable[mask][code]
}

type ChaseSequences struct {
	sequence        [byteLen + 1][byteLen + 1][]SuitHandCode
	sequenceToIndex [byteLen + 1][byteLen + 1][byteCap]int
}

var chase ChaseSequences

func squeezeSecondHandCode(firstHandCode, secondHandCode SuitHandCode, endian Endian) SuitHandCode {
	var i, result, shift SuitHandCode
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

func squeeze(mask, code SuitHandCode, endian Endian) SuitHandCode {
	if endian == little {
		// return squeezeSecondHandTableLittle[firstHandCode][secondHandCode]
		return squeezerLittle.squeeze(mask, code)
	}
	// return squeezeSecondHandTableBig[firstHandCode][secondHandCode]
	return squeezerBig.squeeze(mask, code)
}

func unsqueeze(mask, code SuitHandCode, endian Endian) SuitHandCode {
	if endian == little {
		// return unsqueezeSecondHandTableLittle[firstHandCode][secondHandCode]
		return squeezerLittle.unsqueeze(mask, code)
	}
	// return unsqueezeSecondHandTableBig[firstHandCode][secondHandCode]
	return squeezerBig.unsqueeze(mask, code)
}

func initBinomials() {
	for i := 0; i <= byteLen; i++ {
		BinomialCoefficients[i][0] = 1
		BinomialCoefficients[i][i] = 1
		for j := 1; j < i; j++ {
			BinomialCoefficients[i][j] = BinomialCoefficients[i-1][j] + BinomialCoefficients[i-1][j-1]
		}
	}
}

func initChase() {
	for i := 0; i <= byteLen; i++ {
		for j := 0; j <= byteLen; j++ {
			for k := 0; k < byteCap; k++ {
				chase.sequenceToIndex[i][j][k] = -1
			}
		}
	}
	for n := 0; n <= byteLen; n++ {
		for s := 0; s <= n; s++ {
			chase.sequence[n][s] = make([]SuitHandCode, 0, BinomialCoefficients[n][s])
			var code SuitHandCode
			for i := 0; i < n-s; i++ {
				code += 1 << (s + i)
			}
			w := make([]uint8, n+1)
			for i := 0; i <= n; i++ {
				w[i] = 1
			}
			r := s
			// begin, end := -1, -1
			if r == 0 {
				r = n - s
			}
			for i := 0; i < BinomialCoefficients[n][s]; i++ {
				chase.sequence[n][s] = append(chase.sequence[n][s], code)
				chase.sequenceToIndex[n][s][code] = i
				j := r
				for w[j] == 0 {
					w[j] = 1
					j++
				}
				if j == n {
					break
				}
				w[j] = 0
				if code&(1<<j) > 0 {
					if j%2 == 1 || code&(1<<(j-2)) > 0 {
						code -= 1 << j
						code += 1 << (j - 1)
						// begin, end = n-1-j, n-j
						if r == j && j > 1 {
							r = j - 1
						} else if r == j-1 {
							r = j
						}
					} else {
						code -= 1 << j
						code += 1 << (j - 2)
						// begin, end = n-1-j, n+1-j
						if r == j {
							r = 1
							if j-2 > 1 {
								r = j - 2
							}
						} else if r == j-2 {
							r = j - 1
						}
					}
				} else {
					if j%2 == 0 || code&(1<<(j-1)) > 0 {
						code += 1 << j
						code -= 1 << (j - 1)
						// begin, end = n-1-j, n-j
						if r == j && j > 1 {
							r = j - 1
						} else if r == j-1 {
							r = j
						}
					} else {
						code += 1 << j
						code -= 1 << (j - 2)
						// begin, end = n-1-j, n+1-j
						if r == j-2 {
							r = j
						} else if r == j-1 {
							r = j - 2
						}
					}
				}
			}
		}
	}
}

func init() {
	squeezerLittle.initTables(little)
	squeezerBig.initTables(big)
	initBinomials()
	initChase()
	initContingencyTablesPool()
}

func suitToIndex(handSizes [NumberOfHands]int8, handCodes [NumberOfHands]SuitHandCode, endian Endian) int {
	n := handSizes[0] + handSizes[1] + handSizes[2]
	squeezedSecondHandCode := squeeze(handCodes[0], handCodes[1], endian)
	firstCode := handCodes[0]
	if endian == big {
		firstCode = SuitHandCode(bits.Reverse8(uint8(firstCode)))
		squeezedSecondHandCode = SuitHandCode(bits.Reverse8(uint8(squeezedSecondHandCode)))
	}
	firstIndex := chase.sequenceToIndex[n][n-handSizes[0]][firstCode]
	secondIndex := chase.sequenceToIndex[n-handSizes[0]][handSizes[2]][squeezedSecondHandCode]
	if firstIndex < 0 || secondIndex < 0 {
		panic(fmt.Sprintf("at least one of indices %d and %d is negative\n"+
			"hand sizes = %v, hand codes = %08b\n", firstIndex, secondIndex, handSizes, handCodes))
	}
	return firstIndex*BinomialCoefficients[n-handSizes[0]][handSizes[1]] + secondIndex
}

func indexToSuitCodes(handSizes [NumberOfHands]int8, index int, endian Endian) [NumberOfHands]SuitHandCode {
	if index >= multinomial(handSizes) {
		panic(fmt.Sprintf("index %d is out of bound %d = multinomial(%v)",
			index, multinomial(handSizes), handSizes))
	}
	n := handSizes[0] + handSizes[1] + handSizes[2]
	firstIndex := index / BinomialCoefficients[n-handSizes[0]][handSizes[1]]
	secondIndex := index % BinomialCoefficients[n-handSizes[0]][handSizes[1]]
	var result [NumberOfHands]SuitHandCode
	result[0] = chase.sequence[n][n-handSizes[0]][firstIndex]
	if endian == big {
		result[0] = SuitHandCode(bits.Reverse8(uint8(result[0])))
	}
	secondCode := chase.sequence[n-handSizes[0]][handSizes[2]][secondIndex]
	if endian == big {
		secondCode = SuitHandCode(bits.Reverse8(uint8(secondCode)))
	}
	result[1] = unsqueeze(result[0], secondCode, endian)
	sum := (1 << n) - 1
	if endian == big {
		sum <<= 8 - n
	}
	result[2] = SuitHandCode(sum) - result[0] - result[1]
	return result
}
