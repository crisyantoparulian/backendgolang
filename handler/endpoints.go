package handler

import (
	"fmt"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	apperror "github.com/SawitProRecruitment/UserService/utils/app_error"
	httphelper "github.com/SawitProRecruitment/UserService/utils/http_helper"
	"github.com/SawitProRecruitment/UserService/utils/ptr"
	utilvalidator "github.com/SawitProRecruitment/UserService/utils/validator"
	"github.com/go-playground/validator/v10"
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
	if err := s.Validator.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, httphelper.ErrorResponse{
			Message: "Validation failed",
			Errors:  utilvalidator.FormatValidationErrors(validationErrors),
		})
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
	if err := s.Validator.Struct(params); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, httphelper.ErrorResponse{
			Message: "Validation failed",
			Errors:  utilvalidator.FormatValidationErrors(validationErrors),
		})
	}

	estate, err := s.Repository.GetEstateWithAllDetails(ctx, id)
	if err != nil {
		return httphelper.HttpRespError(c, err)
	}

	var (
		droneRest *Coordinate
		distance  int64
	)

	// currently if param max_distance is not provided,
	// we can directly get it from estate_stats table (pre calculate on create tree)
	// else must calculate it manually
	if params.MaxDistance == nil {
		distance = estate.Stats.DroneDistance
	} else {
		distance, droneRest = calculateDroneDistance(estate, params.MaxDistance)
	}

	var resp generated.DronePlanResponse
	resp.Distance = &distance
	if droneRest != nil {
		resp.Rest = &struct {
			X *int "json:\"x,omitempty\""
			Y *int "json:\"y,omitempty\""
		}{X: &droneRest.X, Y: &droneRest.Y}
	}

	return c.JSON(http.StatusOK, resp)
}

// Get estate statistics
// (GET /estate/{id}/stats)
func (s *Server) GetEstateIdStats(c echo.Context, id openapi_types.UUID) error {
	ctx := c.Request().Context()

	var resp generated.EstateStats
	result, err := s.Repository.GetEstateWithAllDetails(ctx, id, repository.RELATION_TREES)
	if err != nil {
		return httphelper.HttpRespError(c, err)
	}

	resp = generated.EstateStats{
		Count:  ptr.ToPointer[int64](result.Stats.TreeCount),
		Max:    ptr.ToPointer[int](result.Stats.MaxHeight),
		Median: ptr.ToPointer[int](result.Stats.MedianHeight),
		Min:    ptr.ToPointer[int](result.Stats.MinHeight),
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
	if err := s.Validator.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, httphelper.ErrorResponse{
			Message: "Validation failed",
			Errors:  utilvalidator.FormatValidationErrors(validationErrors),
		})
	}

	// #1. Insert the tree to DB
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

	// #2. Calculated the stats
	if err = s.calculateStats(ctx, id); err != nil {
		return httphelper.HttpRespError(c, err)
	}

	// #3 return response
	var resp generated.CreateEstateResponse
	openapiUUID := openapi_types.UUID(treeId)
	resp.Id = &openapiUUID

	return c.JSON(http.StatusCreated, resp)
}
