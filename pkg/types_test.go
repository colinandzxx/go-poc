package pkg

import (
	"bytes"
	"fmt"
	"testing"
)

func TestByte32_String(t *testing.T) {
	expected := "0x0102030405060708090a0b0c0d0e0f1015161718191a1b1c1dd2d3d4d5d6d7d8"
	input := Byte32{
		1,  2,  3,  4,  5,  6,  7,  8,  9,  10,  11,  12,  13,  14,  15,  16,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 210, 211, 212, 213, 214, 215, 216,
	}
	output := input.String()
	if output != expected {
		t.Errorf("Expected %s got %s", expected, output)
	}

	//t.Logf("String(): %s", b32.String())
	//t.Logf("fmt ori: %v", b32)
	//t.Logf("fmt hex: %x", b32)
}

func TestByte32_Format(t *testing.T) {
	{
		expected := "[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,21,22,23,24,25,26,27,28,29,210,211,212,213,214,215,216]"
		input := Byte32{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 210, 211, 212, 213, 214, 215, 216,
		}
		output := bytes.Buffer{}
		fmt.Fprintf(&output, "%v", input)
		if output.String() != expected {
			t.Errorf("Expected %s got %s", expected, output.String())
		}
	}

	{
		expected := "0x0102030405060708090a0b0c0d0e0f1015161718191a1b1c1dd2d3d4d5d6d7d8"
		input := Byte32{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 210, 211, 212, 213, 214, 215, 216,
		}
		output := bytes.Buffer{}
		fmt.Fprintf(&output, "%x", input)
		if output.String() != expected {
			t.Errorf("Expected %s got %s", expected, output.String())
		}
	}
}

func TestByte64_String(t *testing.T) {
	expected := "0x0102030405060708090a0b0c0d0e0f1015161718191a1b1c1dd2d3d4d5d6d7d8" +
		"0102030405060708090a0b0c0d0e0f1015161718191a1b1c1dd2d3d4d5d6d7d8"
	input := Byte64{
		1,  2,  3,  4,  5,  6,  7,  8,  9,  10,  11,  12,  13,  14,  15,  16,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 210, 211, 212, 213, 214, 215, 216,
		1,  2,  3,  4,  5,  6,  7,  8,  9,  10,  11,  12,  13,  14,  15,  16,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 210, 211, 212, 213, 214, 215, 216,
	}
	output := input.String()
	if output != expected {
		t.Errorf("Expected %s got %s", expected, output)
	}

	//t.Logf("String(): %s", b32.String())
	//t.Logf("fmt ori: %v", b32)
	//t.Logf("fmt hex: %x", b32)
}

func TestByte64_Format(t *testing.T) {
	{
		expected := "[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,21,22,23,24,25,26,27,28,29,210,211,212,213,214,215,216," +
			"1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,21,22,23,24,25,26,27,28,29,210,211,212,213,214,215,216]"
		input := Byte64{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 210, 211, 212, 213, 214, 215, 216,
			1,  2,  3,  4,  5,  6,  7,  8,  9,  10,  11,  12,  13,  14,  15,  16,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 210, 211, 212, 213, 214, 215, 216,
		}
		output := bytes.Buffer{}
		fmt.Fprintf(&output, "%v", input)
		if output.String() != expected {
			t.Errorf("Expected %s got %s", expected, output.String())
		}
	}

	{
		expected := "0x0102030405060708090a0b0c0d0e0f1015161718191a1b1c1dd2d3d4d5d6d7d8" +
			"0102030405060708090a0b0c0d0e0f1015161718191a1b1c1dd2d3d4d5d6d7d8"
		input := Byte64{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 210, 211, 212, 213, 214, 215, 216,
			1,  2,  3,  4,  5,  6,  7,  8,  9,  10,  11,  12,  13,  14,  15,  16,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 210, 211, 212, 213, 214, 215, 216,
		}
		output := bytes.Buffer{}
		fmt.Fprintf(&output, "%x", input)
		if output.String() != expected {
			t.Errorf("Expected %s got %s", expected, output.String())
		}
	}
}
