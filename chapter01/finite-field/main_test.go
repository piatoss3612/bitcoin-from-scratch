package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_FieldElement_New(t *testing.T) {
	tcs := []struct {
		num      int
		prime    int
		hasError bool
		msg      string
	}{
		{num: 2, prime: 31, hasError: false, msg: ""},
		{num: -1, prime: 31, hasError: true, msg: "num -1 not in field range 0 to 30"},
		{num: 31, prime: 31, hasError: true, msg: "num 31 not in field range 0 to 30"},
	}

	for _, tc := range tcs {
		_, err := New(tc.num, tc.prime)
		if tc.hasError && err == nil {
			t.Error("expected error but got nil")
		}

		if !tc.hasError && err != nil {
			t.Error("expected nil but got error")
		}

		if tc.hasError && err != nil && err.Error() != tc.msg {
			t.Errorf("expected %s but got %s", tc.msg, err.Error())
		}
	}
}

func Test_FieldElement_String(t *testing.T) {
	tcs := []struct {
		num      int
		prime    int
		expected string
	}{
		{num: 2, prime: 31, expected: "FieldElement_31(2)"},
		{num: 17, prime: 31, expected: "FieldElement_31(17)"},
		{num: 29, prime: 31, expected: "FieldElement_31(29)"},
	}

	for _, tc := range tcs {
		f, _ := New(tc.num, tc.prime)
		if f.String() != tc.expected {
			t.Errorf("expected %s but got %s", tc.expected, f.String())
		}
	}
}

func Test_FieldElement_Equals(t *testing.T) {
	tcs := []struct {
		num      int
		prime    int
		other    FieldElement
		expected bool
	}{
		{num: 2, prime: 31, other: FieldElement{num: 2, prime: 31}, expected: true},
		{num: 2, prime: 31, other: FieldElement{num: 15, prime: 31}, expected: false},
		{num: 2, prime: 31, other: FieldElement{num: 2, prime: 17}, expected: false},
	}

	for _, tc := range tcs {
		f, _ := New(tc.num, tc.prime)
		if f.Equals(tc.other) != tc.expected {
			t.Errorf("expected %t but got %t", tc.expected, f.Equals(tc.other))
		}
	}
}

func Test_FieldElement_Add(t *testing.T) {
	tcs := []struct {
		num1     int
		prime1   int
		num2     int
		prime2   int
		expected int
		hasError bool
		msg      string
	}{
		{num1: 2, prime1: 31, num2: 15, prime2: 31, expected: 17, hasError: false, msg: ""},
		{num1: 23, prime1: 31, num2: 15, prime2: 31, expected: 7, hasError: false, msg: ""},
		{num1: 2, prime1: 31, num2: 15, prime2: 17, expected: 0, hasError: true, msg: "cannot add two numbers in different Fields 31, 17"},
	}

	for _, tc := range tcs {
		f1, _ := New(tc.num1, tc.prime1)
		f2, _ := New(tc.num2, tc.prime2)
		f3, err := f1.Add(*f2)
		if tc.hasError && err == nil {
			t.Error("expected error but got nil")
		}

		if !tc.hasError && err != nil {
			t.Error("expected nil but got error")
		}

		if tc.hasError && err != nil && err.Error() != tc.msg {
			t.Errorf("expected %s but got %s", tc.msg, err.Error())
		}

		if !tc.hasError && f3.num != tc.expected {
			t.Errorf("expected %d but got %d", tc.expected, f3.num)
		}
	}
}

func Test_FieldElement_Sub(t *testing.T) {
	tcs := []struct {
		num1     int
		prime1   int
		num2     int
		prime2   int
		expected int
		hasError bool
		msg      string
	}{
		{num1: 2, prime1: 31, num2: 15, prime2: 31, expected: 18, hasError: false, msg: ""},
		{num1: 2, prime1: 31, num2: 15, prime2: 17, expected: 0, hasError: true, msg: "cannot add two numbers in different Fields 31, 17"},
	}

	for _, tc := range tcs {
		f1, _ := New(tc.num1, tc.prime1)
		f2, _ := New(tc.num2, tc.prime2)
		f3, err := f1.Sub(*f2)

		if tc.hasError && err == nil {
			t.Error("expected error but got nil")
		}

		if !tc.hasError && err != nil {
			t.Error("expected nil but got error")
		}

		if tc.hasError && err != nil && err.Error() != tc.msg {
			t.Errorf("expected %s but got %s", tc.msg, err.Error())
		}

		if !tc.hasError && f3.num != tc.expected {
			t.Errorf("expected %d but got %d", tc.expected, f3.num)
		}
	}
}

