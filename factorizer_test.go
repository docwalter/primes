package primes

import (
	"log"
	"testing"
	"time"
)

func TestFactorizer(t *testing.T) {
	max := uint64(200000000)
	starttime := time.Now()
	set := NewPrimeSet(max)
	log.Println("prime set initialized after ", time.Since(starttime))
	f := set.Factorizer(max)
	log.Println("factorizer initialized after ", time.Since(starttime))
	it := set.Iterator(0)
	p, ok := it.Next()
	for ok && p <= max {
		if pf, ok := f.LargestFactorOf(p); !ok || pf != p {
			t.Errorf("prime %d incorrectly factorized with factor %d", p, pf)
		}
		p, ok = it.Next()
	}
	if _, fok := f.LargestFactorOf(0); fok {
		t.Error("0 should not have a prime factor")
	}
	testLargestFactor(t, f, 2, 2)
	testLargestFactor(t, f, 1024, 2)
	testLargestFactor(t, f, 3, 3)
	testLargestFactor(t, f, 5, 5)
	testLargestFactor(t, f, 210, 7)
	testLargestFactor(t, f, 37055, 7411)
	testLargestFactor(t, f, 23173, 23173)
	testLargestFactor(t, f, 1664099, 1291)
	testLargestFactor(t, f, 3750000, 5)
}

func testLargestFactor(t *testing.T, factorizer Factorizer, n, factor uint64) {
	f, ok := factorizer.LargestFactorOf(n)
	if !ok {
		t.Errorf("LargestFactorOf(%d) failed", n)
	} else if f != factor {
		t.Errorf("LargestFactorOf(%d) = %d instead of %d", n, f, factor)
	}
}
