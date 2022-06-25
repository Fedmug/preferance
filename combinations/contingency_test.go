package comb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExampleVisitContingencyTables_trivial() {
	VisitContingencyTables([]int8{6}, []int8{1, 2, 3}, func(table ContingencyTable) {
		fmt.Println(table)
	})
	VisitContingencyTables([]int8{8}, []int8{1, 0, 4, 3}, func(table ContingencyTable) {
		fmt.Println(table)
	})
	// Output:
	// [[1 2 3]]
	// [[1 0 4 3]]
}

func ExampleVisitContingencyTables_two() {
	VisitContingencyTables([]int8{2, 3}, []int8{1, 2, 2}, func(table ContingencyTable) {
		fmt.Println(table)
	})
	// Output:
	// [[1 1 0] [0 1 2]]
	// [[0 2 0] [1 0 2]]
	// [[1 0 1] [0 2 1]]
	// [[0 1 1] [1 1 1]]
	// [[0 0 2] [1 2 0]]
}

func ExampleVisitContingencyTables_three() {
	VisitContingencyTables([]int8{1, 1, 4}, []int8{2, 2, 2}, func(table ContingencyTable) {
		fmt.Println(table)
	})
	// Output:
	// [[1 0 0] [1 0 0] [0 2 2]]
	// [[1 0 0] [0 1 0] [1 1 2]]
	// [[1 0 0] [0 0 1] [1 2 1]]
	// [[0 1 0] [1 0 0] [1 1 2]]
	// [[0 1 0] [0 1 0] [2 0 2]]
	// [[0 1 0] [0 0 1] [2 1 1]]
	// [[0 0 1] [1 0 0] [1 2 1]]
	// [[0 0 1] [0 1 0] [2 1 1]]
	// [[0 0 1] [0 0 1] [2 2 0]]
}

func ExampleVisitContingencyTables_many() {
	c := 0
	VisitContingencyTables([]int8{8, 4, 7, 8}, []int8{9, 9, 9}, func(ContingencyTable) {
		c += 1
	})
	fmt.Println(c)
	// Output:
	// 12000
}

func ExampleVisitContingencyTables_groupBySuitSizes() {
	trump := true
	total := 0
	totalReduced := 0
	for nSuits := 1; nSuits <= 4; nSuits++ {
		var stage int8
		for stage = 1; stage <= 10; stage++ {
			bounds := make([]int8, nSuits)
			for i := 0; i < nSuits; i++ {
				bounds[i] = 8
			}
			VisitBoundedCompositions(3*stage, bounds, func(comp composition) {
				for i := 0; i < nSuits; i++ {
					if comp[i] == 0 || !isOrdered(comp, trump) {
						return
					}
				}
				colSums := []int8{stage, stage, stage}
				count := 0
				countReduced := 0
				VisitContingencyTables(comp, colSums, func(table ContingencyTable) {
					count++
					if rowsAreSorted(comp, table, trump) {
						countReduced++
					}
				})
				total += count
				totalReduced += countReduced
				// fmt.Println(comp, count)
			})
		}
	}
	fmt.Println("Total:", total)
	fmt.Println("Total reduced:", totalReduced)
	// Output:
	// Total: 859391
	// Total reduced: 564020
}

func ExampleTableToInt() {
	table := ContingencyTable{{1, 0, 1}, {1, 2, 1}}
	code := TableToInt([]int8{2, 4}, table)
	fmt.Println(code)
	fmt.Println(IntToTable(code, triplet, []int8{2, 4}))

	rowSums := []int8{7, 6, 6, 8}
	table = ContingencyTable{{1, 6, 0}, {0, 1, 5}, {2, 0, 4}, {6, 2, 0}}
	code = TableToInt(rowSums, table)
	fmt.Println(code)
	fmt.Println(IntToTable(code, triplet, rowSums))
	// Output:
	// 45
	// [[1 0 1] [1 2 1]]
	// 79566
	// [[1 6 0] [0 1 5] [2 0 4] [6 2 0]]
}

// Output:
// 1795
// [[1 0 1] [1 2 1]]
// 35002886
// [[1 6 0] [0 1 5] [2 0 4] [6 2 0]]

func TestTableInt(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip long table int test")
	}
	trump := true
	var maxCode int
	for nSuits := 1; nSuits <= 4; nSuits++ {
		var stage int8
		for stage = 1; stage <= 10; stage++ {
			bounds := make([]int8, nSuits)
			for i := 0; i < nSuits; i++ {
				bounds[i] = 8
			}
			VisitBoundedCompositions(3*stage, bounds, func(comp composition) {
				for i := 0; i < nSuits; i++ {
					if comp[i] == 0 || !isOrdered(comp, trump) {
						return
					}
				}
				colSums := []int8{stage, stage, stage}
				VisitContingencyTables(comp, colSums, func(table ContingencyTable) {
					code := TableToInt(comp, table)
					if code > maxCode {
						maxCode = code
					}
					tableFromCode := IntToTable(code, triplet, comp)
					require.EqualValues(t, table, tableFromCode)
				})
			})
		}
	}
	fmt.Println("Max code:", maxCode)
}

func BenchmarkTableToInt(b *testing.B) {
	rowSums := []int8{7, 6, 6, 8}
	table := ContingencyTable{{1, 6, 0}, {0, 1, 5}, {2, 0, 4}, {6, 2, 0}}
	for i := 0; i < b.N; i++ {
		TableToInt(rowSums, table)
	}
}

func BenchmarkIntToTable(b *testing.B) {
	rowSums := []int8{7, 6, 6, 8}
	for i := 0; i < b.N; i++ {
		IntToTable(79566, triplet, rowSums)
		// IntToTable(35002886, triplet, rowSums)
	}
}

func BenchmarkWriteIntToTable(b *testing.B) {
	rowSums := []int8{7, 6, 6, 8}
	table := make(ContingencyTable, len(rowSums))
	for i := 0; i < len(rowSums); i++ {
		table[i] = make(composition, 3)
	}
	for i := 0; i < b.N; i++ {
		writeIntToTable(79566, triplet, rowSums, table)
		// writeIntToTable(35002886, triplet, rowSums, table)
	}
}
