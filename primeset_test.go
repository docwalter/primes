package primes

import (
	"fmt"
	"testing"
	"time"
)

func TestInitialization(t *testing.T) {
	start := time.Now()
	size := uint64(100000000)
	p := NewPrimeSet(size)
	elapsed := time.Since(start)
	t.Logf("primes up to %d sieved in %s using %d kB", size, elapsed, p.MemoryUsage()>>10)
}

func TestIterators(t *testing.T) {
	set := NewPrimeSet(1000000)
	it := set.Iterator(0)
	count := 0
	max := uint64(0)
	for p, found := it.Next(); found; p, found = it.Next() {
		max = p
		count++
	}
	if max != 1000121 || count != 78506 {
		t.Error("prime sieve incorrect")
	}
	p, found := set.Iterator(1000121).Next()
	if !found || p != 1000121 {
		t.Error("iterator should have returned the last prime in the set")
	}
	p, found = set.Iterator(1000122).Next()
	if found {
		t.Errorf("iterator returned %d beyond the end of the set", p)
	}
}

func BenchmarkInitialization(b *testing.B) {
	NewPrimeSet(uint64(b.N))
}

func ExampleIterator() {
	set := NewPrimeSet(100) // creates a set reaching at least up to 100
	fmt.Println("largest number in set: ", set.LargestNumber())
	fmt.Println("largest prime number in set: ", set.LargestPrime())
	fmt.Println("all prime numbers:")
	it := set.Iterator(0) // creates an iterator starting with the first prime number
	p, ok := it.Next()
	for ok {
		fmt.Print(p, " ")
		p, ok = it.Next()
	}
	// Output:
	// largest number in set:  191
	// largest prime number in set:  191
	// all prime numbers:
	// 2 3 5 7 11 13 17 19 23 29 31 37 41 43 47 53 59 61 67 71 73 79 83 89 97 101 103 107 109 113 127 131 137 139 149 151 157 163 167 173 179 181 191
}
