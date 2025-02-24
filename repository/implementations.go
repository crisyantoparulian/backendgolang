package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	apperror "github.com/SawitProRecruitment/UserService/utils/app_error"
	"github.com/google/uuid"
)

type Relation string

const (
	RELATION_STATS Relation = "stats"
	RELATION_TREES Relation = "trees"
)

func (r *Repository) CreateEstate(ctx context.Context, input CreateEstateInput) (err error) {
	err = r.createEstateSql(ctx, input)
	if err != nil {
		return apperror.WrapWithCode(fmt.Errorf("failed to create estate: %w", err), http.StatusInternalServerError)
	}

	return
}

func (r *Repository) GetEstateWithAllDetails(ctx context.Context, id uuid.UUID, exludeRelations ...Relation) (estate *Estate, err error) {
	estate, err = r.getEstateByIdSql(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = apperror.WrapWithCode(fmt.Errorf("estate with ID %s not found", id), http.StatusNotFound)
			return nil, err
		}
		err = apperror.WrapWithCode(fmt.Errorf("failed to get estate: %w", err), http.StatusInternalServerError)
		return nil, err
	}

	// Convert exludeRelations slice to a map for quick lookup
	excludeMap := make(map[Relation]bool)
	for _, relation := range exludeRelations {
		excludeMap[relation] = true
	}

	if !excludeMap[RELATION_TREES] {
		estate.Trees, err = r.getTreesByEstateId(ctx, id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.WrapWithCode(fmt.Errorf("failed to get estate tress: %w", err), http.StatusInternalServerError)
		}
	}

	if !excludeMap[RELATION_STATS] {
		estate.Stats, err = r.getEstateStatsSQL(ctx, estate.Id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.WrapWithCode(fmt.Errorf("failed to get estate stats: %w", err), http.StatusInternalServerError)
		}
	}

	return estate, nil
}

func (r *Repository) CreateTree(ctx context.Context, input CreateTreeInput) (err error) {

	estate, err := r.getEstateByIdSql(ctx, input.EstateId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.WrapWithCode(fmt.Errorf("estate with ID %s not found", input.EstateId), http.StatusNotFound)
		}
		return apperror.WrapWithCode(fmt.Errorf("failed to get estate: %w", err), http.StatusInternalServerError)
	}

	// validate X & Y coordinate
	if input.X > estate.Length {
		return apperror.WrapWithCode(fmt.Errorf("coordinate x cannot greater than %d", estate.Length), http.StatusBadRequest)
	}
	if input.Y > estate.Width {
		return apperror.WrapWithCode(fmt.Errorf("coordinate y cannot greater than %d", estate.Width), http.StatusBadRequest)
	}

	isExist, err := r.checkExistEstateTree(ctx, CheckExistEstateTreeInput{
		EstateId: input.EstateId,
		X:        input.X,
		Y:        input.Y,
	})
	if err != nil {
		return apperror.WrapWithCode(fmt.Errorf("failed to check plot tree in database: %w", err), http.StatusInternalServerError)
	}

	if isExist {
		return apperror.WrapWithCode(errors.New("plot already has a tree"), http.StatusUnprocessableEntity)
	}

	err = r.createTreeSQL(ctx, input)
	if err != nil {
		return apperror.WrapWithCode(fmt.Errorf("failed to create tree: %w", err), http.StatusInternalServerError)
	}

	return
}

func (r *Repository) GetCalculatedEstateStats(ctx context.Context, estateId uuid.UUID) (stats *EstateStats, err error) {
	stats, err = r.getCalculatedEstateStatsSQL(ctx, estateId)
	if err != nil {
		err = apperror.WrapWithCode(fmt.Errorf("failed to get calculated stats: %w", err), http.StatusInternalServerError)
		return
	}
	return
}

func (r *Repository) UpsertEstateStats(ctx context.Context, estateID uuid.UUID, stats *EstateStats) error {
	err := r.upsertEstateStatsSQL(ctx, estateID, stats)
	if err != nil {
		return apperror.WrapWithCode(fmt.Errorf("failed to get upsert stats: %w", err), http.StatusInternalServerError)
	}

	return nil
}
