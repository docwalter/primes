# Prime numbers for Go

**This is work in progress, "unstable".** As I find time, I will enhance this package.

Package primes provides functionality for prime numbers.

Prime numbers are precalculated and efficiently stored in a Set: on a MacBook Pro, all primes up to 100,000,000
are calculated in 600ms and use 4M memory, i.e. 6ns and 1/3 bit per number. All methods are accessible through the
Set interface, an implementation of which is obtained using NewPrimeSet():

```go
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
```

The most common use case for prime numbers is factorization of numbers. I created the library mainly to solve some
http://projecteuler.net problems, where there is usually a range of numbers to be factorized. So there is a Factorizer which
can, after some precalculations, factorize numbers up to a given limit:

```go
max := uint64(200000000)
set := NewPrimeSet(max)
factorizer := set.Factorizer(max)
f, ok := factorizer.LargestFactorOf(123456)
```
