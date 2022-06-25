package comb

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func compositionsPrinter(numbers composition) {
	var sum int8
	for _, n := range numbers {
		sum += n
	}
	fmt.Printf("%d = ", sum)
	for i, n := range numbers {
		if i != 0 {
			fmt.Printf("+")
		}
		fmt.Printf("%d", n)
	}
	fmt.Printf("\n")
}

func Example_countCompositions() {
	fmt.Println(countCompositions(8, 2))
	fmt.Println(countCompositions(8, 3))
	fmt.Println(countCompositions(8, 4))
	fmt.Println(countCompositions(13, 4))
	// Output:
	// 44
	// 164
	// 494
	// 2379
}

func Example_triplets_count() {
	for n := 0; n < maxSuitSize; n++ {
		fmt.Println(compTriplets[n].len())
	}
	// Output:
	// 3
	// 6
	// 10
	// 15
	// 21
	// 28
	// 36
	// 45
}

func Example_sortTriplets() {
	triplets := make([][triplet]int8, totalTriplets)
	i := 0
	for n := 0; n < maxSuitSize; n++ {
		for j := 0; j < compTriplets[n].len(); j++ {
			compTriplets[n].writeComp(j, triplets[i][:])
			i++
		}
	}
	sort.Slice(triplets, func(i, j int) bool {
		firstSum := triplets[i][0] + triplets[i][1] + triplets[i][2]
		secondSum := triplets[j][0] + triplets[j][1] + triplets[j][2]
		if firstSum < secondSum {
			return true
		}
		if firstSum > secondSum {
			return false
		}
		for k := 0; k < triplet; k++ {
			if triplets[i][k] < triplets[j][k] {
				return true
			}
			if triplets[i][k] > triplets[j][k] {
				return false
			}
		}
		return true
	})
	for _, comp := range triplets {
		fmt.Println(comp)
	}
}

func Example_compositionsIndexer_1() {
	size := 1
	for j := 0; j < compTriplets[size-1].len(); j++ {
		fmt.Println(compTriplets[size-1].getComp(j))
	}
	// Output:
	// [1 0 0]
	// [0 1 0]
	// [0 0 1]
}

func Example_compositionsIndexer_2() {
	size := 2
	for j := 0; j < compTriplets[size-1].len(); j++ {
		fmt.Println(compTriplets[size-1].getComp(j))
	}
	// Output:
	// [2 0 0]
	// [1 1 0]
	// [0 2 0]
	// [1 0 1]
	// [0 1 1]
	// [0 0 2]
}

func Example_compositionsIndexer_3() {
	size := 3
	for j := 0; j < compTriplets[size-1].len(); j++ {
		fmt.Println(compTriplets[size-1].getComp(j))
	}
	// Output:
	// [3 0 0]
	// [2 1 0]
	// [1 2 0]
	// [0 3 0]
	// [2 0 1]
	// [1 1 1]
	// [0 2 1]
	// [1 0 2]
	// [0 1 2]
	// [0 0 3]
}

func Example_compositionsIndexer_4() {
	size := 4
	for j := 0; j < compTriplets[size-1].len(); j++ {
		fmt.Println(compTriplets[size-1].getComp(j))
	}
	// Output:
	// [4 0 0]
	// [3 1 0]
	// [2 2 0]
	// [1 3 0]
	// [0 4 0]
	// [3 0 1]
	// [2 1 1]
	// [1 2 1]
	// [0 3 1]
	// [2 0 2]
	// [1 1 2]
	// [0 2 2]
	// [1 0 3]
	// [0 1 3]
	// [0 0 4]
}

func Example_compositionsIndexer_5() {
	size := 5
	for j := 0; j < compTriplets[size-1].len(); j++ {
		fmt.Println(compTriplets[size-1].getComp(j))
	}
	// Output:
	// [5 0 0]
	// [4 1 0]
	// [3 2 0]
	// [2 3 0]
	// [1 4 0]
	// [0 5 0]
	// [4 0 1]
	// [3 1 1]
	// [2 2 1]
	// [1 3 1]
	// [0 4 1]
	// [3 0 2]
	// [2 1 2]
	// [1 2 2]
	// [0 3 2]
	// [2 0 3]
	// [1 1 3]
	// [0 2 3]
	// [1 0 4]
	// [0 1 4]
	// [0 0 5]
}