func Test_FieldElement_Mul(t *testing.T) {
	tcs := []struct {
		num1     int
		prime1   int
		num2     int
		prime2   int
		expected int
		hasError bool
		msg      string
	}{
		{num1: 24, prime1: 31, num2: 19, prime2: 31, expected: 22, hasError: false, msg: ""},
		{num1: 24, prime1: 31, num2: 14, prime2: 17, expected: 0, hasError: true, msg: "cannot add two numbers in different Fields 31, 17"},
	}

	for _, tc := range tcs {
		f1, _ := New(tc.num1, tc.prime1)
		f2, _ := New(tc.num2, tc.prime2)
		f3, err := f1.Mul(*f2)

		if tc.hasError && err == nil {
			t.Error("expected error but got nil")
		}

		if !tc.hasError && err != nil {
			t.Error("expected nil but got error")
		}

		if tc.hasError && err != nil && err.Error() != tc.msg {
			t.Errorf("expected %s but got %s", tc.msg, err.Error())
		}

		if !tc.hasError && f3.num != tc.expected {
			t.Errorf("expected %d but got %d", tc.expected, f3.num)
		}
	}
}

func Test_FieldElement_Pow(t *testing.T) {
	tcs := []struct {
		num      int
		prime    int
		exp      int
		expected int
		hasError bool
		msg      string
	}{
		{num: 3, prime: 31, exp: 4, expected: 19, hasError: false, msg: ""},
		{num: 3, prime: 31, exp: -3, expected: 23, hasError: false, msg: ""},
		{num: 3, prime: 31, exp: -17, expected: 24, hasError: false, msg: ""},
		{num: 3, prime: 31, exp: 0, expected: 1, hasError: false, msg: ""},
		{num: 3, prime: 31, exp: 30, expected: 1, hasError: false, msg: ""},
		{num: 3, prime: 31, exp: 60, expected: 1, hasError: false, msg: ""},
	}

	for _, tc := range tcs {
		f, _ := New(tc.num, tc.prime)
		f2, err := f.Pow(tc.exp)

		if tc.hasError && err == nil {
			t.Error("expected error but got nil")
		}

		if !tc.hasError && err != nil {
			t.Error("expected nil but got error")
		}

		if tc.hasError && err != nil && err.Error() != tc.msg {
			t.Errorf("expected %s but got %s", tc.msg, err.Error())
		}

		if !tc.hasError && f2.num != tc.expected {
			t.Errorf("expected %d but got %d", tc.expected, f2.num)
		}
	}
}

func Test_FieldElement_Div(t *testing.T) {
	tcs := []struct {
		num1     int
		prime1   int
		num2     int
		prime2   int
		expected int
		hasError bool
		msg      string
	}{
		{num1: 3, prime1: 31, num2: 24, prime2: 31, expected: 4, hasError: false, msg: ""},
		{num1: 17, prime1: 31, num2: 3, prime2: 31, expected: 16, hasError: false, msg: ""},
		{num1: 4, prime1: 31, num2: 7, prime2: 29, expected: 0, hasError: true, msg: "cannot add two numbers in different Fields 31, 29"},
	}

	for _, tc := range tcs {
		f1, _ := New(tc.num1, tc.prime1)
		f2, _ := New(tc.num2, tc.prime2)
		f3, err := f1.Div(*f2)

		if tc.hasError && err == nil {
			t.Error("expected error but got nil")
		}

		if !tc.hasError && err != nil {
			t.Error("expected nil but got error")
		}

		if tc.hasError && err != nil && err.Error() != tc.msg {
			t.Errorf("expected %s but got %s", tc.msg, err.Error())
		}

		if !tc.hasError && f3.num != tc.expected {
			t.Errorf("expected %d but got %d", tc.expected, f3.num)
		}
	}
}

func Test_main(t *testing.T) {
	oldOut := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	main()

	_ = w.Close()

	os.Stdout = oldOut

	out, _ := io.ReadAll(r)

	outStr := string(out)

	expected := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s",
		"false",
		"true",
		"FieldElement_13(0)",
		"FieldElement_13(1)",
		"FieldElement_13(3)",
		"FieldElement_13(8)",
		"FieldElement_13(12)",
	)

	if strings.EqualFold(outStr, expected) {
		t.Errorf("expected %s but got %s", expected, outStr)
	}
}
