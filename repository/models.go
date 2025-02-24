// This file contains types that are used in the repository layer. for representation table in database
package repository

import (
	"time"

	"github.com/google/uuid"
)

type Estate struct {
	Id     uuid.UUID
	Width  int
	Length int

	CreatedAt time.Time
	UpdatedAt *time.Time
	Trees     []Tree
	Stats     *EstateStats
}

type EstateStats struct {
	Id            uuid.UUID
	EstateID      string
	TreeCount     int64
	MaxHeight     int
	MinHeight     int
	MedianHeight  int
	DroneDistance int64
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type Tree struct {
	Id        uuid.UUID
	EstateId  uuid.UUID
	X         int
	Y         int
	Height    int
	CreatedAt time.Time
	UpdatedAt *time.Time
}
