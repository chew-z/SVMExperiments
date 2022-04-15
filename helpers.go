package main

import (
	"math/rand"
	"time"
)

var userAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.60 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3988.121 Safari/537.35",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36",
}

func randUserAgent() string {
	i := intN(len(userAgents) - 1)
	return userAgents[i]
}

func isEqualInt64(a int64, b int64, delta int64) bool {
	if absDiffInt64(a, b) <= delta {
		return true
	} else {
		return false
	}
}

func absDiffInt64(x, y int64) int64 {
	if x < y {
		return y - x
	}
	return x - y
}

func EagleOrTail() int {
	i := rand.Intn(100)
	if i <= 50 {
		return -1
	} else {
		return 1
	}
}

// Scale 0..1 to 0..100
func Scale(x float64) int {
	return int(100.0 * x)
}

// Normalizes values between min and max
func Normalize(val float64, min float64, max float64) float64 {
	delta := max - min
	return (val - min) / delta
}

func MinMax(array []float64) (float64, float64) {
	max := array[0]
	min := array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

// Generates a pseudo-random int, where 0 <= x < `n`.
func intN(n int) int {
	seed := rand.NewSource(time.Now().UnixNano())
	rnew := rand.New(seed)
	return rnew.Intn(n)
}
