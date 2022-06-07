package combinations

const maxBinom = 8

var binomialCoefficients [maxBinom + 1][maxBinom + 1]int

func init() {
	for i := 0; i <= maxBinom; i++ {
		binomialCoefficients[i][0] = 1
		binomialCoefficients[i][i] = 1
		for j := 1; j < i; j++ {
			binomialCoefficients[i][j] = binomialCoefficients[i-1][j] + binomialCoefficients[i-1][j-1]
		}
	}
}

func Binomial(n, k uint8) uint64 {
	if k > n {
		return 0
	}
	if n <= maxBinom {
		return uint64(binomialCoefficients[n][k])
	}
	if k == 0 || k == n || n <= 1 {
		return 1
	}
	if k == 1 || k == n-1 {
		return uint64(n)
	}
	result := uint64(n - k + 1)
	for i, j := result+1, uint64(2); j <= uint64(k); i, j = i+1, j+1 {
		result = result * i / j
	}
	return result
}

func Multinomial(ns []uint8) uint64 {
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
