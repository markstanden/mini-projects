package argonhasher

import (
	"fmt"
	"testing"
)

/***
*
*	Benchmarks
*
***/

var result string

func benchmarkEncode(cost int, b *testing.B) {
	var hash string
	for n := 0; n < b.N; n++ {
		hash = Encode("secretcode", uint(cost))
		fmt.Println(hash)
	}
	result = hash
}

func BenchmarkEncodeCost1(b *testing.B)  { benchmarkEncode(1, b) }
func BenchmarkEncodeCost2(b *testing.B)  { benchmarkEncode(2, b) }
func BenchmarkEncodeCost3(b *testing.B)  { benchmarkEncode(3, b) }
func BenchmarkEncodeCost4(b *testing.B)  { benchmarkEncode(4, b) }
func BenchmarkEncodeCost5(b *testing.B)  { benchmarkEncode(5, b) }
func BenchmarkEncodeCost6(b *testing.B)  { benchmarkEncode(6, b) }
func BenchmarkEncodeCost7(b *testing.B)  { benchmarkEncode(7, b) }
func BenchmarkEncodeCost8(b *testing.B)  { benchmarkEncode(8, b) }
func BenchmarkEncodeCost9(b *testing.B)  { benchmarkEncode(9, b) }
func BenchmarkEncodeCost10(b *testing.B) { benchmarkEncode(10, b) }

/***
*
*	Examples
*
***/

func ExampleEncode() {
	/*
		the plain text password to hash with argon2
	*/
	pw := "c2BDNoW38DStXvzP"

	/*
		the cost (difficulty) of the computation to produce the hash.
		if 0, the Encode function chooses a sensible, secure default
	*/
	var cost uint = 0

	/*
		Produce the hash.  Each hash generates a unique salt so no two hashes should be equal.
		The hash, salt, and encoding options are all stored in the string,
	*/
	hashedPassword1 := Encode(pw, cost)
	hashedPassword2 := Encode(pw, cost)
	fmt.Println(hashedPassword1 == hashedPassword2)
	// Output:
	// false
}
