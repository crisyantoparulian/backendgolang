package handler

import (
	"testing"

	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/stretchr/testify/assert"
)

func TestCalculateDroneDistance(t *testing.T) {

	tests := []struct {
		name               string
		estate             *repository.Estate
		maxDistance        *int
		expectedDistance   int64
		expectedCoordinate *Coordinate
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
			expectedDistance: 54,
		},
		{
			name: "Estate with different tree heights, empty plot at (1,1) & (4,5)",
			estate: &repository.Estate{
				Width:  1,
				Length: 6,
				Trees: []repository.Tree{
					{X: 2, Y: 1, Height: 5},
					{X: 3, Y: 1, Height: 3},
					{X: 4, Y: 1, Height: 4},
					{X: 6, Y: 1, Height: 6},
				},
			},
			expectedDistance: 68,
		},
		{
			name: "No trees, flat estate",
			estate: &repository.Estate{
				Width:  3,
				Length: 3,
				Trees:  []repository.Tree{},
			},
			expectedDistance: 82,
		},
		{
			name: "Estate with with different tree heights (2x2)",
			estate: &repository.Estate{
				Width:  2,
				Length: 2,
				Trees: []repository.Tree{
					{X: 1, Y: 1, Height: 5},
					{X: 2, Y: 2, Height: 3},
				},
			},
			expectedDistance: 42,
		},
		{
			name: "Estate with only 1 tree plot",
			estate: &repository.Estate{
				Width:  1,
				Length: 1,
				Trees: []repository.Tree{
					{X: 1, Y: 1, Height: 10},
				},
			},
			expectedDistance: 22,
		},
		{
			name: "Estate with only plot without tree",
			estate: &repository.Estate{
				Width:  1,
				Length: 1,
				Trees:  []repository.Tree{},
			},
			expectedDistance: 2,
		},
		{
			name: "Estate with different tree heights (from test example) and max_distance provided (landed at first plot)",
			estate: &repository.Estate{
				Width:  1,
				Length: 5,
				Trees: []repository.Tree{
					{X: 2, Y: 1, Height: 5},
					{X: 3, Y: 1, Height: 3},
					{X: 4, Y: 1, Height: 4},
				},
			},
			maxDistance:        pointerInt(5),
			expectedCoordinate: &Coordinate{X: 1, Y: 1},
			expectedDistance:   54,
		},
		{
			name: "Estate with different tree heights (from test example) and max_distance provided (landed at last plot)",
			estate: &repository.Estate{
				Width:  1,
				Length: 5,
				Trees: []repository.Tree{
					{X: 2, Y: 1, Height: 5},
					{X: 3, Y: 1, Height: 3},
					{X: 4, Y: 1, Height: 4},
				},
			},
			maxDistance:        pointerInt(54),
			expectedCoordinate: &Coordinate{X: 5, Y: 1},
			expectedDistance:   54,
		},
		{
			name: "Estate with different tree heights (from test example) and max_distance provided (landed at plot 3,1)",
			estate: &repository.Estate{
				Width:  1,
				Length: 5,
				Trees: []repository.Tree{
					{X: 2, Y: 1, Height: 5},
					{X: 3, Y: 1, Height: 3},
					{X: 4, Y: 1, Height: 4},
				},
			},
			maxDistance:        pointerInt(30),
			expectedCoordinate: &Coordinate{X: 2, Y: 1},
			expectedDistance:   54,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance, lastCoordinate := calculateDroneDistance(tt.estate, tt.maxDistance)
			assert.Equal(t, tt.expectedDistance, distance)
			assert.Equal(t, tt.expectedCoordinate, lastCoordinate)
		})
	}
}

func pointerInt(i int) *int {
	return &i
}
