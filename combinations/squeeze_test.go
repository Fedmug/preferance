package combinations

import (
	"fmt"
	"math/bits"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func Example_squeezeSecondHandCode() {
	var firstHand, secondHand uint8 = 41, 80
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
		first := uint8(rand.Intn(byteCap))
		second := uint8(rand.Intn(byteCap))
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

func TestSecondHandMap(t *testing.T) {
	total := 0
	for first := 0; first < byteCap; first++ {
		for second := 0; second < byteCap; second++ {
			if first&second == 0 {
				total++
				squeezedLittle := squeeze(uint8(first), uint8(second), little)
				squeezedBig := squeeze(uint8(first), uint8(second), big)
				unsqueezedLittle := unsqueeze(uint8(first), squeezedLittle, little)
				unsqueezedBig := unsqueeze(uint8(first), squeezedBig, big)
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
