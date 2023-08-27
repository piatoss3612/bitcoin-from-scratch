package main

import (
	"fmt"
	"math"
)

// 타원곡선의 점을 나타내는 구조체
type Point struct {
	x, y, a, b int
}

// 타원곡선의 점을 생성하는 함수
func New(x, y, a, b int) (*Point, error) {
	// 무한원점인지 확인
	if x == math.MaxInt64 && y == math.MaxInt64 {
		return &Point{x: x, y: y, a: a, b: b}, nil
	}

	// 주어진 점이 타원곡선 위에 있는지 확인
	if y*y != x*x*x+a*x+b {
		return nil, fmt.Errorf("(%d, %d) is not on the curve", x, y)
	}

	return &Point{x: x, y: y, a: a, b: b}, nil
}

// 타원곡선의 점을 문자열로 표현하는 함수 (Stringer 인터페이스 구현)
func (p Point) String() string {
	// 무한원점인지 확인
	if p.x == math.MaxInt64 && p.y == math.MaxInt64 {
		return "Point(infinity)"
	}
	return fmt.Sprintf("Point(%d, %d)_%d_%d", p.x, p.y, p.a, p.b)
}

// 두 타원곡선의 점이 같은지 확인하는 함수
func (p Point) Equals(other Point) bool {
	// 두 점의 좌표가 같고 같은 타원곡선 위에 있는지 확인
	return p.x == other.x && p.y == other.y &&
		p.a == other.a && p.b == other.b
}

// 두 타원곡선의 점이 다른지 확인하는 함수
func (p Point) NotEquals(other Point) bool {
	return !(p.x == other.x && p.y == other.y &&
		p.a == other.a && p.b == other.b)
}

// 두 타원곡선의 점을 더하는 함수
func (p Point) Add(other Point) (*Point, error) {
	// 같은 타원곡선 위에 있는지 확인
	if p.a != other.a || p.b != other.b {
		return nil, fmt.Errorf("points %s and %s are not on the same curve", p, other)
	}

	/* case1: 두 점이 x축에 수직인 직선 위에 있는 경우 */

	// p가 무한원점인지 확인
	if p.x == math.MaxInt64 && p.y == math.MaxInt64 {
		return &other, nil
	}

	// other가 무한원점인지 확인
	if other.x == math.MaxInt64 && other.y == math.MaxInt64 {
		return &p, nil
	}

	// 한 점에 그의 역원을 더하는 경우, 무한원점을 반환
	if p.x == other.x && p.y != other.y {
		return New(math.MaxInt64, math.MaxInt64, p.a, p.b)
	}

	// TODO: 이후에 구현
	return nil, nil
}

func main() {
	p1, err := New(-1, -1, 5, 7)
	if err != nil {
		panic(err)
	}

	fmt.Println(p1)

	p2, err := New(18, 77, 5, 7)
	if err != nil {
		panic(err)
	}

	fmt.Println(p2)

	fmt.Println(p1.Equals(*p2))
	fmt.Println(p1.NotEquals(*p2))

	p3, err := New(1, 4, 8, 7)
	if err != nil {
		panic(err)
	}

	fmt.Println(p3)

	p4, err := p1.Add(*p3)
	if err != nil {
		panic(err)
	}

	fmt.Println(p4)
}
