package ecc

import "math/big"

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
