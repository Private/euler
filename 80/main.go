//
// Kristoffer Langeland Knudsen
// July 2016
//

// Note that I do these exercises to LEARN Go; don't judge.

package main

/*
The exercise:

It is well known that if the square root of a natural number is not an integer, then it is irrational. The decimal expansion of such square roots is infinite without any repeating pattern at all.

The square root of two is 1.41421356237309504880..., and the digital sum of the first one hundred decimal digits is 475.

For the first one hundred natural numbers, find the total of the digital sums of the first one hundred decimal digits for all the irrational square roots.

Related reading:
https://en.wikipedia.org/wiki/Methods_of_computing_square_roots
*/

import (
	"fmt"
	"math/big"
)

type Context struct {
	number uint
	groups chan uint
	digits chan uint

	remainder *big.Int
	p         *big.Int
}

func NewContext(number uint) *Context {

	context := Context{
		number: number,

		groups: make(chan uint),
		digits: make(chan uint),

		remainder: new(big.Int),
		p:         new(big.Int),
	}

	go context.makeGroups()
	go context.makeDigits()

	return &context
}

func makeGroupsRecursive(number uint) []uint {

	if number < 100 {
		return []uint{number}
	}

	return append(makeGroupsRecursive(number/100), number%100)
}

func (context *Context) makeGroups() {

	for _, group := range makeGroupsRecursive(context.number) {
		context.groups <- group
	}

	for {
		context.groups <- 0
	}
}

func (context *Context) makeDigits() {

	for {
		p := context.p

		group := new(big.Int)
		group.SetUint64(uint64(<-context.groups))

		c := new(big.Int)
		c = c.Mul(context.remainder, big.NewInt(100))
		c = c.Add(c, group)

		var digit uint = 9

		x := new(big.Int)

		for {
			x.SetUint64(uint64(digit))

			v := big.NewInt(20)
			v = v.Mul(v, p)
			v = v.Add(v, x)
			v = v.Mul(v, x)

			test := new(big.Int)
			test = test.Sub(v, c)

			if test.Sign() > 0 {
				digit--
			} else {
				break
			}
		}

		// fmt.Printf("{c: %d} {p: %d} Digit: %d\n", c, p, digit)

		y := big.NewInt(20)
		y = y.Mul(y, p)
		y = y.Add(y, x)
		y = y.Mul(y, x)

		context.remainder = context.remainder.Sub(c, y)
		context.p = context.p.Mul(context.p, big.NewInt(10))
		context.p = context.p.Add(context.p, x)

		context.digits <- digit
	}

}

func (context *Context) SumDigits(n uint) uint {

	var sum, i uint = 0, 0

	for ; i < n; i++ {
		sum += <-context.digits
	}

	return sum
}

func (context *Context) ListDigits(n uint) []uint {

	if n == 0 {
		return []uint{}
	}

	return append([]uint{<-context.digits}, context.ListDigits(n-1)...)

}

func main() {
	fmt.Println("Project Euler")
	fmt.Println("Problem 80 - Square root digital expansion")
	fmt.Println("")

	var total, i uint = 0, 1

	for i = 1; i <= 100; i++ {

		context := NewContext(i)
		total += context.SumDigits(100)

	}

	// We've summed up not only the irrational digits, but also the handful of
	// rational roots in the mix. That's fine, we just sub them out.

	fmt.Printf("Digit total: %d\n", total)
}

// 1.41421356237309504880168872420969807856967187537694807317667973799
// 1.4142135623730950488016887242096980785696718753769
