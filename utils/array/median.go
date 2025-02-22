package array

import "sort"

// CalculateMedian returns the median of an array of integers.
// If the input slice is empty, it returns 0.
func CalculateMedian(ints []int) int {
	if len(ints) == 0 {
		return 0 // Prevents index out of range panic
	}

	sort.Ints(ints) // Ensure ordering
	n := len(ints)
	if n%2 == 0 {
		return (ints[n/2-1] + ints[n/2]) / 2
	}
	return ints[n/2]
}
