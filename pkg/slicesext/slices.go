package slicesext

import (
	"cmp"
	"slices"
)

// Sum calculates the sum of items
// in the slice s
func Sum[T cmp.Ordered](s []T) T {
	var sum T

	for _, v := range s {
		sum += v
	}

	return sum
}

// Intersection intersects on s0 and s1
func Intersection[T comparable](s0, s1 []T) []T {
	res := make([]T, 0)

	for _, v := range s0 {
		if slices.Contains(s1, v) {
			res = append(res, v)
		}
	}

	return res
}
