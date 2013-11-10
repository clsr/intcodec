// Packade intcodec implements encoding and decoding integers in various bases.
package intcodec

import (
	"errors"
	"math"
)

// Codec implements encodin and decoding for a specific base and character set.
type Codec struct {
	encodeMap string
	decodeMap []int16
	maxlen    int
}

const (
	lowerBase36 = "0123456789abcdefghijklmnopqrstuvwxyz"
	upperBase36 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base64      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

// ErrEmpty means that the string given to a Decode method was empty.
var ErrEmpty = errors.New("intcodec: decoding empty string")

// ErrInvalid means that the string being decoded contains characters not in
// the codec's character set.
var ErrInvalid = errors.New("intcodec: invalid character in string")

// ErrOverflow means that the string given to a Decode method overflows its
// return type.
var ErrOverflow = errors.New("intcodec: integer overflow")

// Base64 is a Codec using the standard base64 charset without padding.
var Base64 = New(base64)

// New creates a new Codec from the given character set.
//
// The encoding argument must be a string of 2 to 256 (inclusive) unique
// characters. The base of the codec's encoding is the length of that string.
//
// When encoding, each byte in the string represents a value matching its
// index. That means that, when decoding, encoding[i] represents the digit i
// and the digit i encodes into encoding[i] when encoded.
func New(encoding string) *Codec {
	if len(encoding) < 2 {
		panic("intcodec: minimum permitted base is 2")
	}
	if len(encoding) > 256 {
		panic("intcodec: maximum permitted base is 256")
	}

	decode := make([]int16, 256)
	for i := range decode {
		decode[i] = -1
	}
	for i := 0; i < len(encoding); i++ {
		if decode[encoding[i]] >= 0 {
			panic("intcodec: duplicate character in encoding")
		}
		decode[encoding[i]] = int16(i)
	}

	// maximum number of digits and the sign
	maxlen := int(math.Log(float64(2<<63))/math.Log(float64(len(encoding)))+0.5) + 1

	return &Codec{encoding, decode, maxlen}
}

// LowerBaseN creates a Codec for base n with the standard numeric-letter
// character set in lower case ("0123456789abcdef..."). The argument n must be
// between 2 and 36 (inclusive).
func LowerBaseN(n int) *Codec {
	if n < 2 || n > 36 {
		panic("intcodec: LowerBaseN argument must be >2 and <36")
	}
	return New(lowerBase36[:n])
}

// LowerBaseN creates a Codec for base n with the standard numeric-letter
// character set in upper case ("0123456789ABCDEF..."). The argument n must be
// between 2 and 36 (inclusive).
func UpperBaseN(n int) *Codec {
	if n < 2 || n > 36 {
		panic("intcodec: UpperBaseN argument must be >2 and <36")
	}
	return New(upperBase36[:n])
}

// Base returns the radix of the codec.
func (c *Codec) Base() int {
	return len(c.encodeMap)
}

// EncodeUint encodes a uint64 into a string using the codec's charset.
func (c *Codec) EncodeUint(n uint64) string {
	return c.encode(n, false)
}

// EncodeInt encodes an int64 into a string using the codec's charset. If the
// integer is negative, the produced string is prefixed with '-'. This may
// collide with a '-' character from the codec's charset, if it is present.
func (c *Codec) EncodeInt(n int64) string {
	if n < 0 {
		return c.encode(uint64(-n), true)
	}
	return c.encode(uint64(n), false)
}

func (c *Codec) encode(n uint64, sign bool) string {
	buf := make([]byte, c.maxlen)
	i := len(buf)

	for {
		i--
		buf[i] = c.encodeMap[n%uint64(len(c.encodeMap))]
		n /= uint64(len(c.encodeMap))
		if n == 0 {
			break
		}
	}

	if sign {
		i--
		buf[i] = '-'
	}

	return string(buf[i:])
}

// DecodeUint decodes a string into a uint64 using the codec's charset.
// If the string is empty, contains characters outside of the charset or
// overflows a uint64, an error is returned.
func (c *Codec) DecodeUint(s string) (n uint64, err error) {
	if s == "" {
		err = ErrEmpty
		return
	}
	for i := 0; i < len(s); i++ {
		d := c.decodeMap[s[i]]
		if d < 0 {
			err = ErrInvalid
			return
		}
		n2 := n*uint64(len(c.encodeMap)) + uint64(d)
		if n2 < n {
			err = ErrOverflow
			return
		}
		n = n2
	}
	return
}

// DecodeInt decodes a string into a uint64 using the codec's charset.
// If the string is empty, contains characters outside of the charset or
// overflows a uint64, an error is returned. If the string begins with '-',
// the integer is decoded as negative. If the charset contains the '-'
// character, this will mis-decode integers where that is the first digit.
func (c *Codec) DecodeInt(s string) (n int64, err error) {
	if s == "" {
		err = ErrEmpty
		return
	}
	sign := false
	if s[0] == '-' {
		s = s[1:]
		sign = true
	}
	u, err := c.DecodeUint(s)
	n = int64(u)
	if err != nil && n < 0 {
		err = ErrOverflow
	}
	if sign {
		n = -n
	}
	return
}
