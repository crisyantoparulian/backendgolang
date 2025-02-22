package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	apperror "github.com/SawitProRecruitment/UserService/utils/app_error"
	"github.com/SawitProRecruitment/UserService/utils/array"
	"github.com/google/uuid"
)

func (r *Repository) CreateEstate(ctx context.Context, input CreateEstateInput) (err error) {
	err = r.createEstateSql(ctx, input)
	if err != nil {
		return apperror.WrapWithCode(fmt.Errorf("failed to create estate: %w", err), http.StatusInternalServerError)
	}

	return
}

func (r *Repository) CreateTree(ctx context.Context, input CreateTreeInput) (err error) {

	estate, err := r.getEstateByIdSql(ctx, input.EstateId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.WrapWithCode(fmt.Errorf("estate with ID %s not found", input.Id), http.StatusNotFound)
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

func (r *Repository) GetStatsEstate(ctx context.Context, estateID uuid.UUID) (stats GetStatsEstateOutput, err error) {

	// check estate exist
	_, err = r.getEstateByIdSql(ctx, estateID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = apperror.WrapWithCode(fmt.Errorf("estate with ID %s not found", estateID), http.StatusNotFound)
			return
		}
		err = apperror.WrapWithCode(fmt.Errorf("failed to get estate: %w", err), http.StatusInternalServerError)
		return
	}

	heights, err := r.getAllTreeHeightEstateSql(ctx, estateID)
	if err != nil {
		err = apperror.WrapWithCode(fmt.Errorf("failed to get all tree height: %w", err), http.StatusInternalServerError)
		return
	}

	// if there's no tree, then return all default value 0
	if len(heights) == 0 {
		return
	}

	// Calculate statistic
	count := len(heights)
	max := heights[count-1]
	min := heights[0]
	median := array.CalculateMedian(heights)

	return GetStatsEstateOutput{
		Count:  count,
		Max:    max,
		Min:    min,
		Median: median,
	}, nil
}

func (r *Repository) GetEstateByIdWithTrees(ctx context.Context, estateID uuid.UUID) (*Estate, error) {

	// get estate
	estate, err := r.getEstateByIdSql(ctx, estateID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.WrapWithCode(fmt.Errorf("estate with ID %s not found", estateID), http.StatusNotFound)
		}
		return nil, apperror.WrapWithCode(fmt.Errorf("failed to get estate: %w", err), http.StatusInternalServerError)
	}

	// get all tree in estate
	estate.Trees, err = r.getTreesByEstateId(ctx, estateID)
	if err != nil {
		return nil, apperror.WrapWithCode(fmt.Errorf("failed to get tress: %w", err), http.StatusInternalServerError)
	}

	return estate, nil
}
