package elliptic_crypto

import (
	"math"
	"math/big"
)

var (
	M = 67
)

type Element struct {
	X int
	Y int
}

type Curve struct {
	A      int
	B      int
	M      int
	C      int
	Gx, Gy int
	Ox, Oy int
}

type CurveOperations interface {
	IsOnCurve(x, y int) bool
	Add(x1, y1, x2, y2 int) (int, int)
	ScalarMultiply(scalar, x, y int) (int, int)
	ScalarOriginMultiply(scalar int) (int, int)
	BuildElements() []Element
}

func findSuitableAandB() (int, int) {
	for b := 0; b < M; b++ {
		if 27*b*b%M != 0 {
			return 0, b
		}
		if (4*b*b*b+27)%M != 0 {
			return b, 1
		}
	}
	return 0, 0
}

func isPrime(x int) bool {
	until := int(math.Sqrt(float64(x)))
	for i := 2; i <= until; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func BuildCurve() *Curve {
	a, b := findSuitableAandB()
	c := &Curve{A: a, B: b, M: M, C: 0, Gx: 0, Gy: 0}
	allElements := c.BuildElements()
	firstElement := allElements[1]
	negFirstElement := Element{X: firstElement.X, Y: -firstElement.Y + c.M}
	c.Ox, c.Oy = c.sumTwoElements(firstElement.X, firstElement.Y, negFirstElement.X, negFirstElement.Y)
	for _, element := range allElements {
		for i := 23; i < c.M; i++ {
			multX, multY := c.ScalarMultiply(i, element.X, element.Y)
			if multX == c.Ox && multY == c.Oy && isPrime(i) {
				c.Gx = element.X
				c.Gy = element.Y
				c.C = i
				return c
			}
		}
	}
	return c
}

func (c *Curve) BuildElements() []Element {
	result := []Element{}
	for i := 0; i < c.M; i++ {
		ySqr := (i*i*i + c.A*i + c.B) % c.M
		bigYSqr := big.NewInt(int64(ySqr))
		bigM := big.NewInt(int64(c.M))
		y1 := int(bigYSqr.ModSqrt(bigYSqr, bigM).Int64())
		y2 := -y1 + c.M
		result = append(result, Element{X: i, Y: y1}, Element{X: i, Y: y2})
	}
	return result
}

func (c *Curve) Add(x1, y1, x2, y2 int) (int, int) {
	var x, y int
	if x1 == x2 && y1 == y2 {
		x, y = c.addSameElement(x1, y1)
	} else {
		x, y = c.sumTwoElements(x1, y1, x2, y2)
	}
	return x, y
}

func (c *Curve) addSameElement(x1, y1 int) (int, int) {
	m := big.NewInt(int64((3*x1*x1 + c.A) % c.M))
	n := big.NewInt(int64((2 * y1) % c.M))
	bigM := big.NewInt(int64(c.M))
	inv := new(big.Int).ModInverse(n, bigM)
	lambda := int(new(big.Int).Mul(inv, m).Int64())
	x := (lambda*lambda - x1 - x1) % c.M
	y := (-y1 + lambda*(x-x1)) % M
	return x, y
}

func (c *Curve) sumTwoElements(x1, y1, x2, y2 int) (int, int) {
	y2y1 := big.NewInt(int64(y2 - y1))
	x2x1 := big.NewInt(int64(x2 - x1))
	bigM := big.NewInt(int64(c.M))
	inv := new(big.Int).ModInverse(y2y1, bigM)
	lambda := int(new(big.Int).Mul(inv, x2x1).Int64())
	x := (lambda*lambda - x1 - x2) % c.M
	y := (-y1 + lambda*(x-x1)) % M
	return x, y
}

func (c *Curve) ScalarMultiply(scalar, x, y int) (int, int) {
	bitsArray := bits(scalar)
	resX, resY := x, y
	for _, bit := range bitsArray {
		if bit == 1 {
			x, y = c.sumTwoElements(resX, resY, x, y)
		}
		resX, resY = c.addSameElement(resX, resY)
	}
	return resX, resY
}

func (c *Curve) ScalarOriginMultiply(scalar int) (int, int) {
	x, y := c.Gx, c.Gy
	for i := 0; i < scalar; i++ {
		x, y = c.addSameElement(x, y)
	}
	return x, y
}

func (c *Curve) IsOnCurve(x, y int) bool {
	return y*y == x*x*x+c.A*x+c.B
}

func bits(n int) []int {
	res := []int{}
	for n > 0 {
		res = append(res, n%2)
		n = n >> 1
	}
	return res
}
