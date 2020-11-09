package elliptic_crypto

import (
	"fmt"
	"math"
	"math/big"
	"sort"
)

var (
	test          = true
	M             = 67
	prime_numbers = []int{191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277}
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
	if test {
		M = 211
		return 0, -4
	}
	for b := 1; b < M; b++ {
		if (27*b*b+4)%M != 0 {
			return 1, b
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
	var c *Curve
	if test {
		c = &Curve{A: a, B: b, M: M, C: 241, Gx: 2, Gy: 2}
	} else {
		c = &Curve{A: a, B: b, M: M, C: 0, Gx: 0, Gy: 0}
	}
	allElements := c.BuildElements()
	fmt.Printf("List of all elliptic group elements: %v\n", allElements)
	fmt.Printf("Size of elliptic group: %v\n", len(allElements))
	if !test {
		for _, element := range allElements {
			if element.X == 0 || element.Y == 0 {
				continue
			}
			for _, i := range prime_numbers {
				multX, multY := c.ScalarMultiply(i, element.X, element.Y)
				if multX == c.Ox && multY == c.Oy {
					c.Gx = element.X
					c.Gy = element.Y
					c.C = i
					return c
				}
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
		tmp := new(big.Int).ModSqrt(bigYSqr, bigM)
		if tmp != nil {
			y1 := int(tmp.Int64())
			y2 := (-y1 + c.M) % c.M
			result = append(result, Element{X: i, Y: y1})
			if y1 != y2 {
				result = append(result, Element{X: i, Y: y2})
			}
		}
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
	diff := (3*x1*x1 + c.A) % c.M
	if diff < 0 {
		diff += c.M
	}
	m := big.NewInt(int64(diff))
	diff = 2 * y1 % c.M
	if diff == 0 {
		return c.Ox, c.Oy
	}
	n := big.NewInt(int64(diff))
	bigM := big.NewInt(int64(c.M))
	inv := new(big.Int).ModInverse(n, bigM)
	lambda := new(big.Int).Mul(inv, m)
	lambdaMod := int(new(big.Int).Mod(lambda, bigM).Int64())
	return c.applyLambda(lambdaMod, x1, y1, x1)
}

func (c *Curve) sumTwoElements(x1, y1, x2, y2 int) (int, int) {
	if x1 == c.Ox && y1 == c.Ox {
		return x2, y2
	}
	if x2 == c.Ox && y2 == c.Ox {
		return x1, y1
	}
	diff := y2 - y1
	if diff < 0 {
		diff += c.M
	}
	y2y1 := big.NewInt(int64(diff))
	diff = x2 - x1
	if diff == 0 {
		return c.Ox, c.Oy
	}
	if diff < 0 {
		diff += c.M
	}
	x2x1 := big.NewInt(int64(diff))
	bigM := big.NewInt(int64(c.M))
	inv := new(big.Int).ModInverse(x2x1, bigM)
	lambda := new(big.Int).Mul(inv, y2y1)
	lambdaMod := int(new(big.Int).Mod(lambda, bigM).Int64())
	return c.applyLambda(lambdaMod, x1, y1, x2)
}

func (c *Curve) applyLambda(lambda, x1, y1, x2 int) (int, int) {
	x := (lambda*lambda - x1 - x2) % c.M
	if x < 0 {
		x += c.M
	}
	y := (-y1 + lambda*(-x+x1)) % c.M
	if y < 0 {
		y += c.M
	}
	return x, y
}

func (c *Curve) ScalarMultiply(scalar, x, y int) (int, int) {
	//bitsArray := bits(scalar - 1)
	resX, resY := c.addSameElement(x, y)
	for i := 2; i < scalar; i++ {
		if resX == x && resY == y {
			resX, resY = c.addSameElement(x, y)
		}
		resX, resY = c.sumTwoElements(resX, resY, x, y)
	}
	/*for _, bit := range bitsArray {
		if bit == 1 && resY != y {
			x, y = c.sumTwoElements(resX, resY, x, y)
		}
		resX, resY = c.addSameElement(resX, resY)
	}*/
	return resX, resY
}

func (c *Curve) ScalarOriginMultiply(scalar int) (int, int) {
	resX, resY := c.addSameElement(c.Gx, c.Gy)
	for i := 2; i < scalar; i++ {
		if resX == c.Gx && resY == c.Gy {
			resX, resY = c.addSameElement(c.Gx, c.Gy)
		}
		resX, resY = c.sumTwoElements(resX, resY, c.Gx, c.Gy)
	}
	return resX, resY
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

func findDividers(x int) []int {
	factors := []int{x}
	for i := 2; i < int(math.Sqrt(float64(x))); i++ {
		if x%i == 0 {
			factors = append(factors, i)
			if int(x/i) != i {
				factors = append(factors, int(x/i))
			}
		}
	}
	sort.Slice(factors, func(i, j int) bool { return factors[i] < factors[j] })
	return factors
}
