package primes

// Iterator allows for traversing a prime set in ascending order.
type Iterator interface {
	Next() (uint64, bool) // next prime number and status
}

// Internal implementation of Iterator.
type iterator struct {
	set       *set   // prime set that is traversed by this iterator
	nextIndex uint   // bit index that will be used for the next Next() call
	nextPrime uint64 // prime number that will be returned by the next Next() call or 0 after the end of the sequence
}

// Next returns the next prime number in the set in ascending order.
func (i *iterator) Next() (uint64, bool) {
	if i.nextPrime == 0 {
		// end of sequence reached
		return 0, false
	}
	if i.nextPrime == 2 {
		// start of sequence at the first prime number
		i.nextPrime = 3
		i.nextIndex = 0
		return 2, true
	}
	r := i.nextPrime
	n, found := nextSetBit(i.set.bits, i.nextIndex+1)
	i.nextIndex = n
	if found {
		i.nextPrime = indexToNumber(n)
	} else {
		i.nextPrime = 0
	}
	return r, true
}
