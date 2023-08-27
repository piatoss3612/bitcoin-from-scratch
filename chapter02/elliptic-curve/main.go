package main

import (
	"fmt"
	"math"
)

// 타원곡선의 점을 나타내는 구조체
type Point struct {
	x, y, a, b float64
}

// 타원곡선의 점을 생성하는 함수
func New(x, y, a, b float64) (*Point, error) {
	// 무한원점인지 확인
	if isInfinity(x, y) {
		return &Point{x: x, y: y, a: a, b: b}, nil
	}

	// 주어진 점이 타원곡선 위에 있는지 확인
	if !isOnCurve(x, y, a, b) {
		return nil, fmt.Errorf("(%.2f, %.2f) is not on the curve", x, y)
	}

	return &Point{x: x, y: y, a: a, b: b}, nil
}

// 타원곡선의 점을 문자열로 표현하는 함수 (Stringer 인터페이스 구현)
func (p Point) String() string {
	// 무한원점인지 확인
	if isInfinity(p.x, p.y) {
		return "Point(infinity)"
	}
	return fmt.Sprintf("Point(%.2f, %.2f)_%.2f_%.2f", p.x, p.y, p.a, p.b)
}

// 두 타원곡선의 점이 같은지 확인하는 함수
func (p Point) Equals(other Point) bool {
	// 두 점의 좌표가 같고 같은 타원곡선 위에 있는지 확인
	return samePoint(p.x, p.y, other.x, other.y) &&
		sameCurve(p.a, p.b, other.a, other.b)
}

// 두 타원곡선의 점이 다른지 확인하는 함수
func (p Point) NotEquals(other Point) bool {
	return !(samePoint(p.x, p.y, other.x, other.y) &&
		sameCurve(p.a, p.b, other.a, other.b))
}

// 두 타원곡선의 점을 더하는 함수
func (p Point) Add(other Point) (*Point, error) {
	// 같은 타원곡선 위에 있는지 확인
	if !sameCurve(p.a, p.b, other.a, other.b) {
		return nil, fmt.Errorf("points %s and %s are not on the same curve", p, other)
	}

	/* case1: 두 점이 x축에 수직인 직선 위에 있는 경우 */

	// p가 무한원점인지 확인
	if isInfinity(p.x, p.y) {
		return &other, nil
	}

	// other가 무한원점인지 확인
	if isInfinity(other.x, other.y) {
		return &p, nil
	}

	// 한 점에 그의 역원을 더하는 경우, 무한원점을 반환
	if areInverse(p.x, other.x, p.y, other.y) {
		return New(math.MaxFloat64, math.MaxFloat64, p.a, p.b)
	}

	/* case2: 두 점이 같은 경우 */

	// p와 other가 같은 점인지 확인
	if samePoint(p.x, p.y, other.x, other.y) {
		// case 2-1 예외 처리: 접선이 x축에 수직인 경우, 무한원점을 반환
		if p.y == 0 {
			return New(math.MaxFloat64, math.MaxFloat64, p.a, p.b)
		}
		// 접선의 기울기 구하기
		s := (3*p.x*p.x + p.a) / (2 * p.y)

		// 접선과 타원곡선의 교점 q의 좌표 구하기
		nx := s*s - 2*p.x
		ny := s*(nx-p.x) + p.y

		// 교점 q의 역원 구하기 (y축 대칭)
		ny = -ny

		return New(nx, ny, p.a, p.b)
	}

	/* case3: 두 점이 서로 다른 경우 */

	// p와 other를 지나는 직선의 기울기 구하기
	s := (other.y - p.y) / (other.x - p.x)

	// p와 other를 지나는 직선이 타원곡선과 만나는 다른 한 점 q의 좌표 구하기
	nx := s*s - p.x - other.x
	ny := s*(nx-p.x) + p.y

	// p와 other의 점 덧셈의 결과인 q의 역원 구하기 (y축 대칭)
	ny = -ny

	return New(nx, ny, p.a, p.b)
}

// 무한원점인지 확인하는 함수
func isInfinity(x, y float64) bool {
	return x == math.MaxFloat64 && y == math.MaxFloat64
}

// 타원곡선 위에 있는지 확인하는 함수
func isOnCurve(x, y, a, b float64) bool {
	return y*y == x*x*x+a*x+b
}

// 두 점이 서로 역원인지 확인하는 함수
func areInverse(x1, x2, y1, y2 float64) bool {
	return x1 == x2 && y1 != y2
}

// 두 타원곡선이 같은지 확인하는 함수
func sameCurve(a1, b1, a2, b2 float64) bool {
	return a1 == a2 && b1 == b2
}

// 두 점이 같은지 확인하는 함수
func samePoint(x1, y1, x2, y2 float64) bool {
	return x1 == x2 && y1 == y2
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

	p3, err := New(1, 4, 8, 7)
	if err != nil {
		panic(err)
	}

	fmt.Println(p3)

	p4, err := New(-1, 1, 5, 7)
	if err != nil {
		panic(err)
	}

	fmt.Println(p4)

	inf, err := p1.Add(*p4)
	if err != nil {
		panic(err)
	}

	fmt.Println(inf)

	p5, err := p1.Add(*p2)
	if err != nil {
		panic(err)
	}

	fmt.Println(p5)

	p6, err := p1.Add(*p1)
	if err != nil {
		panic(err)
	}

	fmt.Println(p6)
}
