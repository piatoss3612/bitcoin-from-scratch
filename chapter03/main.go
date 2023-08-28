package main

import (
	"chapter03/ecc"
	"fmt"
)

func main() {
	prime := 223
	a, _ := ecc.NewFieldElement(0, prime)
	b, _ := ecc.NewFieldElement(7, prime)
	x1, _ := ecc.NewFieldElement(192, prime)
	y1, _ := ecc.NewFieldElement(105, prime)
	x2, _ := ecc.NewFieldElement(17, prime)
	y2, _ := ecc.NewFieldElement(56, prime)

	p1, _ := ecc.New(x1, y1, a, b)
	p2, _ := ecc.New(x2, y2, a, b)
	fmt.Println(p1, p2)

	res, err := p1.Add(*p2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
