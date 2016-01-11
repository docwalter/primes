package primes

// Factorizer holds precalculated factors for a given range of numbers, thus allowing for their factorization in near-constant time.
type Factorizer interface {
	LargestFactorOf(n uint64) (uint64, bool) // largest prime factor of a given number
}

// Internal implementation of Factorizer.
type factorizer struct {
	set           *set     // underlying prime set
	factors       []uint64 // largest prime factors of all numbers not divisible by 2 or 3
	largestNumber uint64   // largest number that can be factorized by this Factorizer
}

// Factorizer returns a new factorizer for numbers in the range up to n.
func (s *set) Factorizer(max uint64) Factorizer {
	return newFactorizerBuilder(s, max).build()
}

// LargestFactorOf returns the largest prime factor of a given number.
// If the factorizer boundaries are exceeded during search, the second result is false.
func (f *factorizer) LargestFactorOf(n uint64) (uint64, bool) {
	n >>= numberOfTrailingZeroes(n)
	if n == 0 {
		return 0, false
	}
	if n == 1 {
		return 2, true
	}
	for n%3 == 0 {
		n /= 3
	}
	if n == 1 {
		return 3, true
	}
	if n > f.largestNumber {
		return 0, false
	}
	i := numberToIndex(n)
	return f.factors[i], true
}

// factorizerBuilder is a temporary structure which creates a factorizer and precalculates its factors.
type factorizerBuilder struct {
	set      *set // underlying prime set
	factors  []uint64
	max      uint64
	stack    []uint64
	sp       int
	maxDepth int
}

const maxuint = uint64(0xffffffffffffffff) // maximum value of an uint64

// newFactorizerBuilder creates a new factorizerBuilder with empty factors.
func newFactorizerBuilder(set *set, max uint64) *factorizerBuilder {

	// create empty factors array
	factors := make([]uint64, numberToIndex(max)+1)

	// determine maximum recursion depth and initialize recursion stack
	maxDepth := 0
	test := uint64(1)
	it := set.Iterator(5)
	p, ok := it.Next()
	for ok && test < max && test <= maxuint/p {
		maxDepth++
		test *= p
		p, ok = it.Next()
	}
	stack := make([]uint64, maxDepth+1)

	return &factorizerBuilder{set, factors, max, stack, 1, maxDepth}
}

// build precalculates the factors in the factorizerBuilder.
func (b *factorizerBuilder) build() *factorizer {
	// fill factors recursively
	it := b.set.Iterator(5)
	p, ok := it.Next()
	for ok && p <= b.max {
		b.factors[numberToIndex(p)] = p // the prime number has itself as the only (and thus the largest) prime factor
		if p < b.max/2 {
			b.stack[0] = p
			b.stack[1] = 5
			b.initRecursively(p)
		}
		p, ok = it.Next()
	}

	// build and return the factorizer
	return &factorizer{b.set, b.factors, b.max}
}

/*
Initializes a part of {@link #factors} by recursively setting the largest prime factor of all multiples of {@code base}.

Simple example with base=18797

Prime factors 2 and 3 are treated transparently by skipping all multiples of 2 and 3 in factors, so we start with prime factor 5.

Recursion level 1: prime=5 => n1, 5 x n1, 5^2 x n1, 5^3 x n1, ...
Recursion level 2: prime=7 => n2, 7 x n2, 7^2 x n2, 7^3 x n2, ...

At every level, recursion is stopped as soon as the result of the next multiplication would exceed max. Nonetheless, this
approach quickly leads to a stack overflow, since every available prime number leads to another recursion. That's why recursion must
be truncated carefully so that there are no missing results in the end.

Determining the longest possible chain of prime factors

Considering that max is the largest number that should be calculated, the longest chain of distinct prime factors that can
appear is 3 x 5 x 7 x 11 x 13 x 17 x ... < max. This number is also the maximum recursion level and is calculated in newFactorizerBuilder()
above.

Reducing recursion levels

In the simple example, prime factors are skipped by pushing an unmodified base downward. This happens regularly, so the stack
size explodes until the first actual multiplication takes place. In the advanced example, prime factors are skipped by iterating over
them at the same recursion level. stack keeps track of the prime factors already used, so that upon entering a
new recursion level, prime factor iteration starts at the next higher prime number after
stack[sp] and stops, of course, when max will be exceeded.

Recursion invariant

- Upon entering initRecursively, stack[sp] is set to the start of prime factor iteration at this recursion level.
- Before entering another recursion level, the calling level is responsible for pushing a new starting prime factor on the stack. If the
  maximum stack size maxDepth is exceeded, recursion is stopped. If the first multiplication on the next level
  would already exceed max, recursion must not be performed.
- After the callee returned, the caller must pop the stack.
- stack[0] always contains the largest prime factor with which all multiples should be marked in factors</li>
*/
func (b *factorizerBuilder) initRecursively(base uint64) {

	p := b.stack[b.sp]
	it := b.set.Iterator(p)
	prime, ok := it.Next()

	for ok && base <= b.max/prime && prime <= b.stack[0] {

		// this is the following prime factor with which the next recursion level will start
		next, nextok := it.Next()
		if b.sp < b.maxDepth {
			b.stack[b.sp+1] = next
		}

		// iterate over powers of prime, starting with prime^1
		for i := base * prime; i <= b.max; i *= prime {

			// mark the current number
			b.factors[numberToIndex(i)] = b.stack[0]

			// stop iteration if there would be an integer overflow at the next recursion level
			if i > maxuint/next {
				break
			}

			// do not recurse if
			// 1. the first mutiplication at the next recursion level would already exceed max,
			// 2. the maximum recursion depth is already reached, or
			// 3. the nextPrime would be larger than the largest prime factor in stack[0].
			if i > b.max/next || b.sp == b.maxDepth || next > b.stack[0] {
				continue
			}

			// recurse
			b.sp++ // stack element is already set
			b.initRecursively(i)
			b.sp--

			// avoid integer overflow in the next iteration of prime powers
			if i > maxuint/prime {
				break
			}
		}

		// avoid integer overflow in the next iteration of prime factors
		if base > maxuint/next {
			break
		}

		prime = next
		ok = nextok
	}
}
