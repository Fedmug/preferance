package combinations

// Visits all compositions (ordered partitions) of the form number = q[1] + ... + q[s]
// with restriction 0 <= q[j] <= bound[j]; see Knuth, Art of Programming, 7.2.1.3, problem 60
func VisitBoundedCompositions(n int8, bound []int8, visitor func([]int8)) {
	var sum int8
	for _, b := range bound {
		sum += b
	}

	// Q1
	if sum < n {
		return
	}

	s := len(bound)
	if s == 1 {
		visitor([]int8{n})
		return
	}

	q := make([]int8, s)
	x := n
	for {
		// Q2
		j := 0
		for x > bound[j] {
			q[j] = bound[j]
			x -= bound[j]
			j++
		}
		q[j] = x

		// Q3
		for {
			visitor(q)
			flag := false

			// Q4
			if j == 0 {
				x = q[0] - 1
				j = 1
			} else if q[0] == 0 {
				x = q[j] - 1
				q[j] = 0
				j++
			} else {
				flag = true
			}

			if !flag {
				// Q5
				for j < s && q[j] == bound[j] {
					x += bound[j]
					q[j] = 0
					j++
				}
				if j >= s {
					return
				}

				// Q6
				q[j]++
				if x == 0 {
					q[0] = 0
					continue
				} else {
					break
				}
			}
			// Q7
			for q[j] == bound[j] {
				j++
				if j >= s {
					return
				}
			}
			q[j]++
			j--
			q[j]--
			if q[0] == 0 {
				j = 1
			}
		}
	}
}

func isOrdered(comp []int8, excludeFirstElement bool) bool {
	if excludeFirstElement && comp[0] == 0 {
		return false
	}
	i := 0
	if excludeFirstElement {
		i = 1
	}
	for ; i < len(comp)-1; i++ {
		if comp[i] > comp[i+1] {
			return false
		}
	}
	return true
}

// Given nCards of nSuits, counts all possible suit size distributions
// if maximal suit size is not greater than maxSuitSize
func SuitSizesCounter(nCards, nSuits, maxSuitSize int8, trump bool) int {
	result := 0
	suitBounds := make([]int8, nSuits)
	var i int8
	for ; i < nSuits; i++ {
		suitBounds[i] = maxSuitSize
	}
	VisitBoundedCompositions(nCards, suitBounds, func(suitSizes []int8) {
		if isOrdered(suitSizes, trump) {
			result++
		}
	})
	return result
}
