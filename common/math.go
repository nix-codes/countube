package common

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
	if N < 2 {
		return 0
	}

	// sieve[i] == true means i is prime
	sieve := make([]bool, N+1)
	for i := 2; i <= N; i++ {
		sieve[i] = true
	}

	for i := 2; i*i <= N; i++ {
		if sieve[i] {
			for j := i * i; j <= N; j += i {
				sieve[j] = false
			}
		}
	}

	// count primes
	count := 0
	for i := 2; i <= N; i++ {
		if sieve[i] {
			count++
		}
	}

	return count
}
