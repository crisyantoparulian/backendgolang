// This file contains types that are used in the repository layer.
package repository

import (
	"time"

	"github.com/google/uuid"
)

type GetStatsEstateOutput struct {
	Count  int
	Max    int
	Min    int
	Median int
}

type CreateEstateInput struct {
	Id     uuid.UUID
	Width  int
	Length int
}

type CreateTreeInput struct {
	Id       uuid.UUID
	EstateId uuid.UUID
	X        int
	Y        int
	Height   int
}

type CheckExistEstateTreeInput struct {
	EstateId uuid.UUID
	X        int
	Y        int
}

type Estate struct {
	Id        uuid.UUID
	Width     int
	Length    int
	CreatedAt time.Time
	UpdatedAt *time.Time
	Trees     []Tree
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
