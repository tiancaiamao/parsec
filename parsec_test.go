package parsec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestISODate(t *testing.T) {
	var tmp [10]Byte
	for i := '0'; i < '9'; i++ {
		tmp[i-'0'] = Byte(i)
	}
	digits := OR(tmp[0], tmp[1], tmp[2], tmp[3], tmp[4], tmp[5], tmp[6], tmp[7], tmp[8], tmp[9])
	year := SEQ(digits, digits, digits, digits)
	month := SEQ(digits, digits)
	day := SEQ(digits, digits)
	xx := Byte('-')
	date := SEQ(year, xx, month, xx, day)

	input := "2021-10-23"
	succ, remain := date.Parse(input)
	require.True(t, succ)
	require.Equal(t, remain, "")
}

type BenchResult struct {
	OP    Number
	Byte  Number
	Alloc Number
}

type Number int

func (n *Number) Parse(input string) (bool, string) {
	var res int
	remain := input
	for len(remain) > 0 {
		s, tmp := remain[0], remain[1:]
		if s < '0' || s > '9' {
			break
		}
		res = res * 10
		res += int(s - '0')
		remain = tmp
	}
	if res > 0 {
		*n = Number(res)
		return true, remain
	}
	return false, input
}

type Funcname struct{}

func (Funcname) Parse(input string) (bool, string) {
	azAZ := OR(Range('a', 'z'), Range('A', 'Z'))
	pattern := SEQ(azAZ, NOT(azAZ))
	succ, remain := pattern.Parse(input)
	if succ {
		return true, remain
	}

	xx := SEQ(azAZ, Funcname{})
	return xx.Parse(input)
}

func (r *BenchResult) Parse(input string) (bool, string) {
	WS := WhiteSpace{}
	Benchmark := Literal("Benchmark")
	Funcname := Funcname{}
	NSOP := Literal("ns/op")
	BytesOP := Literal("bytes/op")
	AllocsOP := Literal("allocs/op")
	pattern := SEQ(Benchmark, WS, Funcname,
		WS, &r.OP, WS, NSOP,
		WS, &r.Byte, WS, BytesOP,
		WS, &r.Alloc, WS, AllocsOP)
	return pattern.Parse(input)
}

func TestGoBenchLog(t *testing.T) {
	input := "Benchmark FuncnameXXX 23132 ns/op     823032 bytes/op    41 allocs/op"
	var r BenchResult
	succ, remain := r.Parse(input)
	require.Equal(t, remain, "")
	require.True(t, succ)
	require.Equal(t, r, BenchResult{
		OP: 23132,
		Byte: 823032,
		Alloc: 41,
	})
}

// func TestExpr() {
// 	// TODO
// 	input := "((3 - 5) + (4 * 7)) / 2"
// }

// func TestSexp() {
// 	// TODO
// 	input := "(/ (+ (- 3 5) (* 4 7))
// 		     4)"
// }
