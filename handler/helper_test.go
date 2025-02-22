package handler_test

import (
	"testing"

	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/stretchr/testify/assert"
)

func TestCalculateDroneDistance(t *testing.T) {

	tests := []struct {
		name     string
		estate   *repository.Estate
		expected int
	}{
		{
			name: "Estate with different tree heights (from test example)",
			estate: &repository.Estate{
				Width:  1,
				Length: 5,
				Trees: []repository.Tree{
					{X: 2, Y: 1, Height: 5},
					{X: 3, Y: 1, Height: 3},
					{X: 4, Y: 1, Height: 4},
				},
			},
			expected: 54, // Adjust expected value as needed
		},
		// {
		// 	name: "No trees, flat estate",
		// 	estate: &repository.Estate{
		// 		Width:  3,
		// 		Length: 3,
		// 		Trees:  []repository.Tree{},
		// 	},
		// 	expected: 82, // Manually calculated
		// },
		// {
		// 	name: "Estate with with different tree heights (2x2)",
		// 	estate: &repository.Estate{
		// 		Width:  2,
		// 		Length: 2,
		// 		Trees: []repository.Tree{
		// 			{X: 1, Y: 1, Height: 5},
		// 			{X: 2, Y: 2, Height: 3},
		// 		},
		// 	},
		// 	expected: 42, // Adjust expected value as needed
		// },
		// {
		// 	name: "Estate with only 1 tree plot",
		// 	estate: &repository.Estate{
		// 		Width:  1,
		// 		Length: 1,
		// 		Trees: []repository.Tree{
		// 			{X: 1, Y: 1, Height: 10},
		// 		},
		// 	},
		// 	expected: 22, // Manually calculated
		// },
		// {
		// 	name: "Estate with only plot",
		// 	estate: &repository.Estate{
		// 		Width:  1,
		// 		Length: 1,
		// 		Trees:  []repository.Tree{},
		// 	},
		// 	expected: 2, // Manually calculated
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.CalculateDroneDistance(tt.estate)
			assert.Equal(t, tt.expected, result)
		})
	}
}
