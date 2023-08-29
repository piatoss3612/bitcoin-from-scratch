package ecc

import (
	"fmt"
	"math/big"
)

var (
	A = 0
	B = 7
	N = "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"
	G Point
)

func init() {
	bigGx, _ := new(big.Int).SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	bigGy, _ := new(big.Int).SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)

	gx, err := NewS256Field(bigGx)
	if err != nil {
		panic(err)
	}

	gy, err := NewS256Field(bigGy)
	if err != nil {
		panic(err)
	}

	g, err := NewS256Point(gx, gy)
	if err != nil {
		panic(err)
	}

	G = g
}

type Point interface {
	fmt.Stringer
	X() FieldElement
	Y() FieldElement
	A() FieldElement
	B() FieldElement
	Equal(other Point) bool
	NotEqual(other Point) bool
	Add(other Point) (Point, error)
	Mul(coefficient *big.Int) (Point, error)
	Verify(z FieldElement, sig Signature) (bool, error)
}

type point struct {
	x, y, a, b FieldElement
}

// 타원곡선의 점을 생성하는 함수
func NewPoint(x, y, a, b FieldElement) (Point, error) {
	// 무한원점인지 확인
	if isInfinity(x, y) {
		return &point{x: x, y: y, a: a, b: b}, nil
	}

	// 주어진 점이 타원곡선 위에 있는지 확인
	if !isOnCurve(x, y, a, b) {
		return nil, fmt.Errorf("(%s, %s) is not on the curve", x, y)
	}

	return &point{x: x, y: y, a: a, b: b}, nil
}

func (p point) X() FieldElement {
	return p.x
}

func (p point) Y() FieldElement {
	return p.y
}

func (p point) A() FieldElement {
	return p.a
}

func (p point) B() FieldElement {
	return p.b
}

// 타원곡선의 점을 문자열로 표현하는 함수 (Stringer 인터페이스 구현)
func (p point) String() string {
	// 무한원점인지 확인
	if isInfinity(p.x, p.y) {
		return "Point(infinity)"
	}
	return fmt.Sprintf("Point(%s, %s)_%s_%s FieldElement(%s)", p.x.Num(), p.y.Num(), p.a.Num(), p.b.Num(), p.x.Prime())
}

// 두 타원곡선의 점이 같은지 확인하는 함수
func (p point) Equal(other Point) bool {
	// 두 점의 좌표가 같고 같은 타원곡선 위에 있는지 확인
	return samePoint(p.x, p.y, other.X(), other.Y()) &&
		sameCurve(p.a, p.b, other.A(), other.B())
}

// 두 타원곡선의 점이 다른지 확인하는 함수
func (p point) NotEqual(other Point) bool {
	return !(samePoint(p.x, p.y, other.X(), other.Y()) &&
		sameCurve(p.a, p.b, other.A(), other.B()))
}

// 두 타원곡선의 점을 더하는 함수
func (p point) Add(other Point) (Point, error) {

	// 같은 타원곡선 위에 있는지 확인
	if !sameCurve(p.a, p.b, other.A(), other.B()) {
		return nil, fmt.Errorf("points %s and %s are not on the same curve", p, other)
	}

	/* case1: 두 점이 x축에 수직인 직선 위에 있는 경우 */

	// p가 무한원점인지 확인
	if isInfinity(p.x, p.y) {
		return other, nil
	}

	// other가 무한원점인지 확인
	if isInfinity(other.X(), other.Y()) {
		return &p, nil
	}

	// 한 점에 그의 역원을 더하는 경우, 무한원점을 반환
	if areInverse(p.x, other.X(), p.y, other.Y()) {
		return NewPoint(nil, nil, p.a, p.b)
	}

	/* case2: 두 점이 서로 다른 경우 */

	if p.x.NotEqual(other.X()) {
		// p와 other를 지나는 직선의 기울기 구하기
		s1, err := other.Y().Sub(p.y)
		if err != nil {
			return nil, err
		}

		s2, err := other.X().Sub(p.x)
		if err != nil {
			return nil, err
		}

		s, err := s1.Div(s2)
		if err != nil {
			return nil, err
		}

		// p와 other를 지나는 직선이 타원곡선과 만나는 다른 한 점 q의 좌표 구하기
		x1, err := s.Pow(big.NewInt(2))
		if err != nil {
			return nil, err
		}

		x2, err := x1.Sub(p.x)
		if err != nil {
			return nil, err
		}

		nx, err := x2.Sub(other.X())
		if err != nil {
			return nil, err
		}

		y1, err := p.x.Sub(nx)
		if err != nil {
			return nil, err
		}

		y2, err := s.Mul(y1)
		if err != nil {
			return nil, err
		}

		ny, err := y2.Sub(p.y)
		if err != nil {
			return nil, err
		}

		return NewPoint(nx, ny, p.a, p.b)
	}

	/* case3: 두 점이 같은 경우 */

	// p와 other가 같은 점인지 확인
	if samePoint(p.x, p.y, other.X(), other.Y()) {
		// case 2-1 예외 처리: 접선이 x축에 수직인 경우, 무한원점을 반환
		if p.y.Num().Cmp(big.NewInt(0)) == 0 {
			return NewPoint(nil, nil, p.a, p.b)
		}
		// 접선의 기울기 구하기
		p1, err := NewFieldElement(big.NewInt(3), p.x.Prime())
		if err != nil {
			return nil, err
		}

		p2, err := p.x.Pow(big.NewInt(2))
		if err != nil {
			return nil, err
		}

		p3, err := p1.Mul(p2)
		if err != nil {
			return nil, err
		}

		p4, err := p3.Add(p.a)
		if err != nil {
			return nil, err
		}

		c1, err := NewFieldElement(big.NewInt(2), p.x.Prime())
		if err != nil {
			return nil, err
		}

		c2, err := c1.Mul(p.y)
		if err != nil {
			return nil, err
		}

		s, err := p4.Div(c2)
		if err != nil {
			return nil, err
		}

		// 접선과 타원곡선의 교점 q의 좌표 구하기
		s2, err := s.Pow(big.NewInt(2))
		if err != nil {
			return nil, err
		}

		x1, err := NewFieldElement(big.NewInt(2), p.x.Prime())
		if err != nil {
			return nil, err
		}

		x2, err := x1.Mul(p.x)
		if err != nil {
			return nil, err
		}

		nx, err := s2.Sub(x2)
		if err != nil {
			return nil, err
		}

		y1, err := p.x.Sub(nx)
		if err != nil {
			return nil, err
		}

		y2, err := s.Mul(y1)
		if err != nil {
			return nil, err
		}

		ny, err := y2.Sub(p.y)
		if err != nil {
			return nil, err
		}

		return NewPoint(nx, ny, p.a, p.b)
	}

	return nil, fmt.Errorf("unhandled case, (%s, %s) + (%s, %s)", p.x, p.y, other.X(), other.Y())
}

