package functions

import (
	"math/rand"
)

var dp = func() []float64 {
	dp := make([]float64, 2, 1024)
	dp[0] = 1
	dp[1] = 0
	return dp
}()

func Subfactorial(n int) float64 {
	if n < len(dp) {
		return dp[n]
	}
	for i := len(dp); i <= n; i++ {
		dp = append(dp, float64(i-1)*(dp[i-1]+dp[i-2]))
	}
	return dp[n]
}

func RandomDerangement(n int, rng *rand.Rand) []int {
	A := make([]int, n)
	mark := make([]bool, n)
	for i := 0; i < n; i++ {
		A[i] = i
		mark[i] = false
	}

	i := n - 1
	u := n

	for u >= 2 {
		if !mark[i] {
			var j int
			for {
				j = rng.Intn(i)
				if !mark[j] {
					break
				}
			}
			A[i], A[j] = A[j], A[i]

			p := rng.Float64()
			threshold := float64(u-1) * Subfactorial(u-2) / Subfactorial(u)
			if p < threshold {
				mark[j] = true
				u--
			}
			u--
		}
		i--
	}

	return A
}
