package comb

import (
	"log"
)

const (
	maxSuitSize   = 8
	triplet       = 3
	totalTriplets = 164
)

type composition []int8

// Visits all compositions (ordered partitions) of integer n into s non-negative summands;
// see Knuth, Art of Programming, 7.2.1.3, problem 3
func VisitCompositions(n, s int8, visitor func(composition)) {
	if s <= 0 {
		log.Fatalf("Number of summands should be positive!")
	}
	if s == 1 {
		visitor(composition{n})
		return
	}
	array := make(composition, s)
	var r int8 = -1
	array[0] = n
	for {
		visitor(array)
		if array[0] == 0 {
			if r == s-1 {
				break
			} else {
				array[0] = array[r] - 1
				array[r] = 0
				r++
			}
		} else {
			array[0]--
			r = 1
		}
		array[r]++
	}
}

func countCompositions(nMax, s int8) int {
	result := 0
	var n int8
	for n = 1; n <= nMax; n++ {
		VisitCompositions(n, s, func(composition) { result++ })
	}
	return result
}

type compInt uint16

const bitsPerPart = 4

// stores all compositions n = q[0] + ... + q[s-1]
// each composition is represented by int N = q[0] + q[1]*(n+1) + ... + q[s-1]*(n+1)^{s-1}
// requirements: s <= 4, n <= 15
type compositionIndexer struct {
	n        int8
	s        int8
	list     []compInt
	indexMap map[compInt]int
}

func compToInt(comp composition, n int8) compInt {
	var result, shift compInt
	for i := range comp {
		result |= compInt(comp[i]) << shift
		shift += bitsPerPart
	}
	// fmt.Printf("int(%v) = %v\n", comp, result)
	return result
}

func intToComp(code compInt, n, s int8, comp composition) {
	for i := 0; i < int(s); i++ {
		comp[i] = int8(code % (1 << bitsPerPart))
		code >>= bitsPerPart
	}
	// fmt.Printf("comp(%v) = %v\n", code, comp)
}

func newCompositionIndexer(n, s int8) compositionIndexer {
	nComps := Binomial(n+s-1, s-1)
	list := make([]compInt, nComps)
	indexMap := make(map[compInt]int, nComps)
	index := 0
	VisitCompositions(n, s, func(comp composition) {
		list[index] = compToInt(comp, n)
		indexMap[list[index]] = index
		index++
	})
	return compositionIndexer{n, s, list, indexMap}
}

func (ci compositionIndexer) getIndex(comp composition) int {
	return ci.indexMap[compToInt(comp, ci.n)]
}

func (ci compositionIndexer) getComp(index int) composition {
	comp := make(composition, ci.s)
	intToComp(ci.list[index], ci.n, ci.s, comp)
	return comp
}

func (ci compositionIndexer) writeComp(index int, comp composition) {
	intToComp(ci.list[index], ci.n, ci.s, comp)
}

func (ci compositionIndexer) len() int {
	return len(ci.list)
}

var compTriplets [maxSuitSize]compositionIndexer

func init() {
	var n int8
	for n = 1; n <= maxSuitSize; n++ {
		compTriplets[n-1] = newCompositionIndexer(n, triplet)
	}
}
