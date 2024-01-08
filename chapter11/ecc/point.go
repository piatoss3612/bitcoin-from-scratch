package ecc

import (
	"chapter11/utils"
	"fmt"
	"math/big"
)

// 타원곡선의 점을 생성하는 함수 타입
type PointGenerator func(x, y, a, b FieldElement) (Point, error)

// 타원곡선의 점 구조체
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

// 타원곡선의 점의 x좌표를 반환하는 함수
func (p point) X() FieldElement {
	return p.x
}

// 타원곡선의 점의 y좌표를 반환하는 함수
func (p point) Y() FieldElement {
	return p.y
}

// 타원곡선의 a 계수를 반환하는 함수
func (p point) A() FieldElement {
	return p.a
}

// 타원곡선의 b 계수를 반환하는 함수
func (p point) B() FieldElement {
	return p.b
}

// 타원곡선의 점을 문자열로 표현하는 함수 (Stringer 인터페이스 구현)
func (p point) String() string {
	// 무한원점인지 확인
	if isInfinity(p.x, p.y) {
		return "Point(infinity)"
	}
	return fmt.Sprintf("Point(%s, %s)_%s_%s FieldElement(%s)",
		p.x.Num().Text(16), p.y.Num().Text(16), p.a.Num().Text(16), p.b.Num().Text(16), p.x.Prime().Text(16))
}

// 두 타원곡선의 점이 같은지 확인하는 함수
func (p point) Equal(other Point) bool {
	return p.equal(other)
}

// 두 타원곡선의 점이 다른지 확인하는 함수
func (p point) NotEqual(other Point) bool {
	return !p.equal(other)
}

// 두 타원곡선의 점이 같은지 확인하는 내부 함수
func (p point) equal(other Point) bool {
	return samePoint(p.x, p.y, other.X(), other.Y()) &&
		sameCurve(p.a, p.b, other.A(), other.B())
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
		return p.addDifferentPoint(other, NewPoint)
	}

	/* case3: 두 점이 같은 경우 */

	// p와 other가 같은 점인지 확인
	if samePoint(p.x, p.y, other.X(), other.Y()) {
		return p.addSamePoint(other, NewPoint)
	}

	return nil, fmt.Errorf("unhandled case, (%s, %s) + (%s, %s)", p.x, p.y, other.X(), other.Y())
}

// 서로 다른 두 타원곡선의 점을 더하는 내부 함수
func (p point) addDifferentPoint(other Point, gen PointGenerator) (Point, error) {
	// p와 other를 지나는 직선의 기울기 구하기
	prime := p.x.Prime()

	// s = (other.y - p.y) * (other.x - p.x)^-1 % prime
	s := big.NewInt(0).Mod(
		// (other.y - p.y) * (other.x - p.x)^-1
		big.NewInt(0).Mul(
			// (other.y - p.y) % prime
			big.NewInt(0).Mod(big.NewInt(0).Sub(other.Y().Num(), p.y.Num()), prime),
			// * (other.x - p.x)^-1 % prime
			big.NewInt(0).ModInverse(big.NewInt(0).Sub(other.X().Num(), p.x.Num()), prime),
		),
		prime,
	)

	// p와 other를 지나는 직선이 타원곡선과 만나는 다른 한 점 q의 좌표 구하기

	// nx = s^2 - p.x - other.x % prime
	nx := big.NewInt(0).Mod(
		big.NewInt(0).Sub(
			big.NewInt(0).Sub(
				big.NewInt(0).Exp(s, big.NewInt(2), nil),
				p.x.Num(),
			),
			other.X().Num(),
		),
		prime,
	)

	// ny = (s * (p.x - nx) - p.y) % prime
	ny := big.NewInt(0).Mod(
		big.NewInt(0).Sub(
			big.NewInt(0).Mod(
				big.NewInt(0).Mul(
					s,
					big.NewInt(0).Sub(p.x.Num(), nx),
				),
				prime,
			),
			p.y.Num(),
		),
		prime,
	)

	nxe, err := NewFieldElement(nx, prime)
	if err != nil {
		return nil, err
	}

	nye, err := NewFieldElement(ny, prime)
	if err != nil {
		return nil, err
	}

	return gen(nxe, nye, p.a, p.b)
}

