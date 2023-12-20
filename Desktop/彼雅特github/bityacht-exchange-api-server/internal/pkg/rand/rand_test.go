package rand

import (
	"math"
	mathrand "math/rand"
	"testing"

	"gonum.org/v1/gonum/stat"
)

func TestLetterAndNumberString(t *testing.T) {
	for i := 0; i <= 256; i++ {
		randStr := LetterAndNumberString(i)
		if len(randStr) != i {
			t.Errorf("bad rand string, n = %d, output is %q", i, randStr)
			return
		}

		for index, r := range randStr {
			var ok bool

			for _, rInLetterAndNumber := range letterAndNumber {
				if r == rInLetterAndNumber {
					ok = true
					break
				}
			}
			if !ok {
				t.Errorf("bad rand string, n = %d, bad char %q at index %d, output is %q", i, string(r), index, randStr)
				return
			}
		}
	}
}

func TestIntn(t *testing.T) {
	intTest := []int{1, 5, 10, 100, 1000, 9223372036854775807}
	int8Test := []int8{1, 5, 10, 100, 127}
	int16Test := []int16{1, 5, 10, 100, 1000, 32767}
	int32Test := []int32{1, 5, 10, 100, 1000, 2147483647}
	int64Test := []int64{1, 5, 10, 100, 1000, 9223372036854775807}

	for _, v := range intTest {
		if randNum := Intn(v); randNum < 0 || randNum >= v {
			t.Errorf("bad rand int, type = int, n = %d, output is %d", v, randNum)
		}
	}
	for _, v := range int8Test {
		if randNum := Intn(v); randNum < 0 || randNum >= v {
			t.Errorf("bad rand int, type = int8, n = %d, output is %d", v, randNum)
		}
	}
	for _, v := range int16Test {
		if randNum := Intn(v); randNum < 0 || randNum >= v {
			t.Errorf("bad rand int, type = int16, n = %d, output is %d", v, randNum)
		}
	}
	for _, v := range int32Test {
		if randNum := Intn(v); randNum < 0 || randNum >= v {
			t.Errorf("bad rand int, type = int32, n = %d, output is %d", v, randNum)
		}
	}
	for _, v := range int64Test {
		if randNum := Intn(v); randNum < 0 || randNum >= v {
			t.Errorf("bad rand int, type = int64, n = %d, output is %d", v, randNum)
		}
	}
}

func TestFloat64(t *testing.T) {
	count := 10000000
	piece := 10
	biasCheck := make([]float64, piece)
	for i := 0; i < count; i++ {
		randFloat := Float64()
		if randFloat < 0.0 || randFloat >= 1.0 {
			t.Errorf("bad rand float, output is %f", randFloat)
			return
		}
		biasCheck[int(randFloat*float64(piece))]++
	}
	_, std := stat.MeanStdDev(biasCheck, nil)
	ci := std / math.Sqrt(float64(count))
	if ci > 0.5 {
		t.Errorf("confidence interval = %v is more thand %v", ci, 0.5)
		return
	}
}

func BenchmarkFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64()
	}
}

func BenchmarkMathFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mathrand.Float64() // #nosec G404
	}
}

func TestBase32String(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
	}{
		{"0", args{0}},
		{"1", args{1}},
		{"5", args{5}},
		{"10", args{10}},
		{"20", args{20}},
		{"1024", args{1024}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base32String(tt.args.n); len(got) != tt.args.n {
				t.Errorf("Base32String() = %q, len = %v, want %v", got, len(got), tt.args.n)
			} else {
				t.Log(got)
			}
		})
	}
}
