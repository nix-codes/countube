package pic

// func countuneFn(n int, y0 int, d0 int) int {
// 	return parameterizedTriangularWaveFn(n-1, 10, 10, y0, d0)
// }

func countuneFn(x, y0, d0 int) int {
	return parameterizedTriangularWaveFn(x, 10, 10, y0, d0)
}

/*
** Parameterized triangular wave function.
**
** x: x coordinate (non-negative)
** s: step (number of steps before changing height of the wave; positive)
** A: amplitude (maximum height of the wave)
** y0: height of the wave at the origin x==0
** d0: direction of the wave at origin x==0 (0: up; 1: down)
**/
func parameterizedTriangularWaveFn(x, s, A, y0, d0 int) int {
	deltaX := y0
	if d0 == 1 {
		deltaX = 2*A - y0
	}

	x = x/s + deltaX
	// fmt.Println(deltaX, x, A)

	return canonicalTriangularWaveFn(x, A)
}

// func parameterizedTriangularWaveFn(x int, s int, A int, y0 int, d0 int) int {
// 	deltaX := y0
// 	if d0 == 1 {
// 		deltaX = 2*A - y0
// 	}

// 	x = x/s + deltaX
// 	fmt.Println(deltaX, x, A)

// 	return canonicalTriangularWaveFn(x, A)
// }

/*
** Canonical triangular wave function.
**
** x: x coordinate (non-negative)
** A: amplitude (maximum height; non-negative)
 */
func canonicalTriangularWaveFn(x int, A int) int {
	T := 2 * A // period
	x = x % T

	if x <= A {
		return x
	}

	return T - x
}
