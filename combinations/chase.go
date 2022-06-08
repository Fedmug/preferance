package combinations

// all n-choose-k Chase sequences are stored as binary codes
// see also Khuth, 7.2.1.3, Algorithm C
type ChaseSequences struct {
	sequence        [byteLen + 1][byteLen + 1][]uint8
	sequenceToIndex [byteLen + 1][byteLen + 1][byteCap]int
}

var chase ChaseSequences

func init() {
	for i := 0; i <= byteLen; i++ {
		for j := 0; j <= byteLen; j++ {
			for k := 0; k < byteCap; k++ {
				chase.sequenceToIndex[i][j][k] = -1
			}
		}
	}
	for n := 0; n <= byteLen; n++ {
		for s := 0; s <= n; s++ {
			chase.sequence[n][s] = make([]uint8, 0, Binomial(uint8(n), uint8(s)))
			var code uint8
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
			for i := 0; i < int(Binomial(uint8(n), uint8(s))); i++ {
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