// 동일한 타원곡선의 점을 더하는 내부 함수
func (p point) addSamePoint(other Point, gen PointGenerator) (Point, error) {
	// case 2-1 예외 처리: 접선이 x축에 수직인 경우, 무한원점을 반환
	if p.y.Num().Cmp(big.NewInt(0)) == 0 {
		return nil, nil
	}
	// 접선의 기울기 구하기
	prime := p.x.Prime()

	// s = (3 * p.x^2 + p.a) * (2 * p.y)^-1 % prime
	s := big.NewInt(0).Mod(
		big.NewInt(0).Mul(
			big.NewInt(0).Add(
				big.NewInt(0).Mul(
					big.NewInt(3),
					big.NewInt(0).Exp(p.x.Num(), big.NewInt(2), prime),
				),
				p.a.Num(),
			),
			big.NewInt(0).ModInverse(
				big.NewInt(0).Mul(
					big.NewInt(2),
					p.y.Num(),
				),
				prime,
			),
		),
		prime,
	)

	// 접선과 타원곡선의 교점 q의 좌표 구하기

	// nx = (s^2 - 2 * p.x) % prime
	nx := big.NewInt(0).Mod(
		big.NewInt(0).Sub(
			big.NewInt(0).Exp(s, big.NewInt(2), prime),
			big.NewInt(0).Mod(
				big.NewInt(0).Mul(
					big.NewInt(2),
					p.x.Num(),
				),
				prime,
			),
		),
		prime,
	)

	// ny = (s * (p.x - nx) - p.y) % prime
	ny := big.NewInt(0).Mod(
		big.NewInt(0).Sub(
			big.NewInt(0).Mod(
				big.NewInt(0).Mul(
					s,
					big.NewInt(0).Sub(p.x.Num(), nx),
				),
				prime,
			),
			p.y.Num(),
		),
		prime,
	)

	nxe, err := NewFieldElement(nx, prime)
	if err != nil {
		return nil, err
	}

	nye, err := NewFieldElement(ny, prime)
	if err != nil {
		return nil, err
	}

	return gen(nxe, nye, p.a, p.b)
}

// 타원곡선 점의 스칼라 곱셈 함수
func (p point) Mul(coefficient *big.Int) (Point, error) {
	current := Point(&p) // 시작점으로 초기화

	result, err := NewPoint(nil, nil, p.a, p.b)
	if err != nil {
		return nil, err
	}

	return p.mul(coefficient, current, result)
}

