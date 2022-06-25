package comb

import (
	"fmt"
	"math/bits"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Example_chase() {
	var n, s = 4, 2
	var i uint64
	for i = 0; i < Binomial(int8(n), int8(s)); i++ {
		fmt.Printf("%04b %d\n", chase.sequence[n][s][i], chase.sequenceToIndex[n][s][chase.sequence[n][s][i]])
	}
	n, s = 3, 0
	for i = 0; i < Binomial(int8(n), int8(s)); i++ {
		fmt.Printf("%04b %d\n", chase.sequence[n][s][i], chase.sequenceToIndex[n][s][chase.sequence[n][s][i]])
	}
	n, s = 3, 3
	for i = 0; i < Binomial(int8(n), int8(s)); i++ {
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

func Example_codes32() {
	var code uint32
	count := 0
	var countBySize [maxSuitSize]int
	var countMap [maxSuitSize]map[uint32]int
	for i := 0; i < maxSuitSize; i++ {
		countMap[i] = make(map[uint32]int, powersOfThree[i+1])
	}
	for ; code < 1<<24; code++ {
		first := uint8(code & 0xFF)
		second := uint8((code >> 8) & 0xFF)
		third := uint8((code >> 16) & 0xFF)
		n := bits.OnesCount32(code)
		// popCounts := composition{
		// 	int8(bits.OnesCount8(first)),
		// 	int8(bits.OnesCount8(second)),
		// 	int8(bits.OnesCount8(third))}
		// intComp := compToInt(popCounts, n)
		if n > 0 && first&second == 0 && second&third == 0 && first&third == 0 {
			mask := ^(first + second + third)
			firstSqueezed := squeeze(mask, first, Little)
			secondSqueezed := squeeze(mask, second, Little)
			thirdSqueezed := squeeze(mask, third, Little)
			codeSqueezed := uint32(firstSqueezed) + (uint32(secondSqueezed) << 8) + (uint32(thirdSqueezed) << 16)
			countMap[n-1][codeSqueezed]++
			count++
			countBySize[n-1]++
		}
	}
	fmt.Println(count)
	fmt.Println(countBySize)
	mapSum := 0
	for i := 0; i < maxSuitSize; i++ {
		mapSum += len(countMap[i])
		fmt.Printf("size = %d, counts = %d\n", i+1, len(countMap[i]))
	}
	fmt.Println(mapSum)
	// Output:
	// 65535
	// [24 252 1512 5670 13608 20412 17496 6561]
	// size = 1, counts = 3
	// size = 2, counts = 9
	// size = 3, counts = 27
	// size = 4, counts = 81
	// size = 5, counts = 243
	// size = 6, counts = 729
	// size = 7, counts = 2187
	// size = 8, counts = 6561
	// 9840
}

func TestChase(t *testing.T) {
	for n := 0; n <= byteLen; n++ {
		for s := 0; s <= n; s++ {
			for i := 0; i < int(Binomial(int8(n), int8(s))); i++ {
				code := chase.sequence[n][s][i]
				require.EqualValues(t, i, chase.sequenceToIndex[n][s][code])
				if i > 0 {
					require.EqualValues(t, 2, bits.OnesCount8(uint8(code^chase.sequence[n][s][i-1])))
				}
			}
		}
	}
}

func BenchmarkChaseSeq(b *testing.B) {
	var SuitLength = 8
	for i := 0; i < b.N; i++ {
		first := rand.Intn(SuitLength)
		second := rand.Intn(SuitLength - first)
		third := rand.Intn(SuitLength - first - second)
		index := rand.Intn(int(MultiBinomial([]int8{int8(first), int8(second), int8(third)})))
		firstIndex := index / int(Binomial(int8(second+third), int8(second)))
		_ = chase.sequence[first+second+third][second+third][firstIndex]
	}
}

func BenchmarkSort(b *testing.B) {
	size := 4
	maxInt := 100
	rand.Seed(time.Now().Unix())
	array := make([]int, size)
	for i := 0; i < b.N; i++ {
		for j := 0; j < size; j++ {
			array[j] = rand.Intn(maxInt)
		}
		sort.Ints(array)
	}
}

func BenchmarkPopcount8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j <= 0xFF; j++ {
			bits.OnesCount8(uint8(j))
		}
	}
}

func BenchmarkPopcount16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j <= 0xFFFF; j++ {
			bits.OnesCount16(uint16(j))
		}
	}
}

func BenchmarkSetBit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j <= 0xFF; j++ {
			n := uint8(j)
			for k := 0; k < 8; k++ {
				n |= (1 << k)
			}
		}
	}
}

func BenchmarkUnsetBit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j <= 0xFF; j++ {
			n := uint8(j)
			var k uint8
			for ; k < 8; k++ {
				n &= ^(1 << k)
			}
		}
	}
}
