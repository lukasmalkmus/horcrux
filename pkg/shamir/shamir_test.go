package shamir

import (
	"reflect"
	"testing"
)

func TestSplit_invalid(t *testing.T) {
	secret := []byte("test")

	_, err := Split(secret, 0, 0)
	assert(t, err != nil, "Expected error")

	_, err = Split(secret, 2, 3)
	assert(t, err != nil, "Expected error")

	_, err = Split(secret, 1000, 3)
	assert(t, err != nil, "Expected error")

	_, err = Split(secret, 10, 1)
	assert(t, err != nil, "Expected error")

	_, err = Split(nil, 3, 2)
	assert(t, err != nil, "Expected error")
}

func TestSplit(t *testing.T) {
	secret := []byte("test")

	out, err := Split(secret, 5, 3)
	ok(t, err)

	equals(t, len(out), 5)

	for _, share := range out {
		equals(t, len(share), len(secret)+1)
	}
}

func TestCombine_invalid(t *testing.T) {
	// Not enough parts
	_, err := Combine(nil)
	assert(t, err != nil, "Expected error")

	// Mis-match in length
	parts := [][]byte{
		[]byte("foo"),
		[]byte("ba"),
	}
	_, err = Combine(parts)
	assert(t, err != nil, "Expected error")

	//Too short
	parts = [][]byte{
		[]byte("f"),
		[]byte("b"),
	}
	_, err = Combine(parts)
	assert(t, err != nil, "Expected error")

	parts = [][]byte{
		[]byte("foo"),
		[]byte("foo"),
	}
	_, err = Combine(parts)
	assert(t, err != nil, "Expected error")
}

func TestCombine(t *testing.T) {
	secret := []byte("test")

	out, err := Split(secret, 5, 3)
	ok(t, err)

	// There is 5*4*3 possible choices,
	// we will just brute force try them all
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if j == i {
				continue
			}
			for k := 0; k < 5; k++ {
				if k == i || k == j {
					continue
				}
				parts := [][]byte{out[i], out[j], out[k]}
				recomb, err := Combine(parts)
				ok(t, err)

				equals(t, recomb, secret)
			}
		}
	}
}

func TestField_Add(t *testing.T) {
	out := add(16, 16)
	equals(t, out, uint8(0))

	out = add(3, 4)
	equals(t, out, uint8(7))
}

func TestField_Mult(t *testing.T) {
	out := mult(3, 7)
	equals(t, out, uint8(9))

	out = mult(3, 0)
	equals(t, out, uint8(0))

	out = mult(0, 3)
	equals(t, out, uint8(0))
}

func TestField_Divide(t *testing.T) {
	out := div(0, 7)
	equals(t, out, uint8(0))

	out = div(3, 3)
	equals(t, out, uint8(1))

	out = div(6, 3)
	equals(t, out, uint8(2))
}

func TestPolynomial_Random(t *testing.T) {
	p, err := makePolynomial(42, 2)
	ok(t, err)

	if p.coefficients[0] != 42 {
		t.Fatalf("bad: %v", p.coefficients)
	}
}

func TestPolynomial_Eval(t *testing.T) {
	p, err := makePolynomial(42, 1)
	ok(t, err)

	out := p.evaluate(0)
	equals(t, out, uint8(42))

	out = p.evaluate(1)
	exp := add(42, mult(1, p.coefficients[1]))
	if out != exp {
		t.Fatalf("bad: %v %v %v", out, exp, p.coefficients)
	}
}

func TestInterpolate_Rand(t *testing.T) {
	for i := 0; i < 256; i++ {
		p, err := makePolynomial(uint8(i), 2)
		ok(t, err)

		xVals := []uint8{1, 2, 3}
		yVals := []uint8{p.evaluate(1), p.evaluate(2), p.evaluate(3)}
		out := interpolatePolynomial(xVals, yVals, 0)
		if out != uint8(i) {
			t.Fatalf("Bad: %v %d", out, i)
		}
	}
}
func BenchmarkSplit1MB5Parts5Threshold(b *testing.B) { benchmarkSplit1MB(b, 5, 5) }
func BenchmarkSplit1MB5Parts4Threshold(b *testing.B) { benchmarkSplit1MB(b, 5, 4) }
func BenchmarkSplit1MB5Parts3Threshold(b *testing.B) { benchmarkSplit1MB(b, 5, 3) }
func BenchmarkSplit1MB5Parts2Threshold(b *testing.B) { benchmarkSplit1MB(b, 5, 2) }

func benchmarkSplit1MB(b *testing.B, parts, threshold int) {
	secret := make([]byte, 1024*1024)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := Split(secret, 5, 3)
		ok(b, err)
	}
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	tb.Helper()
	if !condition {
		tb.Fatalf("\033[31m "+msg+"\033[39m\n\n", v...)
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("\033[31m unexpected error: %s\033[39m\n\n", err.Error())
	}
}

// equals fails the test if got is not equal to want.
func equals(tb testing.TB, got, want interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(got, want) {
		tb.Fatalf("\033[31m\n\n\tgot: %#v\n\n\twant: %#v\033[39m\n\n", got, want)
	}
}
