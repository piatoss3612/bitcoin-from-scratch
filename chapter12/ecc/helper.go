package ecc

import (
	"bytes"
	"errors"
	"math/big"
)

// 두 원소값이 같은지 확인하는 함수
func sameBN(x, y *big.Int) bool {
	return x.Cmp(y) == 0
}

// 두 원소값을 더하는 함수
func addBN(x, y, mod *big.Int) *big.Int {
	return big.NewInt(0).Mod(big.NewInt(0).Add(x, y), mod)
}

// 두 원소값을 빼는 함수
func subBN(x, y, mod *big.Int) *big.Int {
	return big.NewInt(0).Mod(big.NewInt(0).Sub(x, y), mod)
}

// 두 원소값을 곱하는 함수
func mulBN(x, y, mod *big.Int) *big.Int {
	return big.NewInt(0).Mod(big.NewInt(0).Mul(x, y), mod)
}

// 원소값을 거듭제곱하는 함수
func powBN(x, exp, mod *big.Int) *big.Int {
	return big.NewInt(0).Exp(x, exp, mod)
}

// 원소값의 역원을 구하는 함수
func invBN(x, mod *big.Int) *big.Int {
	return big.NewInt(0).ModInverse(x, mod)
}

// 원소값의 제곱근을 구하는 함수
func sqrtBN(x, mod *big.Int) *big.Int {
	return big.NewInt(0).ModSqrt(x, mod)
}

// 무한원점인지 확인하는 함수
func isInfinity(x, y FieldElement) bool {
	return x == nil && y == nil
}

// 타원곡선 위에 있는지 확인하는 함수
func isOnCurve(x, y, a, b FieldElement) bool {
	prime := x.Prime()

	// y^2 == x^3 + ax + b
	left := powBN(y.Num(), big.NewInt(2), prime)
	right := addBN(
		addBN(
			powBN(x.Num(), big.NewInt(3), prime),
			mulBN(a.Num(), x.Num(), prime),
			prime,
		),
		b.Num(),
		prime,
	)

	return sameBN(left, right)
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

// num이 0보다 크거나 같고 prime보다 작은지 확인하는 함수
func inRange(num, prime *big.Int) bool {
	return num.Cmp(big.NewInt(0)) != -1 && num.Cmp(prime) == -1
}

// sec 바이트 슬라이스를 타원곡선 위의 점으로 변환하는 함수
func ParsePoint(sec []byte) (Point, error) {
	// prefix가 0x04인 경우, 비압축 포맷
	if sec[0] == 0x04 {
		x, err := NewS256FieldElement(new(big.Int).SetBytes(sec[1:33]))
		if err != nil {
			return nil, err
		}

		y, err := NewS256FieldElement(new(big.Int).SetBytes(sec[33:65]))
		if err != nil {
			return nil, err
		}

		return NewS256Point(x, y)
	}

	// prefix가 0x02 또는 0x03인 경우, 압축 포맷
	if sec[0] == 0x02 || sec[0] == 0x03 {
		x, err := NewS256FieldElement(new(big.Int).SetBytes(sec[1:]))
		if err != nil {
			return nil, err
		}

		// y^2 = x^3 + 7
		alpha := addBN(powBN(x.Num(), big.NewInt(3), P), big.NewInt(int64(B)), P)
		// y = sqrt(alpha)
		beta := sqrtBN(alpha, P)

		var even, odd *big.Int

		// y의 LSB가 짝수인지 홀수인지 확인
		if byte(beta.Bit(0)) == 0x00 {
			even = beta
			odd = subBN(P, beta, P)
		} else {
			odd = beta
			even = subBN(P, beta, P)
		}

		// prefix가 0x02인 경우, y의 LSB가 짝수인 값을 사용
		if sec[0] == 0x02 {
			y, err := NewS256FieldElement(even)
			if err != nil {
				return nil, err
			}

			return NewS256Point(x, y)
		}

		// prefix가 0x03인 경우, y의 LSB가 홀수인 값을 사용
		y, err := NewS256FieldElement(odd)
		if err != nil {
			return nil, err
		}

		return NewS256Point(x, y)
	}

	return nil, errors.New("invalid sec format")
}

// 서명을 파싱하는 함수
func ParseSignature(der []byte) (Signature, error) {
	buf := bytes.NewBuffer(der) // der을 읽기 위한 버퍼 생성

	compound := buf.Next(1)[0]

	// 서명이 0x30으로 시작하는지 확인
	if compound != 0x30 {
		return nil, errors.New("invalid compound byte")
	}

	length := buf.Next(1)[0] // 서명의 길이

	// 서명의 길이가 der의 길이와 일치하는지 확인
	if int(length)+2 > len(der) {
		return nil, errors.New("invalid length")
	}

	marker := buf.Next(1)[0] // r의 시작을 나타내는 마커
	if marker != 0x02 {      // r의 시작을 나타내는 마커는 0x02 이어야 함
		return nil, errors.New("invalid marker for r")
	}

	rLength := buf.Next(1)[0]   // r의 길이
	r := buf.Next(int(rLength)) // r

	marker = buf.Next(1)[0] // s의 시작을 나타내는 마커
	if marker != 0x02 {     // s의 시작을 나타내는 마커는 0x02 이어야 함
		return nil, errors.New("invalid marker for s")
	}

	sLength := buf.Next(1)[0]   // s의 길이
	s := buf.Next(int(sLength)) // s

	// // 서명의 길이가 der의 길이와 일치하는지 확인
	// if len(der) < 6+int(rLength)+int(sLength) {
	// 	return nil, errors.New("invalid signature")
	// }

	// S256Signature 타입의 인스턴스를 생성하여 반환
	return NewS256Signature(
		big.NewInt(0).SetBytes(r),
		big.NewInt(0).SetBytes(s),
	), nil
}
