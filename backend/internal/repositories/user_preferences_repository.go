package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"stops/backend/internal/models"
)

type UserPreferencesRepository struct {
	db *pgxpool.Pool
}

func NewUserPreferencesRepository(db *pgxpool.Pool) *UserPreferencesRepository {
	return &UserPreferencesRepository{db: db}
}

func (repo *UserPreferencesRepository) UpdatePreferences(ctx context.Context, userID string, homeLocation, workLocation, preferredTransport *string) error {
	const query = `
		UPDATE users
		SET home_location = COALESCE($2, home_location),
		    work_location = COALESCE($3, work_location),
		    preferred_transport_type = COALESCE($4, preferred_transport_type),
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := repo.db.Exec(ctx, query, userID, homeLocation, workLocation, preferredTransport)
	if err != nil {
		return fmt.Errorf("update user preferences: %w", err)
	}
	return nil
}

func (repo *UserPreferencesRepository) GetSavedPlaces(ctx context.Context, userID string) ([]models.SavedPlace, error) {
	const query = `
		SELECT id::text, user_id::text, name, ST_AsText(location), place_type, created_at, updated_at
		FROM saved_places
		WHERE user_id = $1
		ORDER BY place_type, created_at
	`
	rows, err := repo.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get saved places: %w", err)
	}
	defer rows.Close()

	var places []models.SavedPlace
	for rows.Next() {
		var p models.SavedPlace
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Location, &p.PlaceType, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan saved place: %w", err)
		}
		places = append(places, p)
	}
	return places, nil
}

func (repo *UserPreferencesRepository) CreateSavedPlace(ctx context.Context, place models.SavedPlace) error {
	const query = `
		INSERT INTO saved_places (user_id, name, location, place_type)
		VALUES ($1, $2, ST_GeomFromText($3, 4326), $4)
		RETURNING id, created_at, updated_at
	`
	err := repo.db.QueryRow(ctx, query, place.UserID, place.Name, place.Location, place.PlaceType).Scan(&place.ID, &place.CreatedAt, &place.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create saved place: %w", err)
	}
	return nil
}

func (repo *UserPreferencesRepository) DeleteSavedPlace(ctx context.Context, userID, placeID string) error {
	const query = `
		DELETE FROM saved_places
		WHERE id = $1 AND user_id = $2
	`
	result, err := repo.db.Exec(ctx, query, placeID, userID)
	if err != nil {
		return fmt.Errorf("delete saved place: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("saved place not found")
	}
	return nil
}

func (repo *UserPreferencesRepository) GetTripHistory(ctx context.Context, userID string, limit int) ([]models.Trip, error) {
	const query = `
		SELECT id::text, user_id::text, route_id::text, origin, destination, duration, fare, timestamp
		FROM trips
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`
	rows, err := repo.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get trip history: %w", err)
	}
	defer rows.Close()

	var trips []models.Trip
	for rows.Next() {
		var t models.Trip
		if err := rows.Scan(&t.ID, &t.UserID, &t.RouteID, &t.Origin, &t.Destination, &t.Duration, &t.Fare, &t.Timestamp); err != nil {
			return nil, fmt.Errorf("scan trip: %w", err)
		}
		trips = append(trips, t)
	}
	return trips, nil
}

func (repo *UserPreferencesRepository) RecordTrip(ctx context.Context, trip models.Trip) error {
	const query = `
		INSERT INTO trips (user_id, route_id, origin, destination, duration, fare, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	err := repo.db.QueryRow(ctx, query, trip.UserID, trip.RouteID, trip.Origin, trip.Destination, trip.Duration, trip.Fare, trip.Timestamp).Scan(&trip.ID, &trip.Timestamp)
	if err != nil {
		return fmt.Errorf("record trip: %w", err)
	}
	return nil
}
