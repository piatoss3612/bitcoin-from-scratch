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

	res1, err := p1.Add(*p2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res1)

	x3, _ := ecc.NewFieldElement(170, prime)
	y3, _ := ecc.NewFieldElement(142, prime)
	p3, _ := ecc.New(x3, y3, a, b)

	x4, _ := ecc.NewFieldElement(60, prime)
	y4, _ := ecc.NewFieldElement(139, prime)
	p4, _ := ecc.New(x4, y4, a, b)

	res2, err := p3.Add(*p4)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res2)

	x5, _ := ecc.NewFieldElement(47, prime)
	y5, _ := ecc.NewFieldElement(71, prime)
	p5, _ := ecc.New(x5, y5, a, b)

	x6, _ := ecc.NewFieldElement(17, prime)
	y6, _ := ecc.NewFieldElement(56, prime)
	p6, _ := ecc.New(x6, y6, a, b)

	res3, err := p5.Add(*p6)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res3)

	x7, _ := ecc.NewFieldElement(143, prime)
	y7, _ := ecc.NewFieldElement(98, prime)
	p7, _ := ecc.New(x7, y7, a, b)

	x8, _ := ecc.NewFieldElement(76, prime)
	y8, _ := ecc.NewFieldElement(66, prime)
	p8, _ := ecc.New(x8, y8, a, b)

	res4, err := p7.Add(*p8)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res4)
}
