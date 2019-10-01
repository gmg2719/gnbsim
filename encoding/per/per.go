// Package per is implementation for Basic Pckage Encoding Rule (PER) in
// ALIGNED variant.
package per

import (
	"fmt"
	"math/bits"
)

// EncConstrainedWholeNumber is the implementation for
// 10.5 Encoding of constrained whole number.
func EncConstrainedWholeNumber(input, min, max int) (v []uint8, bitlen int, err error) {

	if input < min || input > max {
		err = fmt.Errorf("EncConstrainedWholeNumber: input value=%d is out of range. (should be %d <= %d)", input, min, max)
		return
	}

	inputRange := max - min + 1
	inputEnc := input - min

	switch {
	case inputRange == 1: // empty bit-field
		return
	case inputRange < 256: // the bit-field case
		bitlen = bits.Len(uint(inputRange))
		v = append(v, uint8(inputEnc))
		return
	case inputRange == 256: // the one-octet case
		bitlen = 8
		v = append(v, uint8(inputEnc))
		return
	case inputRange <= 65536: // the two-octet case
		bitlen = 16
		v = append(v, uint8((inputEnc >> 8) & 0xff))
		v = append(v, uint8(inputEnc & 0xff))
		return
	case inputRange > 65537: // the indefinite length case
		// not implemented yet
		err = fmt.Errorf("EncConstrainedWholeNumber: not implemented yet for %d", input)
		return
	}
	err = fmt.Errorf("EncConstrainedWholeNumber: invalid range min=%d, max=%d", min, max)
	return
}

/*
// EncInteger returns multi-byte BIT STRING
// 12. Encoding the integer type
func EncInteger(input, min, max int, extmark bool) (v []uint8, err error) {

}

func encInteger2(input, min, max int) (v []uint8, err error) {
	if input < min || input > max {
		err = fmt.Errorf("EncInteger: input value is out of range.")
		return
	}
	inputRange := max - min + 1
	octLen := (bitLen-1)/8 + 1
	val = make([]uint8, octLen, octLen+1)

	switch {
	case inputRange <= 255:
		v, err = EncBitString(input, 0, 0)
	case inputRange == 256:
		v = make([]uint8, 1, 1)
		v[0] = uint8(input)
	case inputRange <= 65536:
		v = make([]uint8, 2, 2)
		v[1] = uint8(input & 0xff)
		input >>= 8
		v[0] = uint8(input & 0xff)
	case inputRange > 65537:
		err = fmt.Errorf("EncInteger: input value is out of range.")
	}
	return
}

// EncEnumerated return ENUMERATED preamble
// 12. Encoding the enumerated type
func EncEnumerated(input int) (val []uint8) {
	val = make([]uint8, 1, 1)
	val[0] = uint8(input)
	return
}

// EncBitString returns multi-byte BIT STRING
// 15. Encoding the bitstering type
func EncBitString(input, min, max int) (val []uint8, err error) {

	bitLen := bits.Len(uint(input))
	if bitLen > max {
		err = fmt.Errorf("BitString: input value overflow.")
		return
	}

	if bitLen < min || min == max {
		bitLen = min
	}

	octLen := (bitLen-1)/8 + 1
	val = make([]uint8, octLen, octLen+1)
	offset := 0
	if min != max {
		val[offset] = uint8(bitLen - min)
		val = append(val, 0)
		offset++
	}
	for bit := bitLen; bit > 0; bit -= 8 {
		tmp := offset + bit/8
		val[tmp] = uint8(input & 0xff)
		input >>= 8
	}
	return
}

// EncOctetString returns multi-byte OCTET STRING
// 16. Encoding the octetstring type
func EncOctetString(input, min, max int) (val []uint8, err error) {

	bitLen := bits.Len(uint(input))
	octLen := (bitLen-1)/8 + 1

	if octLen > max {
		err = fmt.Errorf("OctetString: input value overflow.")
		return
	}

	if octLen < min || min == max {
		octLen = min
	}

	val = make([]uint8, octLen, octLen+1)
	offset := 0
	if min != max {
		val[offset] = uint8(octLen - min)
		val = append(val, 0)
		octLen++
		offset++
	}
	for oct := octLen + offset - 1; oct > offset-1; oct-- {
		val[oct] = uint8(input & 0xff)
		input >>= 8
	}
	return
}
// EncSequence return Sequence Preamble but it just returns 0x00 for now.
// 18. Encoding the sequence type
func EncSequence(input []uint8, extmark bool, markisexist bool, optnum int) (input []uint8) {
	if extmark == true {
		inputt <<= 1
	}
	if markisexist == true {
		input |= 0x01
	}
	input <<= optnum
	return
}

// EncSequenceOf return Sequence-Of Preamble.
// 19. Encoding the sequence-of type
var EncSequenceOf = EncEnumerated

// EncChoice returns CHOICE preamble
// 22. Encoding the choice type
var EncChoice = EncEnumerated

*/

