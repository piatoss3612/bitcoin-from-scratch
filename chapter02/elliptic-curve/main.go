package main

import "fmt"

// 타원곡선의 점을 나타내는 구조체
type Point struct {
	x, y, a, b int
}

// 타원곡선의 점을 생성하는 함수
func New(x, y, a, b int) (*Point, error) {
	if y*y != x*x*x+a*x+b {
		return nil, fmt.Errorf("(%d, %d) is not on the curve", x, y)
	}

	return &Point{x: x, y: y, a: a, b: b}, nil
}

// 타원곡선의 점을 문자열로 표현하는 함수 (Stringer 인터페이스 구현)
func (p Point) String() string {
	return fmt.Sprintf("Point(%d, %d)_%d_%d", p.x, p.y, p.a, p.b)
}

// 두 타원곡선의 점이 같은지 확인하는 함수
func (p Point) Equals(other Point) bool {
	return p.x == other.x && p.y == other.y &&
		p.a == other.a && p.b == other.b
}

// 두 타원곡선의 점이 다른지 확인하는 함수
func (p Point) NotEquals(other Point) bool {
	return !(p.x == other.x && p.y == other.y &&
		p.a == other.a && p.b == other.b)
}

func main() {
	p1, err := New(-1, -1, 5, 7)
	if err != nil {
		panic(err)
	}

	fmt.Println(p1)

	p2, err := New(2, 5, 5, 7)
	if err != nil {
		panic(err)
	}

	fmt.Println(p2)

	fmt.Println(p1.Equals(*p2))
	fmt.Println(p1.NotEquals(*p2))

	p3, err := New(-1, -2, 5, 7)
	if err != nil {
		panic(err)
	}

	fmt.Println(p3)
}
