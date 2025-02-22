// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type RepositoryInterface interface {
	CreateEstate(ctx context.Context, input CreateEstateInput) (err error)

	CreateTree(ctx context.Context, input CreateTreeInput) (err error)
	GetStatsEstate(ctx context.Context, estateID openapi_types.UUID) (stats GetStatsEstateOutput, err error)
	GetEstateByIdWithTrees(ctx context.Context, estateID uuid.UUID) (*Estate, error)
}
