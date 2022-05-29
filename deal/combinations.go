package deal

import (
	"log"
)

const (
	byteLen = 8
	byteCap = 256
)

type Endian int8

const (
	little Endian = iota
	big
)

var BinomialCoefficients [byteLen + 1][byteLen + 1]int

type SuitHandCodeSqueezer struct {
	endian         Endian
	squeezeTable   [byteCap][byteCap]SuitHandCode
	unsqueezeTable [byteCap][byteCap]SuitHandCode
}

var squeezerLittle SuitHandCodeSqueezer
var squeezerBig SuitHandCodeSqueezer

func (sq *SuitHandCodeSqueezer) initTables(endian Endian) {
	sq.endian = endian
	for first := 0; first < byteCap; first++ {
		for second := 0; second < byteCap; second++ {
			if first&second == 0 {
				squeezedCode := squeezeSecondHandCode(SuitHandCode(first), SuitHandCode(second), endian)
				sq.squeezeTable[first][second] = squeezedCode
				sq.unsqueezeTable[first][squeezedCode] = SuitHandCode(second)
			}
		}
	}
}

// Given mask and code, retrieve squeezed code, e. g.
// 00110100, 01000010 -> 00001010 if endian = little
// 00110100, 01000010 -> 01010000 if endian = big
func (sq *SuitHandCodeSqueezer) squeeze(mask, code SuitHandCode) SuitHandCode {
	return sq.squeezeTable[mask][code]
}

// Given mask and code, retrieve unsqueezed code, e. g.
// 00010100, 00001101 -> 00101001 if endian = little
// 00010100, 00110100 -> 00101001 if endian = big
func (sq *SuitHandCodeSqueezer) unsqueeze(mask, code SuitHandCode) SuitHandCode {
	return sq.unsqueezeTable[mask][code]
}

type ChaseSequences struct {
	sequence        [byteLen + 1][byteLen + 1][]SuitHandCode
	sequenceToIndex [byteLen + 1][byteLen + 1][byteCap]int
}

var chase ChaseSequences

func squeezeSecondHandCode(firstHandCode, secondHandCode SuitHandCode, endian Endian) SuitHandCode {
	var i, result, shift SuitHandCode
	for i = 0; i < byteLen; i++ {
		index := i
		if endian == big {
			index = byteLen - 1 - i
		}
		if firstHandCode&(1<<index) == 0 {
			if endian == little {
				result |= (secondHandCode & (1 << index)) >> shift
			} else if endian == big {
				result |= (secondHandCode & (1 << index)) << shift
			} else {
				log.Fatalf("unknown endian: %v", endian)
			}
		} else {
			shift++
		}
	}
	return result
}

func squeeze(firstHandCode, secondHandCode SuitHandCode, endian Endian) SuitHandCode {
	if endian == little {
		// return squeezeSecondHandTableLittle[firstHandCode][secondHandCode]
		return squeezerLittle.squeeze(firstHandCode, secondHandCode)
	}
	// return squeezeSecondHandTableBig[firstHandCode][secondHandCode]
	return squeezerBig.squeeze(firstHandCode, secondHandCode)
}

func unsqueeze(firstHandCode, secondHandCode SuitHandCode, endian Endian) SuitHandCode {
	if endian == little {
		// return unsqueezeSecondHandTableLittle[firstHandCode][secondHandCode]
		return squeezerLittle.unsqueeze(firstHandCode, secondHandCode)
	}
	// return unsqueezeSecondHandTableBig[firstHandCode][secondHandCode]
	return squeezerBig.unsqueeze(firstHandCode, secondHandCode)
}

func initBinomials() {
	for i := 0; i <= byteLen; i++ {
		BinomialCoefficients[i][0] = 1
		BinomialCoefficients[i][i] = 1
		for j := 1; j < i; j++ {
			BinomialCoefficients[i][j] = BinomialCoefficients[i-1][j] + BinomialCoefficients[i-1][j-1]
		}
	}
}

func initChase() {
	for i := 0; i <= byteLen; i++ {
		for j := 0; j <= byteLen; j++ {
			for k := 0; k < byteCap; k++ {
				chase.sequenceToIndex[i][j][k] = -1
			}
		}
	}
	for n := 0; n <= byteLen; n++ {
		for s := 0; s <= n; s++ {
			chase.sequence[n][s] = make([]SuitHandCode, 0, BinomialCoefficients[n][s])
			var code SuitHandCode
			for i := 0; i < n-s; i++ {
				code += 1 << (s + i)
			}
			w := make([]uint8, n+1)
			for i := 0; i <= n; i++ {
				w[i] = 1
			}
			r := s
			// begin, end := -1, -1
			if r == 0 {
				r = n - s
			}
			for i := 0; i < BinomialCoefficients[n][s]; i++ {
				chase.sequence[n][s] = append(chase.sequence[n][s], code)
				chase.sequenceToIndex[n][s][code] = i
				j := r
				for w[j] == 0 {
					w[j] = 1
					j++
				}
				if j == n {
					break
				}
				w[j] = 0
				if code&(1<<j) > 0 {
					if j%2 == 1 || code&(1<<(j-2)) > 0 {
						code -= 1 << j
						code += 1 << (j - 1)
						// begin, end = n-1-j, n-j
						if r == j && j > 1 {
							r = j - 1
						} else if r == j-1 {
							r = j
						}
					} else {
						code -= 1 << j
						code += 1 << (j - 2)
						// begin, end = n-1-j, n+1-j
						if r == j {
							r = 1
							if j-2 > 1 {
								r = j - 2
							}
						} else if r == j-2 {
							r = j - 1
						}
					}
				} else {
					if j%2 == 0 || code&(1<<(j-1)) > 0 {
						code += 1 << j
						code -= 1 << (j - 1)
						// begin, end = n-1-j, n-j
						if r == j && j > 1 {
							r = j - 1
						} else if r == j-1 {
							r = j
						}
					} else {
						code += 1 << j
						code -= 1 << (j - 2)
						// begin, end = n-1-j, n+1-j
						if r == j-2 {
							r = j
						} else if r == j-1 {
							r = j - 2
						}
					}
				}
			}
		}
	}
}

