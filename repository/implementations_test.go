package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	apperror "github.com/SawitProRecruitment/UserService/utils/app_error"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateEstate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repository{Db: db}

	testInput := CreateEstateInput{
		Id:     uuid.New(),
		Width:  1,
		Length: 5,
	}

	tests := []struct {
		name          string
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Success",
			mockSetup: func() {
				mock.ExpectExec(`INSERT INTO estate`).
					WithArgs(testInput.Id, testInput.Width, testInput.Length).
					WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected
			},
			expectedError: nil,
		},
		{
			name: "Database Error",
			mockSetup: func() {
				mock.ExpectExec(`INSERT INTO estate`).
					WithArgs(testInput.Id, testInput.Width, testInput.Length).
					WillReturnError(errors.New("database error"))
			},
			expectedError: apperror.WrapWithCode(fmt.Errorf("failed to create estate: %w", errors.New("database error")), http.StatusInternalServerError),
		},
		{
			name: "No Rows Affected",
			mockSetup: func() {
				mock.ExpectExec(`INSERT INTO estate`).
					WithArgs(testInput.Id, testInput.Width, testInput.Length).
					WillReturnResult(sqlmock.NewResult(1, 0)) // 0 rows affected
			},
			expectedError: apperror.WrapWithCode(fmt.Errorf("failed to create estate: %w", errors.New("no rows affected")), http.StatusInternalServerError),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			err := repo.CreateEstate(context.Background(), testInput)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet()) // Ensure all expectations were met
		})
	}
}

func TestGetEstateWithAllDetails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repository{Db: db}
	ctx := context.Background()
	estateID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()

	// Mock estate data
	estateRow := mock.NewRows([]string{"id", "width", "length", "created_at", "updated_at"}).
		AddRow(estateID, 100, 200, createdAt, updatedAt)

	// Mock tree data
	treeRow := mock.NewRows([]string{"id", "estate_id", "x", "y", "height", "created_at", "updated_at"}).
		AddRow(uuid.New(), estateID, 10, 20, 5, createdAt, updatedAt)

	// // Mock estate stats data
	statsRow := mock.NewRows([]string{"id", "estate_id", "tree_count", "max_height", "min_height", "median_height", "drone_distance", "created_at", "updated_at"}).
		AddRow(uuid.New(), estateID, 5, 10, 2, 6, 100, createdAt, updatedAt)

	tests := []struct {
		name             string
		mockSetup        func()
		excludeRelations []Relation
		expectedError    error
	}{
		{
			name: "Success - All Details",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length, created_at, updated_at FROM estates WHERE id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(estateRow)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, estate_id, x, y, height, created_at, updated_at FROM trees WHERE estate_id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(treeRow)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, estate_id, tree_count, max_height, min_height, median_height,drone_distance, created_at, updated_at
		    	            			FROM estate_stats
		    	            			WHERE estate_id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(statsRow)
			},
			excludeRelations: []Relation{},
			expectedError:    nil,
		},
		{
			name: "Estate Not Found",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length, created_at, updated_at FROM estates WHERE id = $1;`)).
					WithArgs(estateID).
					WillReturnError(sql.ErrNoRows)
			},
			excludeRelations: []Relation{},
			expectedError:    apperror.WrapWithCode(fmt.Errorf("estate with ID %s not found", estateID), http.StatusNotFound),
		},
		{
			name: "Error Fetching Estate",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length, created_at, updated_at FROM estates WHERE id = $1;`)).
					WithArgs(estateID).
					WillReturnError(errors.New("db error"))
			},
			excludeRelations: []Relation{},
			expectedError:    apperror.WrapWithCode(fmt.Errorf("failed to get estate: %w", errors.New("db error")), http.StatusInternalServerError),
		},
		{
			name: "Error fetching trees",
			mockSetup: func() {
				estateRow := mock.NewRows([]string{"id", "width", "length", "created_at", "updated_at"}).
					AddRow(estateID, 100, 200, createdAt, updatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length, created_at, updated_at FROM estates WHERE id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(estateRow)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, estate_id, x, y, height, created_at, updated_at FROM trees WHERE estate_id = $1;`)).
					WithArgs(estateID).
					WillReturnError(errors.New("db error"))

			},
			excludeRelations: []Relation{},
			expectedError:    apperror.WrapWithCode(fmt.Errorf("failed to get estate tress: %w", errors.New("db error")), http.StatusInternalServerError),
		},
		{
			name: "Error Fetching Stats",
			mockSetup: func() {
				estateRow := mock.NewRows([]string{"id", "width", "length", "created_at", "updated_at"}).
					AddRow(estateID, 100, 200, createdAt, updatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length, created_at, updated_at FROM estates WHERE id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(estateRow)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, estate_id, x, y, height, created_at, updated_at FROM trees WHERE estate_id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(treeRow)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, estate_id, tree_count, max_height, min_height, median_height,drone_distance, created_at, updated_at
		    	            			FROM estate_stats
		    	            			WHERE estate_id = $1;`)).
					WithArgs(estateID).
					WillReturnError(errors.New("db error"))
			},
			excludeRelations: []Relation{},
			expectedError:    apperror.WrapWithCode(fmt.Errorf("failed to get estate stats: %w", errors.New("db error")), http.StatusInternalServerError),
		},
		// 	name: "Exclude Trees and Stats",
		// 	mockSetup: func() {
		// 		mock.ExpectQuery(`SELECT id, width, length, created_at, updated_at FROM estates WHERE id = \$1`).
		// 			WithArgs(estateID).
		// 			WillReturnRows(estateRow)
		// 	},
		// 	excludeRelations: []Relation{RELATION_TREES, RELATION_STATS},
		// 	expectedError:    nil,
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			estate, err := repo.GetEstateWithAllDetails(ctx, estateID, tc.excludeRelations...)

			if tc.expectedError != nil {
				assert.Nil(t, estate)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NotNil(t, estate)
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet()) // Ensure all expectations were met
		})
	}
}

