package primes

// wordMask splits a bit index into an uint64 array index and a bit mask.
func wordMask(i uint) (uint, uint64) {
	return i >> 6, 1 << (i & 63)
}

// setBit sets bit i in the given uint64 array.
func setBit(bits []uint64, i uint) {
	word, mask := wordMask(i)
	bits[word] |= mask
}

// setBit clears bit i in the given uint64 array.
func clearBit(bits []uint64, i uint) {
	word, mask := wordMask(i)
	bits[word] &^= mask
}

// getBit returns true iff bit i is set in the given uint64 array.
func getBit(bits []uint64, i uint) bool {
	word, mask := wordMask(i)
	return bits[word]&mask != 0
}

// setAllBit sets all bits in the given uint64 array.
func setAllBits(bits []uint64) {
	const allbits = uint64(0xffffffffffffffff)
	for i := range bits {
		bits[i] = allbits
	}
}

// nextSetBit returns the index of the next set bit in the given uint64 array, starting from i.
// If there is no bit set at or after index i, the second result is false.
func nextSetBit(bits []uint64, i uint) (uint, bool) {
	word := int(i >> 6)
	if word >= len(bits) {
		return 0, false
	}
	w := bits[word] >> (i & 63)
	if w != 0 {
		return i + numberOfTrailingZeroes(w), true
	}
	word++
	for word < len(bits) {
		if bits[word] != 0 {
			return uint(word)<<6 + numberOfTrailingZeroes(bits[word]), true
		}
		word++
	}
	return 0, false
}

// highestSetBit returns the index of the highest set bit.
// If there is no bit set, the second result is false.
func highestSetBit(bits []uint64) (uint, bool) {
	for word := len(bits) - 1; word >= 0; word-- {
		if bits[word] != 0 {
			return uint(word)<<6 + 63 - numberOfLeadingZeroes(bits[word]), true
		}
	}
	return 0, false
}

// numberOfLeadingZeroes returns the number of leading zero bits (0..64) in the given uint64.
func numberOfLeadingZeroes(i uint64) uint {
	if i == 0 {
		return 64
	}
	n := uint(0)
	x := i
	if x&0xffffffff00000000 == 0 {
		n += 32
		x <<= 32
	}
	if x&0xffff000000000000 == 0 {
		n += 16
		x <<= 16
	}
	if x&0xff00000000000000 == 0 {
		n += 8
		x <<= 8
	}
	if x&0xf000000000000000 == 0 {
		n += 4
		x <<= 4
	}
	if x&0xc000000000000000 == 0 {
		n += 2
		x <<= 2
	}
	if x&0x8000000000000000 == 0 {
		n += 1
		x <<= 1
	}
	return n
}

// numberOfTrailingZeroes returns the number of trailing zero bits (0..64) in the given uint64.
func numberOfTrailingZeroes(i uint64) uint {
	if i == 0 {
		return 64
	}
	n := uint(0)
	x := i
	if x&0xffffffff == 0 {
		n += 32
		x >>= 32
	}
	if x&0xffff == 0 {
		n += 16
		x >>= 16
	}
	if x&0xff == 0 {
		n += 8
		x >>= 8
	}
	if x&0xf == 0 {
		n += 4
		x >>= 4
	}
	if x&0x3 == 0 {
		n += 2
		x >>= 2
	}
	if x&0x1 == 0 {
		n += 1
		x >>= 1
	}
	return n
}
