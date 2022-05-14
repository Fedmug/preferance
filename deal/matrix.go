package deal

import (
	"fmt"
	"math/bits"
	"math/rand"
	"strings"

	"golang.org/x/exp/constraints"
)

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
	return matrix, nil
}

func (dm *DealMatrix) SuitSizes() [NumberOfSuits]int8 {
	var result [NumberOfSuits]int8
	for i := range *dm {
		for j := range dm[i] {
			result[i] += int8(bits.OnesCount8(uint8(dm[i][j])))
		}
	}
	return result
}

func (dm *DealMatrix) ContingencyTable() [NumberOfSuits][NumberOfHands]int8 {
	var result [NumberOfSuits][NumberOfHands]int8
	for i := range *dm {
		for j := range dm[i] {
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