// 타원곡선 점의 스칼라 곱셈 내부 함수
func (p point) mul(coefficient *big.Int, current, result Point) (Point, error) {
	coef := coefficient // 계수

	var err error

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

// 타원곡선 점의 서명 검증 함수
func (p point) Verify(z []byte, sig Signature) (bool, error) {
	// TODO: implement verify
	return false, nil
}

// 타원곡선 점의 직렬화 함수
func (p point) SEC(compressed bool) []byte {
	// TODO: implement SEC
	return nil
}

// 타원곡선 점의 Hash160 함수
func (p point) Hash160(compressed bool) []byte {
	// TODO: implement Hash160
	return nil
}

// 타원곡선 점의 주소 생성 함수
func (p point) Address(compressed bool, testnet bool) string {
	// TODO: implement Address
	return ""
}

// secp256k1 타원곡선의 점 구조체
type s256Point struct {
	point // 상위 구조체를 임베딩하여 기능 상속, 필드 재사용
}

// secp256k1 타원곡선의 점 생성 함수
func NewS256Point(x, y FieldElement) (Point, error) {
	a, err := NewS256FieldElement(big.NewInt(int64(A)))
	if err != nil {
		return nil, err
	}

	b, err := NewS256FieldElement(big.NewInt(int64(B)))
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

// secp256k1 타원곡선의 점을 더하는 PointGenerator 함수
func (p s256Point) Add(other Point) (Point, error) {
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
		return NewS256Point(nil, nil)
	}

	// secp256k1 타원곡선의 점을 생성하는 함수
	gen := func(x, y, _, _ FieldElement) (Point, error) {
		return NewS256Point(x, y)
	}

	/* case2: 두 점이 서로 다른 경우 */

	if p.x.NotEqual(other.X()) {
		return p.addDifferentPoint(other, gen)
	}

	/* case3: 두 점이 같은 경우 */

	// p와 other가 같은 점인지 확인
	if samePoint(p.x, p.y, other.X(), other.Y()) {
		return p.addSamePoint(other, gen)
	}

	return nil, fmt.Errorf("unhandled case, (%s, %s) + (%s, %s)", p.x, p.y, other.X(), other.Y())
}

// secp256k1 타원곡선의 점의 스칼라 곱셈 함수
func (p s256Point) Mul(coefficient *big.Int) (Point, error) {
	// 계수가 N보다 큰 경우, 계수를 N으로 나눈 나머지를 계수로 사용
	// 왜? N*G = O, 무한원점 이기 때문에 N보다 큰 계수는 N으로 나눈 나머지를 계수로 사용해도 결과는 같음
	coef := new(big.Int).Mod(coefficient, N)
	current := Point(&p) // 시작점으로 초기화

	result, err := NewS256Point(nil, nil)
	if err != nil {
		return nil, err
	}

	return p.mul(coef, current, result)
}

// secp256k1 타원곡선의 점의 서명 검증 함수
func (p s256Point) Verify(z []byte, sig Signature) (bool, error) {
	bigZ := utils.BytesToBigInt(z) // z를 big.Int로 변환
	sInv := invBN(sig.S(), N)      // s^-1
	u := mulBN(bigZ, sInv, N)      // u = z * s^-1
	v := mulBN(sig.R(), sInv, N)   // v = r * s^-1

	uG, err := G.Mul(u) // uG
	if err != nil {
		return false, err
	}

	vP, err := p.Mul(v) // vP
	if err != nil {
		return false, err
	}

	R, err := uG.Add(vP) // uG + vP
	if err != nil {
		return false, err
	}

	x := R.X().Num() // res의 x좌표

	return x.Cmp(sig.R()) == 0, nil // x좌표가 r과 같은지 확인
}

// secp256k1 타원곡선의 점의 직렬화 함수
func (p s256Point) SEC(compressed bool) []byte {
	if compressed {
		// y좌표의 LSB가 0인 경우, 0x02를 prefix로 사용
		if p.y.Num().Bit(0) == 0 {
			return append([]byte{0x02}, p.x.Num().FillBytes(make([]byte, 32))...)
		}
		// y좌표의 LSB가 1인 경우, 0x03을 prefix로 사용
		return append([]byte{0x03}, p.x.Num().FillBytes(make([]byte, 32))...)
	}

	return append([]byte{0x04},
		append(
			p.x.Num().FillBytes(make([]byte, 32)),
			p.y.Num().FillBytes(make([]byte, 32))...)...)
}

// secp256k1 타원곡선의 점의 SEC 형식을 160비트 해시로 변환하는 함수
func (p s256Point) Hash160(compressed bool) []byte {
	return utils.Hash160(p.SEC(compressed))
}

// secp256k1 타원곡선의 점을 주소로 변환하는 함수
func (p s256Point) Address(compressed bool, testnet bool) string {
	h160 := p.Hash160(compressed) // 타원곡선 점의 SEC 형식을 160비트 해시로 변환

	if testnet {
		h160 = append([]byte{0x6f}, h160...) // testnet 주소의 prefix는 0x6f
	} else {
		h160 = append([]byte{0x00}, h160...) // mainnet 주소의 prefix는 0x00
	}

	return utils.EncodeBase58Checksum(h160) // Base58Checksum 인코딩
}
