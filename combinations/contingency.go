package comb

import (
	"log"
)

type ContingencyTable [][]int8

// Visits all contingency tables with given sums of rows and columns;
// see Knuth, Art of Programming, 7.2.1.3, problem 62
func VisitContingencyTables(rowSums, colSums []int8, visitor func(ContingencyTable)) {
	var rowSum, colSum int8
	for i := 0; i < len(rowSums); i++ {
		rowSum += rowSums[i]
	}
	for _, colItem := range colSums {
		colSum += colItem
	}
	if rowSum != colSum {
		log.Fatalf("row and column sums must be equal: %d != %d", rowSum, colSum)
	}

	if len(rowSums) == 1 {
		visitor(ContingencyTable{colSums})
	} else {
		VisitBoundedCompositions(rowSums[0], colSums, func(array composition) {
			complement := make([]int8, len(colSums))
			complementReduced := make([]int8, 0, len(colSums))
			for j := range array {
				complement[j] = colSums[j] - array[j]
				if complement[j] > 0 {
					complementReduced = append(complementReduced, complement[j])
				}
			}
			VisitContingencyTables(rowSums[1:], complementReduced, func(subTable ContingencyTable) {
				table := make([][]int8, len(rowSums))
				for i := 0; i < len(rowSums); i++ {
					table[i] = make([]int8, len(colSums))
					if i == 0 {
						table[i] = array
					} else {
						colIdx := 0
						for j := 0; j < len(colSums); j++ {
							if complement[j] == 0 {
								table[i][j] = 0
							} else {
								table[i][j] = subTable[i-1][colIdx]
								colIdx++
							}
						}
					}
				}
				visitor(table)
			})
		})
	}
}

func rowsAreSorted(rowSums []int8, table ContingencyTable, trump bool) bool {
	var i int
	if trump {
		i = 1
	}
	for ; i < len(rowSums)-1; i++ {
		if rowSums[i] == rowSums[i+1] {
			for j := 0; j < len(table[i])-1; j++ {
				if table[i][j] < table[i+1][j] {
					break
				}
				if table[i][j] > table[i+1][j] {
					return false
				}
			}
		}
	}
	return true
}

// type tableInt uint32

// func TableToInt(rowSums []int8, table ContingencyTable) tableInt {
// 	var result, shift tableInt
// 	for i := 0; i < len(rowSums); i++ {
// 		result |= tableInt(compTriplets[rowSums[i]-1].getIndex(table[i])) << shift
// 		shift += 8
// 	}
// 	return result
// }

// func IntToTable(code tableInt, nCols int8, rowSums []int8) ContingencyTable {
// 	table := make(ContingencyTable, len(rowSums))
// 	for i := 0; i < len(rowSums); i++ {
// 		table[i] = make(composition, nCols)
// 		compTriplets[rowSums[i]-1].writeComp(int(code&0xFF), table[i])
// 		code >>= 8
// 	}
// 	return table
// }

// func writeIntToTable(code tableInt, nCols int8, rowSums []int8, table ContingencyTable) {
// 	for i := 0; i < len(rowSums); i++ {
// 		compTriplets[rowSums[i]-1].writeComp(int(code&0xFF), table[i])
// 		code >>= 8
// 	}
// }

func TableToInt(rowSums []int8, table ContingencyTable) int {
	result := 0
	pow := 1
	for i := 0; i < len(rowSums); i++ {
		result += pow * compTriplets[rowSums[i]-1].getIndex(table[i])
		pow *= compTriplets[rowSums[i]-1].len()
		// fmt.Printf("After row %d: result = %d, pow = %d\n", i, result, pow)
	}
	return result
}

func IntToTable(code int, nCols int8, rowSums []int8) ContingencyTable {
	table := make(ContingencyTable, len(rowSums))
	for i := 0; i < len(rowSums); i++ {
		factor := compTriplets[rowSums[i]-1].len()
		table[i] = make(composition, nCols)
		// table[i] = compTriplets[rowSums[i]-1].getComp(code % factor)
		compTriplets[rowSums[i]-1].writeComp(code%factor, table[i])
		code /= factor
	}
	return table
}

func writeIntToTable(code int, nCols int8, rowSums []int8, table ContingencyTable) {
	for i := 0; i < len(rowSums); i++ {
		factor := compTriplets[rowSums[i]-1].len()
		compTriplets[rowSums[i]-1].writeComp(code%factor, table[i])
		code /= factor
	}
}
