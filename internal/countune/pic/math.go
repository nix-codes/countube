package pic

/*
** The Countune function is just a parameterized triangular stepped wave function.
**
** x: x coordinate (non-negative)
** s: step (number of steps before changing height of the wave; positive)
** A: amplitude (maximum height of the wave)
** x0: coordinate of the origin (as shifted from x==0)
** y0: height of the wave at the origin
** d0: direction of the wave at origin (0: up; 1: down)
**/
func countuneFn(x, s, A, x0, y0, d0 int) int {
	x = x - x0
	deltaX := y0
	if d0 == 1 {
		deltaX = 2*A - y0
	}

	x = x/s + deltaX

	return canonicalTriangularWaveFn(x, A)
}

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
