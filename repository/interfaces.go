// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/google/uuid"
)

type RepositoryInterface interface {
	CreateEstate(ctx context.Context, input CreateEstateInput) (err error)
	CreateTree(ctx context.Context, input CreateTreeInput) (err error)
	GetEstateWithAllDetails(ctx context.Context, id uuid.UUID, exludeRelations ...Relation) (estate *Estate, err error)
	GetCalculatedEstateStats(ctx context.Context, estateId uuid.UUID) (stats *EstateStats, err error)
	UpsertEstateStats(ctx context.Context, estateID uuid.UUID, stats *EstateStats) error
}
