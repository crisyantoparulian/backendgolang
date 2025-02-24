package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) createTreeSQL(ctx context.Context, input CreateTreeInput) (err error) {
	res, err := r.Db.ExecContext(ctx, `
		INSERT INTO trees (id, estate_id, x, y, height) 
		VALUES ($1, $2, $3, $4, $5);`,
		input.Id, input.EstateId, input.X, input.Y, input.Height)
	if err != nil {
		return err
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected < 1 {
		return errors.New("no rows affected")
	}

	return
}

func (r *Repository) getEstateByIdSql(ctx context.Context, id uuid.UUID) (estate *Estate, err error) {
	estate = &Estate{}
	query := `
		SELECT id, width, length, created_at, updated_at
		FROM estates
		WHERE id = $1;`
	err = r.Db.QueryRowContext(ctx, query, id).Scan(
		&estate.Id,
		&estate.Width,
		&estate.Length,
		&estate.CreatedAt,
		&estate.UpdatedAt,
	)
	return
}

func (r *Repository) checkExistEstateTree(ctx context.Context, input CheckExistEstateTreeInput) (isExist bool, err error) {
	err = r.Db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM trees 
			WHERE estate_id = $1 AND x = $2 AND y = $3
		);`,
		input.EstateId, input.X, input.Y).Scan(&isExist)
	if err != nil {
		return
	}
	return
}

func (r *Repository) createEstateSql(ctx context.Context, input CreateEstateInput) (err error) {
	var res sql.Result
	res, err = r.Db.ExecContext(ctx, `
		INSERT INTO estates (id, width, length)  
			VALUES ($1, $2, $3);`,
		input.Id, input.Width, input.Length)
	if err != nil {
		return
	}

	var rowAffected int64
	rowAffected, err = res.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected < 1 {
		return errors.New("no rows affected")
	}

	return
}

func (r *Repository) getAllTreeHeightEstateSql(ctx context.Context, estateID uuid.UUID) (heights []int, err error) {

	query := `
		SELECT height
		FROM trees
		WHERE estate_id = $1
		ORDER BY height;`
	rows, err := r.Db.QueryContext(ctx, query, estateID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var height int
		if err = rows.Scan(&height); err != nil {
			return
		}
		heights = append(heights, height)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (r *Repository) getTreesByEstateId(ctx context.Context, estateID uuid.UUID) ([]Tree, error) {
	// Query untuk mendapatkan semua pohon di estate
	queryTrees := `
		SELECT id, estate_id, x, y, height, created_at, updated_at
		FROM trees
		WHERE estate_id = $1;`
	rows, err := r.Db.QueryContext(ctx, queryTrees, estateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Simpan semua pohon dalam slice
	var trees []Tree
	for rows.Next() {
		var tree Tree
		if err := rows.Scan(
			&tree.Id,
			&tree.EstateId,
			&tree.X,
			&tree.Y,
			&tree.Height,
			&tree.CreatedAt,
			&tree.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan tree data: %w", err)
		}
		trees = append(trees, tree)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return trees, nil
}

func (r *Repository) getCalculatedEstateStatsSQL(ctx context.Context, estateId uuid.UUID) (stats *EstateStats, err error) {
	stats = &EstateStats{}
	query := `
		SELECT
			COUNT(*) AS tree_count,
			COALESCE(MAX(height), 0) AS max_height,
			COALESCE(MIN(height), 0) AS min_height,
			COALESCE(ROUND(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY height)), 0) AS median_height
		FROM trees
		WHERE estate_id = $1;`
	err = r.Db.QueryRowContext(ctx, query, estateId).Scan(
		&stats.TreeCount,
		&stats.MaxHeight,
		&stats.MinHeight,
		&stats.MedianHeight,
	)
	return
}

func (r *Repository) getEstateStatsSQL(ctx context.Context, estateId uuid.UUID) (stats *EstateStats, err error) {
	stats = &EstateStats{}
	query := `
		SELECT id, estate_id, tree_count, max_height, min_height, median_height,drone_distance, created_at, updated_at
		FROM estate_stats
		WHERE estate_id = $1;`
	err = r.Db.QueryRowContext(ctx, query, estateId).Scan(
		&stats.Id,
		&stats.EstateID,
		&stats.TreeCount,
		&stats.MaxHeight,
		&stats.MinHeight,
		&stats.MedianHeight,
		&stats.DroneDistance,
		&stats.CreatedAt,
		&stats.UpdatedAt,
	)
	if err != nil {
		return
	}

	return
}

func (r *Repository) upsertEstateStatsSQL(ctx context.Context, estateID uuid.UUID, stats *EstateStats) error {
	// Query untuk menyimpan atau memperbarui statistik
	query := `
		INSERT INTO estate_stats (id, estate_id, tree_count, max_height, min_height, median_height, drone_distance)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (estate_id) DO UPDATE
		SET
			tree_count = EXCLUDED.tree_count,
			max_height = EXCLUDED.max_height,
			min_height = EXCLUDED.min_height,
			median_height = EXCLUDED.median_height,
			drone_distance = EXCLUDED.drone_distance,
			updated_at = CURRENT_TIMESTAMP;`
	_, err := r.Db.ExecContext(ctx, query,
		stats.Id,
		estateID,
		stats.TreeCount,
		stats.MaxHeight,
		stats.MinHeight,
		stats.MedianHeight,
		stats.DroneDistance,
	)
	if err != nil {
		return err
	}

	return nil
}
