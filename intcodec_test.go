package intcodec

import (
	"testing"
)

// test numbers taken from $GOROOT/src/pkg/strconv/itoa_test.go
type test struct {
	num  int64
	base int
	str  string
}

var tests = []test{
	{0, 10, "0"},
	{1, 10, "1"},
	{-1, 10, "-1"},
	{12345678, 10, "12345678"},
	{-987654321, 10, "-987654321"},
	{1<<31 - 1, 10, "2147483647"},
	{-1<<31 + 1, 10, "-2147483647"},
	{1 << 31, 10, "2147483648"},
	{-1 << 31, 10, "-2147483648"},
	{1<<31 + 1, 10, "2147483649"},
	{-1<<31 - 1, 10, "-2147483649"},
	{1<<32 - 1, 10, "4294967295"},
	{-1<<32 + 1, 10, "-4294967295"},
	{1 << 32, 10, "4294967296"},
	{-1 << 32, 10, "-4294967296"},
	{1<<32 + 1, 10, "4294967297"},
	{-1<<32 - 1, 10, "-4294967297"},
	{1 << 50, 10, "1125899906842624"},
	{1<<63 - 1, 10, "9223372036854775807"},
	{-1<<63 + 1, 10, "-9223372036854775807"},
	{-1 << 63, 10, "-9223372036854775808"},

	{0, 2, "0"},
	{10, 2, "1010"},
	{-1, 2, "-1"},
	{1 << 15, 2, "1000000000000000"},

	{-8, 8, "-10"},
	{057635436545, 8, "57635436545"},
	{1 << 24, 8, "100000000"},

	{16, 16, "10"},
	{-0x123456789abcdef, 16, "-123456789abcdef"},
	{1<<63 - 1, 16, "7fffffffffffffff"},
	{1<<63 - 1, 2, "111111111111111111111111111111111111111111111111111111111111111"},

	{16, 17, "g"},
	{25, 25, "10"},
	{(((((17*35+24)*35+21)*35+34)*35+12)*35+24)*35 + 32, 35, "holycow"},
	{(((((17*36+24)*36+21)*36+34)*36+12)*36+24)*36 + 32, 36, "holycow"},
}

func TestEncode(t *testing.T) {
	for _, test := range tests {
		s := LowerBaseN(test.base).EncodeInt(test.num)
		if s != test.str {
			t.Errorf("LowerBaseN(%d).Encode(%d) = %q, expected %q", test.base, test.num, s, test.str)
		}
	}
}

func TestDecode(t *testing.T) {
	for _, test := range tests {
		n, err := LowerBaseN(test.base).DecodeInt(test.str)
		if err != nil {
			t.Errorf("LowerBaseN(%d).DecodeInt(%q) error: %v", test.base, test.str, err)
		} else if n != test.num {
			t.Errorf("LowerBaseN(%d).DecodeInt(%q) = %d, expected %d", test.base, test.str, n, test.num)
		}
	}
}
