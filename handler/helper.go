package handler

import (
	"net/http"
	"strconv"

	"github.com/SawitProRecruitment/UserService/repository"
	httphelper "github.com/SawitProRecruitment/UserService/utils/http_helper"
	utilvalidator "github.com/SawitProRecruitment/UserService/utils/validator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (s *Server) validatePayload(c echo.Context, payload interface{}) (err error) {
	// Validate payload
	if err = s.Validator.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, httphelper.ErrorResponse{
			Message: "Validation failed",
			Errors:  utilvalidator.FormatValidationErrors(validationErrors),
		})
	}
	return
}

type mapTree map[string]repository.Tree

func newMapTree(trees []repository.Tree) mapTree {
	treeMap := make(map[string]repository.Tree)
	for _, tree := range trees {
		key := strconv.Itoa(tree.X) + "," + strconv.Itoa(tree.Y) // "x,y"
		treeMap[key] = tree
	}
	return treeMap
}

func (m mapTree) getTreeByCoordinate(x, y int) (tree repository.Tree, isExist bool) {
	tree, exists := m[getCoordinateKey(x, y)]
	return tree, exists
}

// get x,y format coordinate in string
func getCoordinateKey(x, y int) (res string) {
	return strconv.Itoa(x) + "," + strconv.Itoa(y)
}

// calculate drone distance until reach end destination estate plot
func CalculateDroneDistance(estate *repository.Estate) int {
	var (
		width           = estate.Width
		length          = estate.Length
		totalDistance   = 1
		lastDroneHeight = 1
	)

	// used map for easy access checking availibility tree in each plot
	// key => x,y and val => Tree
	mapTrees := newMapTree(estate.Trees)

	// loop coordinate every Y axis in estate (South to North)
	for y := 1; y <= width; y++ {

		// 1 => then go right
		// -1 => then go left
		direction := 1
		// determine for start, end & direction loop
		start, end := 1, length
		if y%2 == 1 {
			direction = 1
			start = 1
			end = length
		} else {
			direction = -1
			start = length
			end = 1
		}

		// loop coordinate every X axis in estate (West to East)
		// loop every column in estate base on direction
		for x := start; x != end+direction; x += direction {
			tree, isExist := mapTrees.getTreeByCoordinate(x, y)
			if !isExist {
				if x < length {
					totalDistance += 10
				}
				continue
			}

			// calculate total distance & last drone height
			heightForMonitorTree := tree.Height + 1 // for monitoring, need to add 1m
			if heightForMonitorTree > lastDroneHeight {
				distance := heightForMonitorTree - lastDroneHeight
				totalDistance += distance // need to up the drone
				lastDroneHeight += distance
			} else if heightForMonitorTree < lastDroneHeight {
				distance := lastDroneHeight - heightForMonitorTree
				totalDistance += distance
				lastDroneHeight -= distance
			}
			// move next row plot
			if x < length {
				totalDistance += 10
			}
		}

		// need to add 10m, because we need to go to up for next Y
		if y < width {
			totalDistance += 10
		}
	}

	// add distance for grounding drone
	totalDistance += lastDroneHeight

	return totalDistance
}
