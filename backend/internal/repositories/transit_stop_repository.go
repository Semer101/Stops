package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"stops/backend/internal/models"
)

type TransitStopRepository struct {
	db *pgxpool.Pool
}

type NearbyStopsParams struct {
	Latitude     float64
	Longitude    float64
	RadiusMeters int
	Limit        int
}

func NewTransitStopRepository(db *pgxpool.Pool) *TransitStopRepository {
	return &TransitStopRepository{db: db}
}

func (repository *TransitStopRepository) FindNearby(ctx context.Context, params NearbyStopsParams) ([]models.TransitStop, error) {
	const query = `
		SELECT
			id::text,
			name,
			type,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude,
			area,
			verified,
			votes,
			availability_status,
			ST_Distance(
				location::geography,
				ST_SetSRID(ST_MakePoint($2, $1), 4326)::geography
			) AS distance_meters
		FROM transit_stops
		WHERE verified = TRUE
		  AND ST_DWithin(
			location::geography,
			ST_SetSRID(ST_MakePoint($2, $1), 4326)::geography,
			$3
		  )
		ORDER BY distance_meters ASC
		LIMIT $4
	`

	rows, err := repository.db.Query(ctx, query, params.Latitude, params.Longitude, params.RadiusMeters, params.Limit)
	if err != nil {
		return nil, fmt.Errorf("query nearby transit stops: %w", err)
	}
	defer rows.Close()

	stops := make([]models.TransitStop, 0)
	for rows.Next() {
		var stop models.TransitStop
		var stopType string
		if err := rows.Scan(
			&stop.ID,
			&stop.Name,
			&stopType,
			&stop.Latitude,
			&stop.Longitude,
			&stop.Area,
			&stop.Verified,
			&stop.Votes,
			&stop.AvailabilityStatus,
			&stop.DistanceMeters,
		); err != nil {
			return nil, fmt.Errorf("scan nearby transit stop: %w", err)
		}

		stop.Type = models.TransitStopType(stopType)
		stops = append(stops, stop)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate nearby transit stops: %w", err)
	}

	return stops, nil
}