func Example_compositionsIndexer_6() {
	size := 6
	for j := 0; j < compTriplets[size-1].len(); j++ {
		fmt.Println(compTriplets[size-1].getComp(j))
	}
	// Output:
	// [6 0 0]
	// [5 1 0]
	// [4 2 0]
	// [3 3 0]
	// [2 4 0]
	// [1 5 0]
	// [0 6 0]
	// [5 0 1]
	// [4 1 1]
	// [3 2 1]
	// [2 3 1]
	// [1 4 1]
	// [0 5 1]
	// [4 0 2]
	// [3 1 2]
	// [2 2 2]
	// [1 3 2]
	// [0 4 2]
	// [3 0 3]
	// [2 1 3]
	// [1 2 3]
	// [0 3 3]
	// [2 0 4]
	// [1 1 4]
	// [0 2 4]
	// [1 0 5]
	// [0 1 5]
	// [0 0 6]
}

func Example_compositionsIndexer_7() {
	size := 7
	for j := 0; j < compTriplets[size-1].len(); j++ {
		fmt.Println(compTriplets[size-1].getComp(j))
	}
	// Output:
	// [7 0 0]
	// [6 1 0]
	// [5 2 0]
	// [4 3 0]
	// [3 4 0]
	// [2 5 0]
	// [1 6 0]
	// [0 7 0]
	// [6 0 1]
	// [5 1 1]
	// [4 2 1]
	// [3 3 1]
	// [2 4 1]
	// [1 5 1]
	// [0 6 1]
	// [5 0 2]
	// [4 1 2]
	// [3 2 2]
	// [2 3 2]
	// [1 4 2]
	// [0 5 2]
	// [4 0 3]
	// [3 1 3]
	// [2 2 3]
	// [1 3 3]
	// [0 4 3]
	// [3 0 4]
	// [2 1 4]
	// [1 2 4]
	// [0 3 4]
	// [2 0 5]
	// [1 1 5]
	// [0 2 5]
	// [1 0 6]
	// [0 1 6]
	// [0 0 7]
}

func Example_compositionsIndexer_8() {
	size := 8
	for j := 0; j < compTriplets[size-1].len(); j++ {
		fmt.Println(compTriplets[size-1].getComp(j))
	}
	// Output:
	// [8 0 0]
	// [7 1 0]
	// [6 2 0]
	// [5 3 0]
	// [4 4 0]
	// [3 5 0]
	// [2 6 0]
	// [1 7 0]
	// [0 8 0]
	// [7 0 1]
	// [6 1 1]
	// [5 2 1]
	// [4 3 1]
	// [3 4 1]
	// [2 5 1]
	// [1 6 1]
	// [0 7 1]
	// [6 0 2]
	// [5 1 2]
	// [4 2 2]
	// [3 3 2]
	// [2 4 2]
	// [1 5 2]
	// [0 6 2]
	// [5 0 3]
	// [4 1 3]
	// [3 2 3]
	// [2 3 3]
	// [1 4 3]
	// [0 5 3]
	// [4 0 4]
	// [3 1 4]
	// [2 2 4]
	// [1 3 4]
	// [0 4 4]
	// [3 0 5]
	// [2 1 5]
	// [1 2 5]
	// [0 3 5]
	// [2 0 6]
	// [1 1 6]
	// [0 2 6]
	// [1 0 7]
	// [0 1 7]
	// [0 0 8]
}

func ExampleVisitBoundedCompositions_trivial() {
	VisitBoundedCompositions(7, []int8{7}, compositionsPrinter)
	// Output:
	// 7 = 7
}

func ExampleVisitBoundedCompositions() {
	VisitBoundedCompositions(6, []int8{3, 2, 3}, compositionsPrinter)
	// Output:
	// 6 = 3+2+1
	// 6 = 3+1+2
	// 6 = 2+2+2
	// 6 = 3+0+3
	// 6 = 2+1+3
	// 6 = 1+2+3
}

