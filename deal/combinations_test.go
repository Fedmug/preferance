package deal

import (
	"fmt"
	"math/bits"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func Example_squeezeSecondHandCode() {
	var firstHand, secondHand SuitHandCode = 41, 80
	fmt.Printf("%08b (first hand)\n", firstHand)
	fmt.Printf("%08b (second hand)\n", secondHand)

	fmt.Printf("%08b (little)\n", squeeze(firstHand, secondHand, little))
	fmt.Printf("%08b (big)\n", squeeze(firstHand, secondHand, big))
	// Output:
	// 00101001 (first hand)
	// 01010000 (second hand)
	// 00001100 (little)
	// 01100000 (big)
}

func squeezeSecondHandCodeRandom(n int) {
	for i := 0; i < n; i++ {
		first := SuitHandCode(rand.Intn(byteCap))
		second := SuitHandCode(rand.Intn(byteCap))
		if first&second == 0 {
			fmt.Printf("%08b (first hand)\n", first)
			fmt.Printf("%08b (second hand)\n", second)

			fmt.Printf("%08b (little)\n", squeeze(first, second, little))
			fmt.Printf("%08b (big)\n\n", squeeze(first, second, big))
		}
	}
}

func Example_squeezeSecondHandCode_random() {
	rand.Seed(125)
	squeezeSecondHandCodeRandom(30)
	// Output:
	// 10010010 (first hand)
	// 01100101 (second hand)
	// 00011011 (little)
	// 11011000 (big)
	//
	// 11110111 (first hand)
	// 00000000 (second hand)
	// 00000000 (little)
	// 00000000 (big)
	//
	// 10000000 (first hand)
	// 01000011 (second hand)
	// 01000011 (little)
	// 10000110 (big)
	//
	// 01011010 (first hand)
	// 00100001 (second hand)
	// 00000101 (little)
	// 01010000 (big)
}

func Example_chase() {
	var n, s = 4, 2
	for i := 0; i < BinomialCoefficients[n][s]; i++ {
		fmt.Printf("%04b %d\n", chase.sequence[n][s][i], chase.sequenceToIndex[n][s][chase.sequence[n][s][i]])
	}
	n, s = 3, 0
	for i := 0; i < BinomialCoefficients[n][s]; i++ {
		fmt.Printf("%04b %d\n", chase.sequence[n][s][i], chase.sequenceToIndex[n][s][chase.sequence[n][s][i]])
	}
	n, s = 3, 3
	for i := 0; i < BinomialCoefficients[n][s]; i++ {
		fmt.Printf("%04b %d\n", chase.sequence[n][s][i], chase.sequenceToIndex[n][s][chase.sequence[n][s][i]])
	}
	// Output:
	// 1100 0
	// 1001 1
	// 1010 2
	// 0110 3
	// 0101 4
	// 0011 5
	// 0111 0
	// 0000 0
}

func TestSecondHandMap(t *testing.T) {
	total := 0
	for first := 0; first < byteCap; first++ {
		for second := 0; second < byteCap; second++ {
			if first&second == 0 {
				total++
				squeezedLittle := squeeze(SuitHandCode(first), SuitHandCode(second), little)
				squeezedBig := squeeze(SuitHandCode(first), SuitHandCode(second), big)
				unsqueezedLittle := unsqueeze(SuitHandCode(first), squeezedLittle, little)
				unsqueezedBig := unsqueeze(SuitHandCode(first), squeezedBig, big)
				require.EqualValues(t, bits.OnesCount8(uint8(second)), bits.OnesCount8(uint8(squeezedLittle)))
				require.EqualValues(t, bits.OnesCount8(uint8(second)), bits.OnesCount8(uint8(squeezedBig)))
				require.EqualValues(t, second, unsqueezedLittle)
				require.EqualValues(t, second, unsqueezedBig)
				onesCount := bits.OnesCount8(uint8(first))
				require.EqualValues(t, squeezedLittle, squeezedBig>>onesCount)
			}
		}
	}
	require.EqualValues(t, 6561, total)
}

func TestChase(t *testing.T) {
	for n := 0; n <= byteLen; n++ {
		for s := 0; s <= n; s++ {
			for i := 0; i < BinomialCoefficients[n][s]; i++ {
				code := chase.sequence[n][s][i]
				require.EqualValues(t, i, chase.sequenceToIndex[n][s][code])
				if i > 0 {
					require.EqualValues(t, 2, bits.OnesCount8(uint8(code^chase.sequence[n][s][i-1])))
				}
			}
		}
	}
}

func Example_suitToIndex() {
	handSizes := [NumberOfHands]int8{2, 1, 0}
	fmt.Println(suitToIndex(handSizes, [NumberOfHands]SuitHandCode{0x3, 0x4, 0}, little),
		suitToIndex(handSizes, [NumberOfHands]SuitHandCode{0x5, 0x2, 0}, little),
		suitToIndex(handSizes, [NumberOfHands]SuitHandCode{0x6, 0x1, 0}, little))

	handSizes = [NumberOfHands]int8{0, 1, 4}
	fmt.Println(suitToIndex(handSizes, [NumberOfHands]SuitHandCode{0, 0x1, 0x1e}, little),
		suitToIndex(handSizes, [NumberOfHands]SuitHandCode{0, 0x2, 0x1d}, little),
		suitToIndex(handSizes, [NumberOfHands]SuitHandCode{0, 0x4, 0x1b}, little),
		suitToIndex(handSizes, [NumberOfHands]SuitHandCode{0, 0x8, 0x17}, little),
		suitToIndex(handSizes, [NumberOfHands]SuitHandCode{0, 0x10, 0xf}, little))
	// Output:
	// 2 1 0
	// 2 3 1 4 0
}

func TestSuitIndex(t *testing.T) {
	var first, second, third int8
	for first = 0; first <= SuitLength; first++ {
		for second = 0; second <= SuitLength-first; second++ {
			for third = 0; third <= SuitLength-first-second; third++ {
				handSizes := [NumberOfHands]int8{first, second, third}
				maxIndex := BinomialCoefficients[first+second+third][first] *
					BinomialCoefficients[second+third][second]
				for index := 0; index < maxIndex; index++ {
					suitCodes := indexToSuitCodes(handSizes, index, little)
					for i := 0; i < NumberOfHands; i++ {
						require.EqualValues(t, handSizes[i], bits.OnesCount8(uint8(suitCodes[i])),
							fmt.Sprintf("hand size %v does not match code %v", handSizes[i], suitCodes[i]))
					}
					indexFromCodes := suitToIndex(handSizes, suitCodes, little)
					require.GreaterOrEqual(t, indexFromCodes, 0, fmt.Sprintf("index %d is negative", indexFromCodes))
					require.EqualValues(t, index, indexFromCodes,
						fmt.Sprintf("handSizes = %v, suitCodes = %v", handSizes, suitCodes))
				}
			}
		}
	}
}

func BenchmarkIndexToSuitCodes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		first := rand.Intn(SuitLength)
		second := rand.Intn(SuitLength - first)
		third := rand.Intn(SuitLength - first - second)
		index := rand.Intn(BinomialCoefficients[first+second+third][first] *
			BinomialCoefficients[second+third][second])
		indexToSuitCodes([NumberOfHands]int8{int8(first), int8(second), int8(third)}, index, little)
	}
}

func BenchmarkSuitToIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		handSizes := [NumberOfHands]int8{3, 1, 3}
		handCodes := [NumberOfHands]SuitHandCode{0x26, 0x40, 0x19}
		suitToIndex(handSizes, handCodes, little)
	}
}

func BenchmarkSqueezer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for first := 0; first < byteLen; first++ {
			for second := 0; second < byteLen; second++ {
				if first&second == 0 {
					code := squeeze(SuitHandCode(first), SuitHandCode(second), little)
					unsqueeze(SuitHandCode(first), code, little)
				}
			}
		}
	}
}

func BenchmarkChaseSeq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		first := rand.Intn(SuitLength)
		second := rand.Intn(SuitLength - first)
		third := rand.Intn(SuitLength - first - second)
		index := rand.Intn(BinomialCoefficients[first+second+third][first] *
			BinomialCoefficients[second+third][second])
		firstIndex := index / BinomialCoefficients[second+third][second]
		_ = chase.sequence[first+second+third][second+third][firstIndex]
	}
}
