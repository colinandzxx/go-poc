package pkg

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var (
	byte32T = reflect.TypeOf(Byte32{})
	byte64T = reflect.TypeOf(Byte64{})
)

type Byte32 [32]byte

func (self Byte32) String() string {
	return HexEncode(self[:])
}

func (self Byte32) Format(s fmt.State, c rune) {
	switch c {
	case 'x', 'X':
		fmt.Fprintf(s, "%s", HexEncode(self[:]))
	default:
		buf, err := json.Marshal(self)
		if err != nil {
			panic(err)
		}
		s.Write(buf)
	}
}

// UnmarshalText parses a Byte32 in hex syntax.
func (self Byte32) UnmarshalText(input []byte) error {
	return unmarshalFixedText("Byte32", input, self[:])
}

// UnmarshalJSON parses a Byte32 in hex syntax.
func (self Byte32) UnmarshalJSON(input []byte) error {
	return unmarshalFixedJSON(byte32T, input, self[:])
}

// MarshalText implements encoding.TextMarshaler
func (self Byte32) MarshalText() ([]byte, error) {
	return marshalText(self[:])
}

type Byte64 [64]byte

func (self Byte64) String() string {
	return HexEncode(self[:])
}

func (self Byte64) Format(s fmt.State, c rune) {
	switch c {
	case 'x', 'X':
		fmt.Fprintf(s, "%s", HexEncode(self[:]))
	default:
		buf, err := json.Marshal(self)
		if err != nil {
			panic(err)
		}
		s.Write(buf)
	}
}

// UnmarshalText parses a Byte64 in hex syntax.
func (self Byte64) UnmarshalText(input []byte) error {
	return unmarshalFixedText("Byte64", input, self[:])
}

// UnmarshalJSON parses a Byte64 in hex syntax.
func (self Byte64) UnmarshalJSON(input []byte) error {
	return unmarshalFixedJSON(byte64T, input, self[:])
}

// MarshalText implements encoding.TextMarshaler
func (self Byte64) MarshalText() ([]byte, error) {
	return marshalText(self[:])
}