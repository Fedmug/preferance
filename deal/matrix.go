package deal

import (
	"fmt"
	"math/bits"
	"math/rand"
	"strings"

	"golang.org/x/exp/constraints"
)

const FullSuit SuitHandCode = 0xFF

type SuitHandCode uint8

type DealMatrix [NumberOfSuits][NumberOfHands]SuitHandCode

type DensePolicy int8

const (
	All DensePolicy = iota
	Max
	Min
	Random
)

const NullDelimiter rune = '\u0000'

// Parses string with hand cards, e.g. ♠KQJ9♣J♦Q♥AJ107 (delimiter is null)
// or A.KJT.87.QJ97 (delimiter = '.')
func parseHandString(s string, delimiter rune) ([NumberOfSuits]SuitHandCode, error) {
	var result [NumberOfSuits]SuitHandCode
	replacedString := strings.ReplaceAll(s, "10", "T")
	suitIndex := 0
	for _, ch := range replacedString {
		if ch == delimiter {
			suitIndex += 1
			continue
		}
		if suit, ok := SuitSymbolsMap[ch]; ok {
			suitIndex = int(suit)
			continue
		}
		if rank_code, ok := RankSymbolsMap[ch]; !ok {
			return result, fmt.Errorf("unknown rank symbol: %c", ch)
		} else {
			result[suitIndex] |= 1 << rank_code
		}
	}
	return result, nil
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func DealMatrixFromStrings(handStrings []string, delimiter rune) (DealMatrix, error) {
	var matrix DealMatrix
	if len(handStrings) != NumberOfHands {
		return matrix, fmt.Errorf("unexpected number of hands: want %d, got %d",
			NumberOfHands, len(handStrings))
	}
	for i, s := range handStrings {
		handCards, err := parseHandString(s, delimiter)
		if err != nil {
			return matrix, err
		}
		for j, code := range handCards {
			matrix[j][i] = code
		}
	}
	/*
		for i := 0; i < NumberOfSuits; i++ {
			matrix.missedCards[i] = FullSuit
			for j := 0; j < NumberOfHands; j++ {
				matrix.missedCards[i] -= matrix.codes[i][j]
			}
		}*/
	// matrix.squeeze(InvalidSuit)
	return matrix, nil
}

func (dm *DealMatrix) squeeze(endian Endian) DealMatrix {
	var result DealMatrix
	var missedCards [NumberOfSuits]SuitHandCode
	for i := 0; i < NumberOfSuits; i++ {
		for j := 0; j < NumberOfHands; j++ {
			missedCards[i] += dm[i][j]
		}
		missedCards[i] = ^missedCards[i]
	}
	for i := 0; i < NumberOfSuits; i++ {
		for j := 0; j < NumberOfHands; j++ {
			result[i][j] = squeeze(missedCards[i], dm[i][j], endian)
		}
	}
	return result
}

/*
// converts code = x0+y to xy, y = prefix
func squeezeZeros(code, prefixLen, zerosLen uint8) uint8 {
	if prefixLen == 8 {
		return code
	}
	prefix := code % (1 << prefixLen)
	suffix := (code >> (prefixLen + zerosLen)) << prefixLen
	return prefix + suffix
}

func (dm *DealMatrix) squeeze() {
	for i := 0; i < NumberOfSuits; i++ {
		for j := 0; j < NumberOfHands; j++ {
			dm.missedCards[i] |= ^dm.codes[i][j]
			dm.squeezedCodes[i][j] = dm.codes[i][j]
		}
		directCodeSums := ^dm.missedCards[i]
		invertedCodeSums := dm.missedCards[i]
		var zerosLen, prefixLen uint8
		for directCodeSums > 0 {
			if directCodeSums&1 == 0 {
				zerosLen = uint8(bits.TrailingZeros8(uint8(directCodeSums)))
				for j := 0; j < NumberOfHands; j++ {
					dm.squeezedCodes[i][j] =
						SuitHandCode(squeezeZeros(uint8(dm.squeezedCodes[i][j]), prefixLen, zerosLen))
				}
				directCodeSums >>= zerosLen
				invertedCodeSums >>= zerosLen
			} else {
				onesLen := uint8(bits.TrailingZeros8(uint8(invertedCodeSums)))
				directCodeSums >>= onesLen
				invertedCodeSums >>= onesLen
				prefixLen += onesLen
			}
		}
	}
}
*/

func (dm *DealMatrix) SuitSizes() [NumberOfSuits]int8 {
	var result [NumberOfSuits]int8
	for i := 0; i < NumberOfSuits; i++ {
		for j := 0; j < NumberOfHands; j++ {
			result[i] += int8(bits.OnesCount8(uint8(dm[i][j])))
		}
	}
	return result
}

func (dm *DealMatrix) HandSizes() [NumberOfHands]int8 {
	var result [NumberOfHands]int8
	for i := 0; i < NumberOfHands; i++ {
		for j := 0; j < NumberOfSuits; j++ {
			result[i] += int8(bits.OnesCount8(uint8(dm[j][i])))
		}
	}
	return result
}

func (dm *DealMatrix) ContingencyTable() ContingencyTable {
	var result ContingencyTable
	for i := 0; i < NumberOfSuits; i++ {
		for j := 0; j < NumberOfHands; j++ {
			result[i][j] += int8(bits.OnesCount8(uint8(dm[i][j])))
		}
	}
	return result
}

func (dm *DealMatrix) DeckSize() int8 {
	var result int8
	for _, suitSize := range dm.SuitSizes() {
		result += suitSize
	}
	return result
}

func (dm *DealMatrix) String() string {
	contingencyTable := dm.ContingencyTable()
	var result strings.Builder
	for j := 0; j < NumberOfHands; j++ {
		needSpace := false
		for i := 0; i < NumberOfSuits; i++ {
			if contingencyTable[i][j] > 0 {
				if needSpace {
					result.WriteString(" ")
				}
				needSpace = true
				result.WriteRune(SuitSymbols[i])
				for k := SuitLength - 1; k >= 0; k-- {
					if dm[i][j]&(1<<k) != 0 {
						result.WriteRune(RankSymbols[k])
					}
				}
			}
		}
		if j+1 < NumberOfHands {
			result.WriteRune('\n')
		}
	}
	return result.String()
}

func (dm *DealMatrix) Add(card Card, hand HandIndex) {
	suit := card.Suit()
	rank := card.Rank()
	dm[suit][hand] += 1 << rank
}

func (dm *DealMatrix) Remove(card Card, hand HandIndex) {
	suit := card.Suit()
	rank := card.Rank()
	dm[suit][hand] -= 1 << rank
}

func (dm *DealMatrix) GetMoves(suit Suit, hand HandIndex, policy DensePolicy, isTrump bool) []Card {
	directCode := uint8(dm[suit][hand])
	invertedCode := ^directCode
	moves := make([]Card, 0, bits.OnesCount8(directCode))
	offset := 0
	for directCode > 0 {
		shift := 0
		if directCode&1 == 0 {
			shift = min(bits.TrailingZeros8(directCode), SuitLength-offset)
		} else {
			shift = min(bits.TrailingZeros8(invertedCode), SuitLength-offset)
			if shift == 1 {
				moves = append(moves, NewCard(suit, Rank(offset), isTrump))
			} else {
				switch policy {
				case Min:
					moves = append(moves, NewCard(suit, Rank(offset), isTrump))
				case Max:
					moves = append(moves, NewCard(suit, Rank(offset+shift-1), isTrump))
				case Random:
					moves = append(moves, NewCard(suit, Rank(offset+rand.Intn(shift)), isTrump))
				case All:
					for idx := offset; idx < offset+shift; idx++ {
						moves = append(moves, NewCard(suit, Rank(idx), isTrump))
					}
				}
			}
		}
		offset += shift
		directCode >>= shift
		invertedCode >>= shift
	}
	return moves
}

func (dm *DealMatrix) Index(endian Endian) int64 {
	var result int64
	var coef int64 = 1
	squeezedMatrix := dm.squeeze(endian)
	contingencyTable := squeezedMatrix.ContingencyTable()
	for i := 0; i < NumberOfSuits; i++ {
		suitIndex := int64(suitToIndex(contingencyTable[i], squeezedMatrix[i], endian))
		result += suitIndex * coef
		coef *= int64(multinomial(contingencyTable[i]))
	}
	return result
}

func DealMatrixFromIndex(table ContingencyTable, index int64, endian Endian) DealMatrix {
	var result DealMatrix
	for i := 0; i < NumberOfSuits; i++ {
		coef := int64(multinomial(table[i]))
		result[i] = indexToSuitCodes(table[i], int(index%coef), endian)
		index /= coef
	}
	return result
}
