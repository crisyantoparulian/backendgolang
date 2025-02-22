package array_test

import (
	"testing"

	"github.com/SawitProRecruitment/UserService/utils/array"
)

func TestCalculateMedian(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"odd number of elements", []int{1, 3, 2}, 2},
		{"even number of elements", []int{1, 2, 3, 4}, 2},
		{"unsorted array", []int{5, 3, 1, 4, 2}, 3},
		{"negative numbers", []int{-5, -3, -1, -4, -2}, -3},
		{"all same numbers", []int{7, 7, 7, 7, 7}, 7},
		{"single element", []int{42}, 42},
		{"two elements", []int{10, 20}, 15},
		{"empty slice", []int{}, 0}, // Edge case
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := array.CalculateMedian(tt.nums)
			if got != tt.want {
				t.Errorf("CalculateMedian(%v) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
}
