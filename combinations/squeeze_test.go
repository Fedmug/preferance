package comb

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

	fmt.Printf("%08b (little)\n", squeeze(firstHand, secondHand, Little))
	fmt.Printf("%08b (big)\n", squeeze(firstHand, secondHand, Big))
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

			fmt.Printf("%08b (little)\n", squeeze(first, second, Little))
			fmt.Printf("%08b (big)\n\n", squeeze(first, second, Big))
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

func Example_initRows() {
	fmt.Println("rowsLittle:", len(rowsLittle), len(rowToIndexLittle))
	fmt.Println("rowsBig:", len(rowsBig), len(rowToIndexBig))
	// Output:
	// rowsLittle: 9841 65536
	// rowsBig: 9841 65536
}

func Example_initRowsIndex() {
	var counters [totalRows + 1]int
	for _, v := range rowToIndexLittle {
		counters[v]++
	}
	fmt.Println(counters[:9])
	// Output:
	// [1 8 28 56 70 56 28 8 1]
}

func TestSecondHandMap(t *testing.T) {
	total := 0
	for first := 0; first < byteCap; first++ {
		for second := 0; second < byteCap; second++ {
			if first&second == 0 {
				total++
				squeezedLittle := squeeze(uint8(first), uint8(second), Little)
				squeezedBig := squeeze(uint8(first), uint8(second), Big)
				unsqueezedLittle := unsqueeze(uint8(first), squeezedLittle, Little)
				unsqueezedBig := unsqueeze(uint8(first), squeezedBig, Big)
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

func TestRowIndex(t *testing.T) {
	for i := 0; i < len(rowsLittle); i++ {
		require.EqualValues(t, i, rowToIndexLittle[rowsLittle[i]],
			fmt.Sprintf("little index is wrong: got %d, want %d", rowToIndexLittle[rowsLittle[i]], i))
		require.EqualValues(t, i, rowToIndexBig[rowsBig[i]],
			fmt.Sprintf("big index is wrong: got %d, want %d", rowToIndexBig[rowsBig[i]], i))
		require.LessOrEqual(t, bits.OnesCount32(rowsLittle[i]), 8,
			fmt.Sprintf("popcount in %x is too big: %d", rowsLittle[i], bits.OnesCount32(rowsLittle[i])))
		require.LessOrEqual(t, bits.OnesCount32(rowsBig[i]), 8,
			fmt.Sprintf("popcount in %x is too big: %d", rowsBig[i], bits.OnesCount32(rowsBig[i])))
	}
}

func BenchmarkSqueezer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for first := 0; first < byteCap; first++ {
			for second := 0; second < byteCap; second++ {
				if first&second == 0 {
					code := squeeze(uint8(first), uint8(second), Little)
					unsqueeze(uint8(first), code, Little)
				}
			}
		}
	}
}
