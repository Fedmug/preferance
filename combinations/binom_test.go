package combinations

import (
	"fmt"
	"testing"
)

func ExampleBinomial() {
	fmt.Println(Binomial(5, 2))
	fmt.Println(Binomial(20, 10))
	// Output:
	// 10
	// 184756
}

type testBinomCase struct {
	n, k   uint8
	result uint64
}

var binomCases = []testBinomCase{
	{n: 0, k: 0, result: 1},
	{n: 5, k: 10, result: 0},
	{n: 12, k: 0, result: 1},
	{n: 6, k: 1, result: 6},
	{n: 65, k: 65, result: 1},
	{n: 15, k: 14, result: 15},
	{n: 4, k: 2, result: 6},
	{n: 8, k: 4, result: 70},
	{n: 20, k: 10, result: 184756},
	{n: 32, k: 10, result: 64512240},
	{n: 52, k: 13, result: 635013559600},
}

func TestBinomSimple(t *testing.T) {
	for _, input := range binomCases {
		if out := Binomial(input.n, input.k); out != input.result {
			t.Errorf("wrong binomial(%v, %v): got %d, want %d", input.n, input.k, out, input.result)
		}
	}
}

type testMultinomialCase struct {
	ns     []uint8
	result uint64
}

func TestMulinomial(t *testing.T) {
	for _, input := range []testMultinomialCase{
		{ns: []uint8{5, 0}, result: 1},
		{ns: []uint8{2, 3}, result: 10},
		{ns: []uint8{2, 3, 3}, result: 560},
		{ns: []uint8{2, 2, 2, 2}, result: 2520},
		{ns: []uint8{10, 10, 10, 2}, result: 2753294408504640},
	} {
		if out := Multinomial(input.ns); out != input.result {
			t.Errorf("multinomial(%v) = %d != %d", input.ns, out, input.result)
		}
		if out := MultiBinomial(input.ns); out != input.result {
			t.Errorf("multinomial(%v) = %d != %d", input.ns, out, input.result)
		}
	}
}

func BenchmarkBinomial_8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Binomial(8, 3)
	}
}

func BenchmarkBinomial_32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Binomial(32, 10)
	}
}

func BenchmarkMultinomial(b *testing.B) {
	array := [3]uint8{3, 2, 3}
	for i := 0; i < b.N; i++ {
		if x := Multinomial(array[:3]); x < 1 {
			b.Fatalf("To small: %d", x)
		}
	}
}

func BenchmarkMultiBinomial(b *testing.B) {
	array := [3]uint8{3, 2, 3}
	for i := 0; i < b.N; i++ {
		if x := MultiBinomial(array[:3]); x < 1 {
			b.Fatalf("To small: %d", x)
		}
	}
}
