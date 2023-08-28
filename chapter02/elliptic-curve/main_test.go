package main

import (
	"io"
	"math"
	"os"
	"strings"
	"testing"
)

func Test_New(t *testing.T) {
	tcs := []struct {
		x        float64
		y        float64
		a        float64
		b        float64
		hasError bool
		msg      string
	}{
		{0, 1, 5, 7, true, "(0.00, 1.00) is not on the curve"},
		{-1, -1, 5, 7, false, ""},
		{math.MaxFloat64, math.MaxFloat64, 5, 7, false, ""},
	}

	for _, tc := range tcs {
		_, err := New(tc.x, tc.y, tc.a, tc.b)

		if tc.hasError && err == nil {
			t.Errorf("Expected error but got nil")
		}

		if !tc.hasError && err != nil {
			t.Errorf("Expected no error but got %v", err)
		}

		if tc.hasError && err != nil && err.Error() != tc.msg {
			t.Errorf("Expected error message %s but got %s", tc.msg, err.Error())
		}
	}
}

func Test_String(t *testing.T) {
	tcs := []struct {
		x        float64
		y        float64
		a        float64
		b        float64
		expected string
	}{
		{-1, -1, 5, 7, "Point(-1.00, -1.00)_5.00_7.00"},
		{math.MaxFloat64, math.MaxFloat64, 5, 7, "Point(infinity)"},
	}

	for _, tc := range tcs {
		p, _ := New(tc.x, tc.y, tc.a, tc.b)
		actual := p.String()

		if actual != tc.expected {
			t.Errorf("Expected %s but got %s", tc.expected, actual)
		}
	}
}

func Test_Equals(t *testing.T) {
	tcs := []struct {
		x1       float64
		y1       float64
		a1       float64
		b1       float64
		x2       float64
		y2       float64
		a2       float64
		b2       float64
		expected bool
	}{
		{-1, -1, 5, 7, -1, -1, 5, 7, true},
		{-1, -1, 5, 7, -1, 1, 5, 7, false},
		{-1, -1, 5, 7, math.MaxFloat64, math.MaxFloat64, 5, 7, false},
		{math.MaxFloat64, math.MaxFloat64, 5, 7, math.MaxFloat64, math.MaxFloat64, 5, 7, true},
	}

	for _, tc := range tcs {
		p1, _ := New(tc.x1, tc.y1, tc.a1, tc.b1)
		p2, _ := New(tc.x2, tc.y2, tc.a2, tc.b2)
		actual := p1.Equals(*p2)

		if actual != tc.expected {
			t.Errorf("Expected %v but got %v", tc.expected, actual)
		}
	}
}

func Test_NotEquals(t *testing.T) {
	tcs := []struct {
		x1       float64
		y1       float64
		a1       float64
		b1       float64
		x2       float64
		y2       float64
		a2       float64
		b2       float64
		expected bool
	}{
		{-1, -1, 5, 7, -1, -1, 5, 7, false},
		{-1, -1, 5, 7, -1, 1, 5, 7, true},
		{-1, -1, 5, 7, math.MaxFloat64, math.MaxFloat64, 5, 7, true},
		{math.MaxFloat64, math.MaxFloat64, 5, 7, math.MaxFloat64, math.MaxFloat64, 5, 7, false},
	}

	for _, tc := range tcs {
		p1, _ := New(tc.x1, tc.y1, tc.a1, tc.b1)
		p2, _ := New(tc.x2, tc.y2, tc.a2, tc.b2)
		actual := p1.NotEquals(*p2)

		if actual != tc.expected {
			t.Errorf("Expected %v but got %v", tc.expected, actual)
		}
	}
}

func Test_Add(t *testing.T) {
	tcs := []struct {
		x1       float64
		y1       float64
		a1       float64
		b1       float64
		x2       float64
		y2       float64
		a2       float64
		b2       float64
		expected *Point
		hasError bool
		msg      string
	}{
		{-1, -1, 5, 7, 1, 4, 8, 7, nil, true, "points Point(-1.00, -1.00)_5.00_7.00 and Point(1.00, 4.00)_8.00_7.00 are not on the same curve"},
		{math.MaxFloat64, math.MaxFloat64, 8, 7, 1, 4, 8, 7, &Point{1, 4, 8, 7}, false, ""},
		{-1, -1, 5, 7, math.MaxFloat64, math.MaxFloat64, 5, 7, &Point{-1, -1, 5, 7}, false, ""},
		{-1, -1, 5, 7, -1, 1, 5, 7, &Point{math.MaxFloat64, math.MaxFloat64, 5, 7}, false, ""},
		{-1, 0, 2, 3, -1, 0, 2, 3, &Point{math.MaxFloat64, math.MaxFloat64, 2, 3}, false, ""},
		{-1, -1, 5, 7, -1, -1, 5, 7, &Point{18, 77, 5, 7}, false, ""},
		{-1, -1, 5, 7, 2, 5, 5, 7, &Point{3, -7, 5, 7}, false, ""},
	}

	for _, tc := range tcs {
		p1, _ := New(tc.x1, tc.y1, tc.a1, tc.b1)
		p2, _ := New(tc.x2, tc.y2, tc.a2, tc.b2)
		actual, err := p1.Add(*p2)

		if tc.hasError && err == nil {
			t.Errorf("Expected error but got nil")
		}

		if !tc.hasError && err != nil {
			t.Errorf("Expected no error but got %v", err)
		}

		if tc.hasError && err != nil && err.Error() != tc.msg {
			t.Errorf("Expected error message %s but got %s", tc.msg, err.Error())
		}

		if tc.expected != nil && !actual.Equals(*tc.expected) {
			t.Errorf("Expected %v but got %v", tc.expected, actual)
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

	expected := `
	Point(-1.00, -1.00)_5.00_7.00
	Point(2.00, 5.00)_5.00_7.00
	false
	true
	Point(1.00, 4.00)_8.00_7.00
	Point(-1.00, 1.00)_5.00_7.00
	Point(infinity)
	Point(3.00, -7.00)_5.00_7.00
	Point(18.00, 77.00)_5.00_7.00
	`

	if strings.EqualFold(outStr, expected) {
		t.Errorf("Expected %s but got %s", expected, outStr)
	}
}
