package handler

import (
	"context"
	"strconv"

	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

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

type Plot struct {
	X          int
	Y          int
	treeHeight *int // will null if tree not exist
	nextPlot   *Plot
}

type Coordinate struct {
	X int
	Y int
}

// calculate drone distance until reach end destination estate plot
func calculateDroneDistance(estate *repository.Estate, maxDistance *int) (int64, *Coordinate) {
	var (
		totalDistance            int64 = 1 // initial value
		lastDroneHeight                = 1 // initial value
		lastDroneCoordinate      *Coordinate
		restDroneBatteryDistance = 0
	)

	// init last drone coordinate
	if maxDistance != nil {
		restDroneBatteryDistance = *maxDistance - 1 // because for the first ditance groun 1m
		lastDroneCoordinate = &Coordinate{X: 1, Y: 1}
	}

	current := createLinkedListByEstate(estate)

	// pre calculate for first head
	if current.treeHeight != nil {
		totalDistance += int64(*current.treeHeight)
		lastDroneHeight += *current.treeHeight
		if restDroneBatteryDistance > 0 {
			restDroneBatteryDistance -= *current.treeHeight
		}
	}

	for current != nil {
		// previous loop is last plot, stopping if next plot is nil
		if current.nextPlot == nil {
			current = current.nextPlot
			continue
		}

		var distanceUpDown int // distance for up or down drone
		// check for next plot tree exist, for go up or down the drone
		if current.nextPlot.treeHeight != nil {
			heightForMonitorTree := *current.nextPlot.treeHeight + 1
			if heightForMonitorTree > lastDroneHeight {
				distanceUpDown = heightForMonitorTree - lastDroneHeight
				totalDistance += int64(distanceUpDown)
				lastDroneHeight += distanceUpDown
			} else if heightForMonitorTree < lastDroneHeight {
				distanceUpDown = lastDroneHeight - heightForMonitorTree
				totalDistance += int64(distanceUpDown)
				lastDroneHeight -= distanceUpDown
			}
		}

		// travel to next plot
		totalDistance += 10

		// calculate rest drone battery if > 0
		if restDroneBatteryDistance > 0 {
			restDroneBatteryDistance -= distanceUpDown

			// check if better to rest at current plot or continue next plot (consider length drone to ground for avoid crash)
			nextPlotWithSafeLanding := 10 + lastDroneHeight
			if restDroneBatteryDistance < nextPlotWithSafeLanding {
				lastDroneCoordinate = &Coordinate{X: current.X, Y: current.Y}
				restDroneBatteryDistance = 0 // bacause we dont want for next iterate to update it
			} else {
				lastDroneCoordinate = &Coordinate{X: current.nextPlot.X, Y: current.nextPlot.Y}
			}

			restDroneBatteryDistance -= 10
		}

		current = current.nextPlot
	}

	// add distance for grounding drone
	totalDistance += int64(lastDroneHeight)

	return totalDistance, lastDroneCoordinate
}

// create linked list to easy calculate the distance
func createLinkedListByEstate(estate *repository.Estate) (head *Plot) {
	var (
		current *Plot
		width   = estate.Width
		length  = estate.Length
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
			// creating plot for every x,y coordinate
			plot := &Plot{X: x, Y: y}

			tree, isExist := mapTrees.getTreeByCoordinate(x, y)
			if isExist {
				plot.treeHeight = &tree.Height
			}

			if head == nil {
				head = plot
				current = plot
			} else {
				current.nextPlot = plot
				current = plot
			}
		}
	}

	return
}

// calculate & save stats into estate stats
func (s *Server) calculateStats(ctx context.Context, estateId uuid.UUID) (err error) {
	var calculatedStats *repository.EstateStats
	var estate *repository.Estate
	g, groupCtx := errgroup.WithContext(ctx)
	g.Go(func() (errCalc error) {
		calculatedStats, errCalc = s.Repository.GetCalculatedEstateStats(groupCtx, estateId)
		return
	})

	g.Go(func() (errGetState error) {
		estate, errGetState = s.Repository.GetEstateWithAllDetails(groupCtx, estateId)
		return
	})

	// Wait for both goroutines to complete
	if err := g.Wait(); err != nil {
		return err
	}

	// pre calculate the distance after inserting the tree
	distance, _ := calculateDroneDistance(estate, nil)

	// set stats for saving to DB
	if estate.Stats.Id == uuid.Nil {
		estate.Stats.Id = uuid.New()
	}
	estate.Stats.TreeCount = calculatedStats.TreeCount
	estate.Stats.DroneDistance = int64(distance)
	estate.Stats.MaxHeight = calculatedStats.MaxHeight
	estate.Stats.MedianHeight = calculatedStats.MedianHeight
	estate.Stats.MinHeight = calculatedStats.MinHeight

	if err = s.Repository.UpsertEstateStats(ctx, estateId, estate.Stats); err != nil {
		return
	}

	return
}
