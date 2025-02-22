package handler

import (
	"fmt"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	apperror "github.com/SawitProRecruitment/UserService/utils/app_error"
	httphelper "github.com/SawitProRecruitment/UserService/utils/http_helper"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Create a new estate
// (POST /estate)
func (s *Server) PostEstate(c echo.Context) error {
	ctx := c.Request().Context()
	payload := generated.CreateEstateRequest{}
	if err := c.Bind(&payload); err != nil {
		return httphelper.HttpRespError(c,
			apperror.WrapWithCode(fmt.Errorf("failed to unmarshall request: %w", err),
				http.StatusBadRequest))
	}

	// Validate payload
	if err := s.validatePayload(c, payload); err != nil {
		return err
	}

	estateID := uuid.New()
	err := s.Repository.CreateEstate(ctx, repository.CreateEstateInput{
		Id:     estateID,
		Width:  payload.Width,
		Length: payload.Length,
	})
	if err != nil {
		return httphelper.HttpRespError(c, err)
	}

	var resp generated.CreateEstateResponse
	openapiUUID := openapi_types.UUID(estateID)
	resp.Id = &openapiUUID

	return c.JSON(http.StatusCreated, resp)
}

// Get drone travel plan
// (GET /estate/{id}/drone-plan)
func (s *Server) GetEstateIdDronePlan(c echo.Context, id openapi_types.UUID, params generated.GetEstateIdDronePlanParams) error {
	ctx := c.Request().Context()

	// Validate payload
	if err := s.validatePayload(c, params); err != nil {
		return err
	}

	estate, err := s.Repository.GetEstateByIdWithTrees(ctx, id)
	if err != nil {
		return httphelper.HttpRespError(c, err)
	}

	distance := CalculateDroneDistance(estate)

	var resp generated.DronePlanResponse
	resp.Distance = &distance

	return c.JSON(http.StatusOK, resp)
}

// Get estate statistics
// (GET /estate/{id}/stats)
func (s *Server) GetEstateIdStats(c echo.Context, id openapi_types.UUID) error {
	ctx := c.Request().Context()

	var resp generated.EstateStats
	result, err := s.Repository.GetStatsEstate(ctx, id)
	if err != nil {
		return httphelper.HttpRespError(c, err)
	}

	resp = generated.EstateStats{
		Count:  &result.Count,
		Max:    &result.Max,
		Median: &result.Median,
		Min:    &result.Min,
	}

	return c.JSON(http.StatusOK, resp)
}

// Add a tree to an estate
// (POST /estate/{id}/tree)
func (s *Server) PostEstateIdTree(c echo.Context, id openapi_types.UUID) error {

	ctx := c.Request().Context()
	payload := generated.AddTreeRequest{}
	if err := c.Bind(&payload); err != nil {
		return httphelper.HttpRespError(c,
			apperror.WrapWithCode(fmt.Errorf("failed to unmarshall request: %w", err),
				http.StatusBadRequest))
	}

	// Validate payload
	if err := s.validatePayload(c, payload); err != nil {
		return err
	}

	treeId := uuid.New()
	err := s.Repository.CreateTree(ctx, repository.CreateTreeInput{
		Id:       treeId,
		EstateId: id,
		X:        payload.X,
		Y:        payload.Y,
		Height:   payload.Height,
	})
	if err != nil {
		return httphelper.HttpRespError(c, err)
	}

	var resp generated.CreateEstateResponse
	openapiUUID := openapi_types.UUID(treeId)
	resp.Id = &openapiUUID

	return c.JSON(http.StatusCreated, resp)
}
