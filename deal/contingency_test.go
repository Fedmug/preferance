package deal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Example_db() {
	count := countTableLines("suit_sizes")
	fmt.Println(count)
	count = countTableLines("suit_sizes_v3")
	fmt.Println(count)
	count = countTableLines("contingency_tables")
	fmt.Println(count)
	// Output:
	// 602
	// 4680
	// 564020
}

func Example_db_matrix() {
	matrices := getMatrix()
	fmt.Println(len(matrices))
	// Output:
	// 550727
}

func Example_init_tables() {
	// initContingencyTables(nTables)
	fmt.Println(len(contingencyTables))
	table := contingencyTables[1983]
	fmt.Println(table, contingencyTableMap[table])
	table = contingencyTables[65535]
	fmt.Println(table, contingencyTableMap[table])
	// Output:
	// 564020
	// [[0 1 2] [0 3 3] [5 1 0] [0 0 0]] 1983
	// [[1 1 1] [1 1 1] [2 5 0] [3 0 5]] 65535
}

func Example_init_pool() {
	// initContingencyTablesPool()
	fmt.Println(len(contingencyTables))
	table := contingencyTables[1983]
	fmt.Println(table, contingencyTableMap[table])
	table = contingencyTables[65535]
	fmt.Println(table, contingencyTableMap[table])
	// Output:
	// 564020
	// [[0 1 2] [0 3 3] [5 1 0] [0 0 0]] 1983
	// [[1 1 1] [1 1 1] [2 5 0] [3 0 5]] 65535
}

func TestContingencyTables(t *testing.T) {
	// initContingencyTablesPool()
	for i := 0; i < nTables; i++ {
		table := contingencyTables[i]
		index := contingencyTableMap[table]
		require.EqualValues(t, i, index, fmt.Sprintf("table index is wrong: got %d, want %d", index, i))
		var deckSize int8
		for suit := 0; suit < NumberOfSuits; suit++ {
			for hand := 0; hand < NumberOfHands; hand++ {
				deckSize += table[suit][hand]
				require.LessOrEqual(t, table[suit][hand], int8(SuitLength))
				require.GreaterOrEqual(t, table[suit][hand], int8(0))
			}
		}
		require.EqualValues(t, 0, deckSize%NumberOfHands)
	}
}