func TestCreateTree(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repository{Db: db}
	ctx := context.Background()
	estateID := uuid.New()
	treeID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()

	// Mock estate data
	estateRow := mock.NewRows([]string{"id", "width", "length", "created_at", "updated_at"}).
		AddRow(estateID, 100, 200, createdAt, updatedAt)

	tests := []struct {
		name          string
		mockSetup     func()
		input         CreateTreeInput
		expectedError error
	}{
		{
			name: "Success - Create Tree",
			mockSetup: func() {
				// Mock getEstateByIdSql
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length, created_at, updated_at FROM estates WHERE id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(estateRow)

				// Mock checkExistEstateTree
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS ( SELECT 1 FROM trees WHERE estate_id = $1 AND x = $2 AND y = $3 );`)).
					WithArgs(estateID, 10, 20).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

				// Mock createTreeSQL
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO trees (id, estate_id, x, y, height) VALUES ($1, $2, $3, $4, $5);`)).
					WithArgs(treeID, estateID, 10, 20, 15).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: CreateTreeInput{
				Id:       treeID,
				EstateId: estateID,
				X:        10,
				Y:        20,
				Height:   15,
			},
			expectedError: nil,
		},
		{
			name: "Estate Not Found",
			mockSetup: func() {
				// Mock getEstateByIdSql
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length, created_at, updated_at FROM estates WHERE id = $1;`)).
					WithArgs(estateID).
					WillReturnError(sql.ErrNoRows)
			},
			input: CreateTreeInput{
				Id:       treeID,
				EstateId: estateID,
				X:        10,
				Y:        20,
				Height:   15,
			},
			expectedError: apperror.WrapWithCode(fmt.Errorf("estate with ID %s not found", estateID), http.StatusNotFound),
		},
		{
			name: "Invalid X Coordinate",
			mockSetup: func() {
				estateRow := mock.NewRows([]string{"id", "width", "length", "created_at", "updated_at"}).
					AddRow(estateID, 100, 200, createdAt, updatedAt)

				// Mock getEstateByIdSql
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length, created_at, updated_at FROM estates WHERE id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(estateRow)
			},
			input: CreateTreeInput{
				Id:       treeID,
				EstateId: estateID,
				X:        300, // X coordinate greater than estate length
				Y:        20,
				Height:   15,
			},
			expectedError: apperror.WrapWithCode(fmt.Errorf("coordinate x cannot greater than %d", 200), http.StatusBadRequest),
		},
		{
			name: "Invalid Y Coordinate",
			mockSetup: func() {
				estateRow := mock.NewRows([]string{"id", "width", "length", "created_at", "updated_at"}).
					AddRow(estateID, 100, 200, createdAt, updatedAt)
				// Mock getEstateByIdSql
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length, created_at, updated_at FROM estates WHERE id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(estateRow)
			},
			input: CreateTreeInput{
				Id:       treeID,
				EstateId: estateID,
				X:        10,
				Y:        300, // Y coordinate greater than estate width
				Height:   15,
			},
			expectedError: apperror.WrapWithCode(fmt.Errorf("coordinate y cannot greater than %d", 100), http.StatusBadRequest),
		},
		{
			name: "Tree Already Exists",
			mockSetup: func() {
				estateRow := mock.NewRows([]string{"id", "width", "length", "created_at", "updated_at"}).
					AddRow(estateID, 100, 200, createdAt, updatedAt)
				// Mock getEstateByIdSql
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length, created_at, updated_at FROM estates WHERE id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(estateRow)

				// Mock checkExistEstateTree
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS ( SELECT 1 FROM trees WHERE estate_id = $1 AND x = $2 AND y = $3 );`)).
					WithArgs(estateID, 10, 20).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			input: CreateTreeInput{
				Id:       treeID,
				EstateId: estateID,
				X:        10,
				Y:        20,
				Height:   15,
			},
			expectedError: apperror.WrapWithCode(errors.New("plot already has a tree"), http.StatusUnprocessableEntity),
		},
		{
			name: "Database Error - Check Tree Existence",
			mockSetup: func() {
				estateRow := mock.NewRows([]string{"id", "width", "length", "created_at", "updated_at"}).
					AddRow(estateID, 100, 200, createdAt, updatedAt)

				// Mock getEstateByIdSql
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length, created_at, updated_at FROM estates WHERE id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(estateRow)

				// Mock checkExistEstateTree
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS ( SELECT 1 FROM trees WHERE estate_id = $1 AND x = $2 AND y = $3 );`)).
					WithArgs(estateID, 10, 20).
					WillReturnError(errors.New("db error"))
			},
			input: CreateTreeInput{
				Id:       treeID,
				EstateId: estateID,
				X:        10,
				Y:        20,
				Height:   15,
			},
			expectedError: apperror.WrapWithCode(fmt.Errorf("failed to check plot tree in database: %w", errors.New("db error")), http.StatusInternalServerError),
		},
		{
			name: "Database Error - Create Tree",
			mockSetup: func() {
				estateRow := mock.NewRows([]string{"id", "width", "length", "created_at", "updated_at"}).
					AddRow(estateID, 100, 200, createdAt, updatedAt)

				// Mock getEstateByIdSql
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, width, length, created_at, updated_at FROM estates WHERE id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(estateRow)

				// Mock checkExistEstateTree
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS ( SELECT 1 FROM trees WHERE estate_id = $1 AND x = $2 AND y = $3 );`)).
					WithArgs(estateID, 10, 20).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

				// Mock createTreeSQL
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO trees (id, estate_id, x, y, height) VALUES ($1, $2, $3, $4, $5);`)).
					WithArgs(treeID, estateID, 10, 20, 15).
					WillReturnError(errors.New("db error"))
			},
			input: CreateTreeInput{
				Id:       treeID,
				EstateId: estateID,
				X:        10,
				Y:        20,
				Height:   15,
			},
			expectedError: apperror.WrapWithCode(fmt.Errorf("failed to create tree: %w", errors.New("db error")), http.StatusInternalServerError),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			err := repo.CreateTree(ctx, tc.input)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet()) // Ensure all expectations were met
		})
	}
}

func TestGetCalculatedEstateStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repository{Db: db}
	ctx := context.Background()
	estateID := uuid.New()
	emptyStats := &EstateStats{}

	tests := []struct {
		name          string
		mockSetup     func()
		estateId      uuid.UUID
		expectedStats *EstateStats
		expectedError error
	}{
		{
			name: "Success - Get Calculated Stats",
			mockSetup: func() {
				// Mock getCalculatedEstateStatsSQL
				rows := sqlmock.NewRows([]string{"tree_count", "max_height", "min_height", "median_height"}).
					AddRow(10, 20, 5, 12)
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT
						COUNT(*) AS tree_count,
						COALESCE(MAX(height), 0) AS max_height,
						COALESCE(MIN(height), 0) AS min_height,
						COALESCE(ROUND(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY height)), 0) AS median_height
					FROM trees
					WHERE estate_id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(rows)
			},
			estateId: estateID,
			expectedStats: &EstateStats{
				TreeCount:    10,
				MaxHeight:    20,
				MinHeight:    5,
				MedianHeight: 12,
			},
			expectedError: nil,
		},
		{
			name: "Database Error",
			mockSetup: func() {
				// Mock getCalculatedEstateStatsSQL
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT
						COUNT(*) AS tree_count,
						COALESCE(MAX(height), 0) AS max_height,
						COALESCE(MIN(height), 0) AS min_height,
						COALESCE(ROUND(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY height)), 0) AS median_height
					FROM trees
					WHERE estate_id = $1;`)).
					WithArgs(estateID).
					WillReturnError(errors.New("db error"))
			},
			estateId:      estateID,
			expectedStats: emptyStats,
			expectedError: apperror.WrapWithCode(fmt.Errorf("failed to get calculated stats: %w", errors.New("db error")), http.StatusInternalServerError),
		},
		{
			name: "No Trees in Estate",
			mockSetup: func() {
				// Mock getCalculatedEstateStatsSQL
				rows := sqlmock.NewRows([]string{"tree_count", "max_height", "min_height", "median_height"}).
					AddRow(0, 0, 0, 0)
				mock.ExpectQuery(regexp.QuoteMeta(`
					SELECT
						COUNT(*) AS tree_count,
						COALESCE(MAX(height), 0) AS max_height,
						COALESCE(MIN(height), 0) AS min_height,
						COALESCE(ROUND(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY height)), 0) AS median_height
					FROM trees
					WHERE estate_id = $1;`)).
					WithArgs(estateID).
					WillReturnRows(rows)
			},
			estateId: estateID,
			expectedStats: &EstateStats{
				TreeCount:    0,
				MaxHeight:    0,
				MinHeight:    0,
				MedianHeight: 0,
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			stats, err := repo.GetCalculatedEstateStats(ctx, tc.estateId)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Equal(t, tc.expectedStats, stats)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedStats, stats)
			}

			assert.NoError(t, mock.ExpectationsWereMet()) // Ensure all expectations were met
		})
	}
}

func TestUpsertEstateStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repository{Db: db}
	ctx := context.Background()
	estateID := uuid.New()

	stats := &EstateStats{
		Id:            uuid.New(),
		TreeCount:     10,
		MaxHeight:     20,
		MinHeight:     5,
		MedianHeight:  12,
		DroneDistance: 100,
	}

	tests := []struct {
		name          string
		mockSetup     func()
		estateId      uuid.UUID
		stats         *EstateStats
		expectedError error
	}{
		{
			name: "Success - Upsert Estate Stats",
			mockSetup: func() {
				// Mock upsertEstateStatsSQL
				mock.ExpectExec(regexp.QuoteMeta(`
					INSERT INTO estate_stats (id, estate_id, tree_count, max_height, min_height, median_height, drone_distance)
					VALUES ($1, $2, $3, $4, $5, $6, $7)
					ON CONFLICT (estate_id) DO UPDATE
					SET
						tree_count = EXCLUDED.tree_count,
						max_height = EXCLUDED.max_height,
						min_height = EXCLUDED.min_height,
						median_height = EXCLUDED.median_height,
						drone_distance = EXCLUDED.drone_distance,
						updated_at = CURRENT_TIMESTAMP;`)).
					WithArgs(
						stats.Id,
						estateID,
						stats.TreeCount,
						stats.MaxHeight,
						stats.MinHeight,
						stats.MedianHeight,
						stats.DroneDistance,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			estateId:      estateID,
			stats:         stats,
			expectedError: nil,
		},
		{
			name: "Database Error",
			mockSetup: func() {
				// Mock upsertEstateStatsSQL
				mock.ExpectExec(regexp.QuoteMeta(`
					INSERT INTO estate_stats (id, estate_id, tree_count, max_height, min_height, median_height, drone_distance)
					VALUES ($1, $2, $3, $4, $5, $6, $7)
					ON CONFLICT (estate_id) DO UPDATE
					SET
						tree_count = EXCLUDED.tree_count,
						max_height = EXCLUDED.max_height,
						min_height = EXCLUDED.min_height,
						median_height = EXCLUDED.median_height,
						drone_distance = EXCLUDED.drone_distance,
						updated_at = CURRENT_TIMESTAMP;`)).
					WithArgs(
						stats.Id,
						estateID,
						stats.TreeCount,
						stats.MaxHeight,
						stats.MinHeight,
						stats.MedianHeight,
						stats.DroneDistance,
					).
					WillReturnError(errors.New("db error"))
			},
			estateId:      estateID,
			stats:         stats,
			expectedError: apperror.WrapWithCode(fmt.Errorf("failed to get upsert stats: %w", errors.New("db error")), http.StatusInternalServerError),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			err := repo.UpsertEstateStats(ctx, tc.estateId, tc.stats)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet()) // Ensure all expectations were met
		})
	}
}
