/*

Package primes provides functionality for prime numbers.

Prime numbers are precalculated and efficiently stored in a Set: on a MacBook Pro, all primes up to 100,000,000
are calculated in 600ms and use 4M memory, i.e. 6ns and 1/3 bit per number. All methods are accessible through the
Set interface, an implementation of which is obtained using NewPrimeSet().

Example usage:

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

*/
package primes

// Set is a set of prime numbers.
type Set interface {
	IsPrime(n uint64) bool          // true iff n is prime
	Iterator(start uint64) Iterator // allows for traversing the set
	LargestNumber() uint64          // largest number in the set
	LargestPrime() uint64           // largest prime number in the set
	MemoryUsage() uint              // number of bytes used for the prime bits
}

// set is the internal implementation of Set.
type set struct {
	bits         []uint64 // bits for prime number candidates that are not divisible by 2 and 3
	largestPrime uint64   // largest prime number in the set
}

// NewPrimeSet creates a new set of prime numbers up to a given limit.
func NewPrimeSet(limit uint64) Set {
	if limit < 5 {
		panic("prime set must have at least a size of 5")
	}
	s := new(set)
	i := numberToIndex(limit)
	words := i >> 6
	if i&63 != 0 {
		words++
	}
	s.bits = make([]uint64, words)
	calculatePrimeBitSet(s.bits)
	h, _ := highestSetBit(s.bits)
	s.largestPrime = indexToNumber(h)
	return s
}

// Iterator returns an iterator over the prime set that returns all primes in ascending order.
func (s *set) Iterator(start uint64) Iterator {
	if start <= 2 {
		// starting from the beginning, so the next prime number is 2
		return &iterator{s, 0, 2}
	}
	i := numberToIndex(start)
	for {
		n, found := nextSetBit(s.bits, i)
		if !found {
			// there is no next prime number in the set, so return an iterator that is already finished
			return &iterator{s, 0, 0}
		}
		p := indexToNumber(n)
		if p >= start {
			return &iterator{s, n, p}
		}
		i++
	}
}

// IsPrime returns true iff n is a prime number.
func (s *set) IsPrime(n uint64) bool {
	if n <= 63 {
		const quickcheck = uint64(0x816d129a64b4cb6f)
		quickCheckMask := uint64(1) << ((n - 1) >> 1)
		return quickcheck&quickCheckMask != 0
	}
	if n&1 == 0 || n%3 == 0 {
		return false
	}
	i := numberToIndex(n)
	if i > 0 {
		return getBit(s.bits, i)
	}
	return false
}

// LargestPrime returns the largest prime number in the set, i.e. the upper limit for IsPrime() etc.
func (s *set) LargestPrime() uint64 {
	return s.largestPrime
}

// LargestNumber returns the largest prime number in the set, i.e. the upper limit for IsPrime() etc.
func (s *set) LargestNumber() uint64 {
	return indexToNumber(uint(len(s.bits)<<6 - 1))
}

func (s *set) MemoryUsage() uint {
	return uint(len(s.bits) << 3)
}

// calculatePrimeBitSet initializes the prime bit set using a simple prime sieve.
func calculatePrimeBitSet(bits []uint64) {
	setAllBits(bits)
	highestbitindex := uint(len(bits)<<6 - 1)
	limit := indexToNumber(highestbitindex)
	for i, found := uint(0), true; found; i, found = nextSetBit(bits, i+1) {
		p := indexToNumber(i)
		n := p * 5
		d := p * 2
		for n <= limit {
			if n%3 != 0 {
				clearBit(bits, numberToIndex(n))
			}
			n += d
		}
	}
}

// numberToIndex returns for a given n the index of the bit in the bit set which marks primality of n.
func numberToIndex(n uint64) uint {
	if n < 5 {
		return 0
	}
	x2 := n / 2
	x3 := n / 3
	x23 := n / 6
	return uint(n - (x2 + x3 - x23) - 1)
}

// indexToNumber determines which prime candidate is associated with a given index.
func indexToNumber(i uint) uint64 {
	if i == 0 {
		return 3
	}
	n := uint64(5)
	n += uint64(6 * ((i - 1) >> 1))
	if i&1 == 0 {
		n += 2
	}
	return n
}