func ExampleVisitBoundedCompositions_count() {
	c := 0
	VisitBoundedCompositions(15, []int8{8, 8, 8, 8}, func(comp composition) { c++ })
	fmt.Println(c)
	// Output:
	// 480
}

func ExampleSuitSizesCounter() {
	var handSize int8
	total := 0
	for handSize = 1; handSize <= 10; handSize++ {
		count := SuitSizesCounter(3*handSize, 4, 8, false)
		total += count
		fmt.Printf("c(%d) = %d\n", handSize, count)
	}
	fmt.Printf("total = %d\n", total)
	// Output:
	// c(1) = 3
	// c(2) = 9
	// c(3) = 17
	// c(4) = 27
	// c(5) = 31
	// c(6) = 31
	// c(7) = 23
	// c(8) = 15
	// c(9) = 6
	// c(10) = 2
	// total = 164
}

func ExampleSuitSizesCounter_trump() {
	var handSize int8
	total := 0
	for handSize = 1; handSize <= 10; handSize++ {
		count := SuitSizesCounter(3*handSize, 4, 8, true)
		total += count
		fmt.Printf("c(%d) = %d\n", handSize, count)
	}
	fmt.Printf("total = %d\n", total)
	// Output:
	// c(1) = 4
	// c(2) = 16
	// c(3) = 40
	// c(4) = 69
	// c(5) = 90
	// c(6) = 90
	// c(7) = 69
	// c(8) = 40
	// c(9) = 16
	// c(10) = 4
	// total = 438
}

func Example_countRows() {
	var total uint64
	var n int8
	for n = 1; n <= maxSuitSize; n++ {
		var variants uint64
		for _, compCode := range compTriplets[n-1].list {
			var comp [triplet]int8
			intToComp(compCode, n, triplet, comp[:])
			variants += MultiBinomial(comp[:])
		}
		fmt.Printf("row sum = %d: variants = %d\n", n, variants)
		total += variants
	}
	fmt.Println(total)
	// Output:
	// row sum = 1: variants = 3
	// row sum = 2: variants = 9
	// row sum = 3: variants = 27
	// row sum = 4: variants = 81
	// row sum = 5: variants = 243
	// row sum = 6: variants = 729
	// row sum = 7: variants = 2187
	// row sum = 8: variants = 6561
	// 9840
}

func TestCompInt(t *testing.T) {
	var n int8
	compFromCode := make(composition, triplet)
	for n = 1; n <= maxSuitSize; n++ {
		VisitCompositions(n, triplet, func(comp composition) {
			code := compToInt(comp, int8(n))
			intToComp(code, n, triplet, compFromCode)
			require.EqualValues(t, comp, compFromCode)
		})
	}
}

func TestBoundedCompositions(t *testing.T) {
	var n int8
	for n = 1; n <= 32; n++ {
		VisitBoundedCompositions(n, []int8{8, 8, 8, 8}, func(comb composition) {
			var sum int8
			require.EqualValues(t, 4, len(comb))
			for _, c := range comb {
				sum += c
				require.LessOrEqual(t, c, int8(8))
			}
			require.EqualValues(t, n, sum)
		})
	}
}

func BenchmarkBoundedCompostionsVisit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := 0
		VisitBoundedCompositions(15, []int8{8, 8, 8, 8}, func(composition) { c++ })
	}
}

func BenchmarkCompToInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for n := 1; n <= maxSuitSize; n++ {
			VisitCompositions(int8(n), triplet, func(comp composition) {
				compToInt(comp, int8(n))
			})
		}
	}
}

func BenchmarkIntToComp(b *testing.B) {
	s := 3
	comp := make(composition, s)
	for i := 0; i < b.N; i++ {
		for n := 1; n <= maxSuitSize; n++ {
			for j := 0; j < int(Binomial(int8(n+s-1), int8(s-1))); j++ {
				intToComp(compInt(j), int8(n), int8(s), comp)
			}
		}
	}
}
