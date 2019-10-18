// Copyright 2019 hhorai. All rights reserved.
// Use of this source code is governed by a MIT license that can be found
// in the LICENSE file.

// Package per is implementation for Basic Pckage Encoding Rule (PER) in
// ALIGNED variant.
package per

import (
	"fmt"
	"math/bits"
)

// MergeBitField is utility function for merging bit-field.
// e.g. preamble or short integer value is not octet alined value, so
// those fields need to be packed in same octets.
func MergeBitField(in1 []uint8, inlen1 int, in2 []uint8, inlen2 int) (
	out []uint8, outlen int) {
	/*
	   ex1.
	   in1(len=4)  nil
	   in2(len=14) bxx11 1010 1111 0000
	   out(len=18) b1110 1011 1100 00xx

	   ex2.
	   in1(len=4)  b1010 xxxx
	   in2(len=14) b1110 1011 1100 00xx
	   out(len=18) b1010 1110 1011 1100 00xx
	*/

	if in1 == nil {
		out, outlen = ShiftLeftMost(in2, inlen2)
		return
	}

	out = make([]uint8, len(in1), len(in1))
	out = append(out, in2...)
	out = ShiftLeft(out, len(in1)*8-inlen1)
	for n := 0; n < len(in1); n++ {
		out[n] |= in1[n]
	}
	outlen = inlen1 + inlen2

	octetlen := (outlen-1)/8 + 1
	out = out[:octetlen]
	return
}

// ShiftLeft is utility function to left shift the octet values.
func ShiftLeft(in []uint8, shiftlen int) (out []uint8) {
	out = in
	for n := 0; n < shiftlen; n++ {
		overflow := false
		for m := len(out) - 1; m >= 0; m-- {
			do := false
			if overflow == true {
				do = true
				overflow = false
			}
			if out[m]&0x80 == 0x80 {
				overflow = true
			}
			out[m] <<= 1
			if do == true {
				out[m] |= 0x01
			}
		}
	}
	return
}

// ShiftRight is utility function to right shift the octet values.
func ShiftRight(in []uint8, shiftlen int) (out []uint8) {
	out = in
	for n := 0; n < shiftlen; n++ {
		underflow := false
		for m := 0; m < len(out); m++ {
			do := false
			if underflow == true {
				do = true
				underflow = false
			}
			if out[m]&0x1 == 0x1 {
				underflow = true
			}
			out[m] >>= 1
			if do == true {
				out[m] |= 0x80
			}
		}
	}
	return
}

// ShiftLeftMost is utility function to shift the octet values to the leftmost.
func ShiftLeftMost(in []uint8, inlen int) (out []uint8, outlen int) {
	out = in
	outlen = inlen
	out = ShiftLeft(out, len(in)*8-inlen)
	return
}

// EncConstrainedWholeNumber is the implementation for
// 10.5 Encoding of constrained whole number.
func EncConstrainedWholeNumber(input, min, max int) (
	v []uint8, bitlen int, err error) {

	if input < min || input > max {
		err = fmt.Errorf("EncConstrainedWholeNumber: "+
			"input value=%d is out of range. "+
			"(should be %d <= %d)", input, min, max)
		return
	}

	inputRange := max - min + 1
	inputEnc := input - min

	switch {
	case inputRange == 1: // empty bit-field
		return
	case inputRange < 256: // the bit-field case
		bitlen = bits.Len(uint(inputRange))
		//v = append(v, uint8(inputEnc << uint((8 - bitlen))))
		v = append(v, uint8(inputEnc))
		return
	case inputRange == 256: // the one-octet case
		bitlen = 8
		v = append(v, uint8(inputEnc))
		return
	case inputRange <= 65536: // the two-octet case
		bitlen = 16
		v = append(v, uint8((inputEnc>>8)&0xff))
		v = append(v, uint8(inputEnc&0xff))
		return
	case inputRange > 65537: // the indefinite length case
		// not implemented yet
		err = fmt.Errorf("EncConstrainedWholeNumber: "+
			"not implemented yet for %d", input)
		return
	}
	err = fmt.Errorf("EncConstrainedWholeNumber: "+
		"invalid range min=%d, max=%d", min, max)
	return
}

func encConstrainedWholeNumberWithExtmark(input, min, max int, extmark bool) (
	v []uint8, bitlen int, err error) {
	v, bitlen, err = EncConstrainedWholeNumber(input, min, max)
	if err != nil {
		return
	}
	if extmark == true {
		switch {
		case bitlen%8 == 0:
			bitlen += 8
			v = append([]uint8{0x00}, v...)
		case bitlen < 8:
			bitlen++
		}
	}
	ShiftLeftMost(v, bitlen)
	return
}

