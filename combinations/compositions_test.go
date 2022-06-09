package combinations

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func compositionsPrinter(numbers []int8) {
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
	VisitBoundedCompositions(15, []int8{8, 8, 8, 8}, func(comp []int8) { c++ })
	fmt.Println(c)
	// Output:
	// 480
}

func ExampleVisitBoundedCompositions_count_triplets() {
	total := 0
	var i int8
	for i = 1; i <= 8; i++ {
		c := 0
		VisitBoundedCompositions(i, []int8{8, 8, 8}, func([]int8) { c++ })
		total += c
		fmt.Println(c)
	}
	fmt.Println("Total:", total)
	// Output:
	// 3
	// 6
	// 10
	// 15
	// 21
	// 28
	// 36
	// 45
	// Total: 164
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

func TestBoundedCompositions(t *testing.T) {
	var n int8
	for n = 1; n <= 32; n++ {
		VisitBoundedCompositions(n, []int8{8, 8, 8, 8}, func(comb []int8) {
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
		VisitBoundedCompositions(15, []int8{8, 8, 8, 8}, func([]int8) { c++ })
	}
}
