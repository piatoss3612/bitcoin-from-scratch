package ecc

import (
	"fmt"
)

type Point struct {
	x, y, a, b *FieldElement
}

// 타원곡선의 점을 생성하는 함수
func New(x, y, a, b *FieldElement) (*Point, error) {
	// 무한원점인지 확인
	if isInfinity(x, y) {
		return &Point{x: x, y: y, a: a, b: b}, nil
	}

	// 주어진 점이 타원곡선 위에 있는지 확인
	if !isOnCurve(x, y, a, b) {
		return nil, fmt.Errorf("(%s, %s) is not on the curve", x, y)
	}

	return &Point{x: x, y: y, a: a, b: b}, nil
}

// 타원곡선의 점을 문자열로 표현하는 함수 (Stringer 인터페이스 구현)
func (p Point) String() string {
	// 무한원점인지 확인
	if isInfinity(p.x, p.y) {
		return "Point(infinity)"
	}
	return fmt.Sprintf("Point(%d, %d)_%d_%d FieldElement(%d)", p.x.num, p.y.num, p.a.num, p.b.num, p.x.prime)
}

// 두 타원곡선의 점이 같은지 확인하는 함수
func (p Point) Equal(other Point) bool {
	// 두 점의 좌표가 같고 같은 타원곡선 위에 있는지 확인
	return samePoint(p.x, p.y, other.x, other.y) &&
		sameCurve(p.a, p.b, other.a, other.b)
}

// 두 타원곡선의 점이 다른지 확인하는 함수
func (p Point) NotEqual(other Point) bool {
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
		return New(nil, nil, p.a, p.b)
	}

	/* case2: 두 점이 같은 경우 */

	// p와 other가 같은 점인지 확인
	if samePoint(p.x, p.y, other.x, other.y) {
		// case 2-1 예외 처리: 접선이 x축에 수직인 경우, 무한원점을 반환
		if p.y.num == 0 {
			return New(nil, nil, p.a, p.b)
		}
		// 접선의 기울기 구하기
		p1, err := NewFieldElement(3, p.x.prime)
		if err != nil {
			return nil, err
		}

		p2, err := p.x.Pow(2)
		if err != nil {
			return nil, err
		}

		p3, err := p1.Mul(*p2)
		if err != nil {
			return nil, err
		}

		p4, err := p3.Add(*p.a)
		if err != nil {
			return nil, err
		}

		c1, err := NewFieldElement(2, p.x.prime)
		if err != nil {
			return nil, err
		}

		c2, err := c1.Mul(*p.y)
		if err != nil {
			return nil, err
		}

		s, err := p4.Div(*c2)
		if err != nil {
			return nil, err
		}

		// 접선과 타원곡선의 교점 q의 좌표 구하기
		s2, err := p.x.Pow(2)
		if err != nil {
			return nil, err
		}

		x1, err := NewFieldElement(2, p.x.prime)
		if err != nil {
			return nil, err
		}

		x2, err := x1.Mul(*p.x)
		if err != nil {
			return nil, err
		}

		nx, err := s2.Sub(*x2)
		if err != nil {
			return nil, err
		}

		y1, err := p.x.Sub(*nx)
		if err != nil {
			return nil, err
		}

		y2, err := s.Mul(*y1)
		if err != nil {
			return nil, err
		}

		ny, err := y2.Sub(*p.y)
		if err != nil {
			return nil, err
		}

		return New(nx, ny, p.a, p.b)
	}

	/* case3: 두 점이 서로 다른 경우 */

	// p와 other를 지나는 직선의 기울기 구하기
	s1, err := other.y.Sub(*p.y)
	if err != nil {
		return nil, err
	}

	s2, err := other.x.Sub(*p.x)
	if err != nil {
		return nil, err
	}

	s, err := s1.Div(*s2)
	if err != nil {
		return nil, err
	}

	// p와 other를 지나는 직선이 타원곡선과 만나는 다른 한 점 q의 좌표 구하기
	x1, err := s.Pow(2)
	if err != nil {
		return nil, err
	}

	x2, err := x1.Sub(*p.x)
	if err != nil {
		return nil, err
	}

	nx, err := x2.Sub(*other.x)
	if err != nil {
		return nil, err
	}

	y1, err := p.x.Sub(*nx)
	if err != nil {
		return nil, err
	}

	y2, err := s.Mul(*y1)
	if err != nil {
		return nil, err
	}

	ny, err := y2.Sub(*p.y)
	if err != nil {
		return nil, err
	}

	return New(nx, ny, p.a, p.b)
}

// 무한원점인지 확인하는 함수
func isInfinity(x, y *FieldElement) bool {
	return x == nil && y == nil
}

// 타원곡선 위에 있는지 확인하는 함수
func isOnCurve(x, y, a, b *FieldElement) bool {
	left, err := y.Pow(2)
	if err != nil {
		return false
	}

	r1, err := x.Pow(3)
	if err != nil {
		return false
	}

	r2, err := a.Mul(*x)
	if err != nil {
		return false
	}

	r3, err := r1.Add(*r2)
	if err != nil {
		return false
	}

	right, err := r3.Add(*b)
	if err != nil {
		return false
	}

	return left.Equal(*right)
}

// 두 점이 서로 역원인지 확인하는 함수
func areInverse(x1, x2, y1, y2 *FieldElement) bool {
	return x1.Equal(*x2) && y1.NotEqual(*y2)
}

// 두 타원곡선이 같은지 확인하는 함수
func sameCurve(a1, b1, a2, b2 *FieldElement) bool {
	return a1.Equal(*a2) && b1.Equal(*b2)
}

// 두 점이 같은지 확인하는 함수
func samePoint(x1, y1, x2, y2 *FieldElement) bool {
	return x1.Equal(*x2) && y1.Equal(*y2)
}