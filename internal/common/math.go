package common

import (
	"math/rand"
)

func Chance(p float64) bool {
	return rand.Float64() < p
}

func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func CountPrimesBelow(N int) int {
	if N <= 2 {
		return 0
	}

	sieve := make([]bool, N)
	for i := 2; i < N; i++ {
		sieve[i] = true
	}

	for i := 2; i*i < N; i++ {
		if sieve[i] {
			for j := i * i; j < N; j += i {
				sieve[j] = false
			}
		}
	}

	count := 0
	for i := 2; i < N; i++ {
		if sieve[i] {
			count++
		}
	}
	return count
}

/*
lo: lower limit (inclusive)
hi: upper limit (exclusive)
*/
func CountPrimesInRange(lo, hi int) int {
	if hi <= 2 || lo >= hi {
		return 0
	}

	sieve := make([]bool, hi)
	for i := 2; i < hi; i++ {
		sieve[i] = true
	}

	for i := 2; i*i < hi; i++ {
		if sieve[i] {
			for j := i * i; j < hi; j += i {
				sieve[j] = false
			}
		}
	}

	if lo < 2 {
		lo = 2
	}
	count := 0
	for i := lo; i < hi; i++ {
		if sieve[i] {
			count++
		}
	}
	return count
}
