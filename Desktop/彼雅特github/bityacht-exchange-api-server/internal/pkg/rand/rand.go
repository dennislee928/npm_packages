package rand

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"math/big"

	"golang.org/x/exp/constraints"
)

const onlyNumber = "0123456789"
const onlyNumberLength = len(onlyNumber)

const letterAndNumber = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const letterAndNumberLength = len(letterAndNumber)

// LetterAndNumberString returns, as a string, random string in the char set [0-9a-zA-Z] from crypto/rand. It panics if n < 0.
func LetterAndNumberString(n int) string {
	return randString(letterAndNumber, letterAndNumberLength, n)
}

// NumberString returns, as a string, random string in the char set [0-9] from crypto/rand. It panics if n < 0.
func NumberString(n int) string {
	return randString(onlyNumber, onlyNumberLength, n)
}

func randString(charSet string, setSize int, n int) string {
	if n < 0 {
		panic("invalid argument to Intn")
	} else if n == 0 {
		return ""
	}

	output := make([]byte, n)
	if _, err := rand.Read(output); err != nil {
		panic(err)
	}

	for i, randByte := range output {
		output[i] = charSet[int(randByte)%setSize]
	}

	return string(output)
}

// Intn returns, as a T, a non-negative random number in the half-open interval [0,n) from crypto/rand. It panics if n <= 0.
func Intn[T constraints.Signed](n T) T {
	if n < 0 {
		panic("invalid argument to Intn")
	}

	if output, err := rand.Int(rand.Reader, big.NewInt(int64(n))); err != nil {
		panic(err)
	} else {
		return T(output.Int64())
	}
}

// Float64 returns, as a float64, a pseudo-random number in the half-open interval [0.0,1.0)
func Float64() float64 {
	var buf [8]byte
	// 7 bytes is enough for 53 bits of precision.
	if _, err := rand.Read(buf[:7]); err != nil {
		panic(err)
	}
	// Mask off the unwanted bits to keep only 53 bits for precision.
	u64 := binary.BigEndian.Uint64(buf[:]) & ((1 << 53) - 1)
	return float64(u64) / float64(1<<53)
}

// Base32String returns, as a string, random string in base32 from crypto/rand. It panics if n < 0.
func Base32String(n int) string {
	if n < 0 {
		panic("invalid argument to Intn")
	} else if n == 0 {
		return ""
	}

	bufLen := n / 8 * 5 // 5 byte -> 8 char
	if n%8 != 0 {
		bufLen += 5
	}

	buf := make([]byte, bufLen)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}

	return base32.StdEncoding.EncodeToString(buf)[:n]
}
