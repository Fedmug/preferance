package combinations

import (
	"fmt"
	"math/bits"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func Example_chase() {
	var n, s = 4, 2
	var i uint64
	for i = 0; i < Binomial(uint8(n), uint8(s)); i++ {
		fmt.Printf("%04b %d\n", chase.sequence[n][s][i], chase.sequenceToIndex[n][s][chase.sequence[n][s][i]])
	}
	n, s = 3, 0
	for i = 0; i < Binomial(uint8(n), uint8(s)); i++ {
		fmt.Printf("%04b %d\n", chase.sequence[n][s][i], chase.sequenceToIndex[n][s][chase.sequence[n][s][i]])
	}
	n, s = 3, 3
	for i = 0; i < Binomial(uint8(n), uint8(s)); i++ {
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

func TestChase(t *testing.T) {
	for n := 0; n <= byteLen; n++ {
		for s := 0; s <= n; s++ {
			for i := 0; i < int(Binomial(uint8(n), uint8(s))); i++ {
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
		index := rand.Intn(int(MultiBinomial([]uint8{uint8(first), uint8(second), uint8(third)})))
		firstIndex := index / int(Binomial(uint8(second+third), uint8(second)))
		_ = chase.sequence[first+second+third][second+third][firstIndex]
	}
}