func init() {
	squeezerLittle.initTables(little)
	squeezerBig.initTables(big)
	initBinomials()
	initChase()
}

func suitToIndex(handSizes [NumberOfHands]int8, handCodes [NumberOfHands]SuitHandCode, endian Endian) int {
	n := handSizes[0] + handSizes[1] + handSizes[2]
	squeezedSecondHandCode := squeeze(handCodes[0], handCodes[1], endian)
	firstIndex := chase.sequenceToIndex[n][n-handSizes[0]][handCodes[0]]
	secondIndex := chase.sequenceToIndex[n-handSizes[0]][handSizes[2]][squeezedSecondHandCode]
	return firstIndex*BinomialCoefficients[n-handSizes[0]][handSizes[1]] + secondIndex
}

func indexToSuitCodes(handSizes [NumberOfHands]int8, index int, endian Endian) [NumberOfHands]SuitHandCode {
	n := handSizes[0] + handSizes[1] + handSizes[2]
	firstIndex := index / BinomialCoefficients[n-handSizes[0]][handSizes[1]]
	secondIndex := index % BinomialCoefficients[n-handSizes[0]][handSizes[1]]
	var result [NumberOfHands]SuitHandCode
	result[0] = chase.sequence[n][n-handSizes[0]][firstIndex]
	result[1] = unsqueeze(result[0], chase.sequence[n-handSizes[0]][handSizes[2]][secondIndex], endian)
	sum := (1 << n) - 1
	if endian == big {
		sum <<= 8 - n
	}
	result[2] = SuitHandCode(sum) - result[0] - result[1]
	return result
}

/*
func multinomial(ns []uint8) uint64 {
	//sort.Slice(ns, func(i, j int) bool { return ns[i] < ns[j] })
	var sum uint8
	for _, n := range ns {
		sum += n
	}
	var result uint64 = 1
	i := sum
	for _, n := range ns[:len(ns)-1] {
		var j uint8
		for j = 1; j <= n; j++ {
			result *= uint64(i)
			result /= uint64(j)
			i -= 1
		}
	}
	return result
}

/*
func multipermToIndex(multiperm []uint8, ns []uint8) uint64 {

	   Calculates index of a multiset permutation,
	   see Knuth, Art of Programming, 7.2.1.2, ex. 4
	   :param ns: sizes of elements of the multiset, ns = (n_1, ..,, n_k)
	   :param multiperm: permutation of the multiset {n_1*0, ..., n_k*(k-1)}
	   :return: uint64, index, 0 <= index < (n; ns) (multinomial coefficient)

	var result uint64
	for _, element := range multiperm {
		multi_coef := multinomial(ns)
		var n uint8
		for _, size := range ns {
			n += size
		}
		var j uint8
		for ; j < element; j++ {
			result += multi_coef * uint64(ns[j]) / uint64(n)
		}
		ns[element]--
	}
	return result
}

/*
func indexToMultiperm(index uint64, ns []uint8) []uint8 {

        Calculates multiset permutation by its index,
        see Knuth, Art of Programming, 7.2.1.2, ex. 4
        :param ns: sizes of elements of the multiset, ns = (n_1, ..,, n_k)
        :param index: 0 <= index < (n; ns) (multinomial coefficient)
        :return: permutation of the multiset {n_1*0, ..., n_k*(k-1)}

	var n uint8
	var i int = -1
	for _, size := range ns {
		n += size
	}
    if n == 1 {
		result := make([]uint8, 1)
        return np.array(np.argmax(ns))
	}
    multi_coef = multinomial(ns)
    assert index < multi_coef, "Index should be less than multinomial coefficient!"
    cum_sum = 0
    for j in range(len(ns)):
        N_j = multi_coef * ns[j] // ns.sum()
        if index < cum_sum + N_j:
            sub_ns = np.copy(ns)
            sub_ns[j] -= 1
            if sub_ns[j] == 0:
                np.delete(sub_ns, j)
            sub_result = index2multiperm(index - cum_sum, sub_ns)
            return np.append(j, sub_result)
        cum_sum += N_j
}
*/