func (p point) Mul(coefficient *big.Int) (Point, error) {
	coef := coefficient  // 계수
	current := Point(&p) // 시작점으로 초기화

	result, err := NewPoint(nil, nil, p.a, p.b)
	if err != nil {
		return nil, err
	}

	// 이진수 전개법을 이용하여 타원곡선의 점 곱셈
	for coef.Cmp(big.NewInt(0)) == 1 {
		// 가장 오른쪽 비트가 1인지 확인
		if coef.Bit(0) == 1 {
			result, err = result.Add(current) // 현재 점을 결과에 더하기
			if err != nil {
				return nil, err
			}
		}
		current, err = current.Add(current) // 현재 점을 두 배로 만들기
		if err != nil {
			return nil, err
		}

		coef.Rsh(coef, 1) // 비트를 오른쪽으로 한 칸씩 이동
	}

	return result, nil
}

func (p point) Verify(z FieldElement, sig Signature) (bool, error) {
	// TODO: implement verify
	return false, nil
}

type s256Point struct {
	point
}

func NewS256Point(x, y FieldElement) (Point, error) {
	a, err := NewS256Field(big.NewInt(int64(A)))
	if err != nil {
		return nil, err
	}

	b, err := NewS256Field(big.NewInt(int64(B)))
	if err != nil {
		return nil, err
	}

	if isInfinity(x, y) {
		return &s256Point{point{x: x, y: y, a: a, b: b}}, nil
	}

	if !isOnCurve(x, y, a, b) {
		return nil, fmt.Errorf("(%s, %s) is not on the curve", x, y)
	}

	return &s256Point{point{x: x, y: y, a: a, b: b}}, nil
}

func (p s256Point) Mul(coefficient *big.Int) (Point, error) {
	n, ok := new(big.Int).SetString(N, 16)
	if !ok {
		return nil, fmt.Errorf("failed to convert N to big.Int")
	}

	coef := new(big.Int).Mod(coefficient, n)

	return p.point.Mul(coef)
}

func (p s256Point) Verify(z FieldElement, sig Signature) (bool, error) {
	n, ok := new(big.Int).SetString(N, 16)
	if !ok {
		return false, fmt.Errorf("failed to convert N to big.Int")
	}

	sInv, err := sig.S().Pow(big.NewInt(0).Sub(n, big.NewInt(2)))
	if err != nil {
		return false, err
	}

	u, err := z.Mul(sInv)
	if err != nil {
		return false, err
	}

	v, err := sig.R().Mul(sInv)
	if err != nil {
		return false, err
	}

	uG, err := G.Mul(u.Num())
	if err != nil {
		return false, err
	}

	vP, err := p.Mul(v.Num())
	if err != nil {
		return false, err
	}

	x, err := uG.Add(vP)
	if err != nil {
		return false, err
	}

	return x.X().Num().Cmp(sig.R().Num()) == 0, nil
}

// 무한원점인지 확인하는 함수
func isInfinity(x, y FieldElement) bool {
	return x == nil && y == nil
}

// 타원곡선 위에 있는지 확인하는 함수
func isOnCurve(x, y, a, b FieldElement) bool {
	left, err := y.Pow(big.NewInt(2))
	if err != nil {
		return false
	}

	r1, err := x.Pow(big.NewInt(3))
	if err != nil {
		return false
	}

	r2, err := a.Mul(x)
	if err != nil {
		return false
	}

	r3, err := r1.Add(r2)
	if err != nil {
		return false
	}

	right, err := r3.Add(b)
	if err != nil {
		return false
	}

	return left.Equal(right)
}

// 두 점이 서로 역원인지 확인하는 함수
func areInverse(x1, x2, y1, y2 FieldElement) bool {
	return x1.Equal(x2) && y1.NotEqual(y2)
}

// 두 타원곡선이 같은지 확인하는 함수
func sameCurve(a1, b1, a2, b2 FieldElement) bool {
	return a1.Equal(a2) && b1.Equal(b2)
}

// 두 점이 같은지 확인하는 함수
func samePoint(x1, y1, x2, y2 FieldElement) bool {
	return x1.Equal(x2) && y1.Equal(y2)
}