// EncInteger is the implementation for
// 12. Encoding the integer type
// but it is only for the case of single value and constrained whole nuber.
func EncInteger(input, min, max int, extmark bool) (
	v []uint8, bitlen int, err error) {

	if min == max { // 12.2.1 single value
		if extmark == true {
			bitlen = 1
			v = make([]uint8, 1, 1)
		}
		return
	}

	// 12.2.2 constrained whole number
	v, bitlen, err = encConstrainedWholeNumberWithExtmark(input,
		min, max, extmark)
	return
}

// EncEnumerated is the implementation for
// 13. Encoding the enumerated type
func EncEnumerated(input, min, max int, extmark bool) (
	v []uint8, bitlen int, err error) {
	v, bitlen, err =
		encConstrainedWholeNumberWithExtmark(input, min, max, extmark)
	return
}

// EncBitString returns multi-byte BIT STRING
// 15. Encoding the bitstering type
func EncBitString(input []uint8, inputlen, min, max int, extmark bool) (
	v []uint8, bitlen int, err error) {

	//bitLen := bits.Len(uint(input))
	if inputlen < min || inputlen > max {
		err = fmt.Errorf("EncBitString: "+
			"input len(value)=%d is out of range. "+
			"(should be %d <= %d)", inputlen, min, max)
		return
	}

	if min == max && min != inputlen {
		err = fmt.Errorf("EncBitString: "+
			"input len(value)=%d must be %d", inputlen, min)
		return
	}

	if len(input)*8 < inputlen {
		err = fmt.Errorf("EncBitString: "+
			"input len(value)=%d is too short.", len(input))
		return
	}

	v, bitlen = ShiftLeftMost(input, inputlen)

	if min == max {
		// fixed length case. not implemented yet.
		switch {
		case min < 17:
		case min > 16 && min < 65537:
		}
		return
	}

	// range is constrained whole number.
	pv, plen, _ := encConstrainedWholeNumberWithExtmark(inputlen,
		min, max, extmark)

	v, bitlen = MergeBitField(pv, plen, v, bitlen)
	return
}

// EncOctetString returns multi-byte OCTET STRING
// 16. Encoding the octetstring type
//
// - the length of returned value can be calculated by len().
// - returned value can be len(value) == 0 if the specified octet string has
//   fixed length and the lenght is less than 3. And then the octet string is
//   encoded as bit field.
func EncOctetString(input []uint8, min, max int, extmark bool) (
	pv []uint8, plen int, v []uint8, err error) {

	inputlen := len(input)
	if inputlen < min || inputlen > max {
		err = fmt.Errorf("EncOctetString: "+
			"input len(value)=%d is out of range. "+
			"(should be %d <= %d)", inputlen, min, max)
		return
	}

	v = input
	plen = 0

	if min == max {
		if extmark == false {
			return
		}

		pv = []uint8{0x00}

		switch {
		case min < 3:
			pv = append(pv, v...)
			plen = inputlen*8 + 1
			v = []uint8{}
		case min < 65537:
			plen = 1
		}
		pv, plen = ShiftLeftMost(pv, plen)
		return
	}

	// range is constrained whole number.
	pv, plen, perr :=
		encConstrainedWholeNumberWithExtmark(inputlen,
			min, max, extmark)

	if perr != nil {
		err = fmt.Errorf("EncOctetString: unexpected error.")
		return
	}
	return
}

// EncSequence return Sequence Preamble but it just returns 0x00 for now.
// 18. Encoding the sequence type
func EncSequence(extmark bool, optnum int, optflag uint) (
	pv []uint8, plen int, err error) {
	if optnum > 7 {
		err = fmt.Errorf("EncSequence: "+
			"optnum=%d is not implemented yet. (should be < 8)",
			optnum)
		return
	}
	if extmark == true {
		plen++
	}
	plen += optnum
	pv = make([]uint8, 1, 1)
	pv[0] |= uint8(optflag)
	pv, plen = ShiftLeftMost(pv, plen)
	return
}

// EncSequenceOf return Sequence-Of Preamble.
// 19. Encoding the sequence-of type
var EncSequenceOf = EncEnumerated

// EncChoice is the implementation for
// 22. Encoding the choice type
func EncChoice(input, min, max int, extmark bool) (
	pv []uint8, plen int, err error) {
	pv, plen, err = EncInteger(input, min, max, extmark)
	pv, plen = ShiftLeftMost(pv, plen)
	return
}
