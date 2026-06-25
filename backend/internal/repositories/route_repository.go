package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"stops/backend/internal/models"
)

type RouteRepository struct {
	db *pgxpool.Pool
}

func NewRouteRepository(db *pgxpool.Pool) *RouteRepository {
	return &RouteRepository{db: db}
}

func (repository *RouteRepository) FindAll(ctx context.Context) ([]models.Route, error) {
	const query = `
		SELECT id::text, route_name, start_point, destination_point, COALESCE(fare, 0), city, COALESCE(estimated_time, 0), status, created_at, updated_at
		FROM routes
		WHERE status = 'Active'
	`
	rows, err := repository.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all routes: %w", err)
	}
	defer rows.Close()

	var routes []models.Route
	for rows.Next() {
		var r models.Route
		if err := rows.Scan(&r.ID, &r.RouteName, &r.StartPoint, &r.DestinationPoint, &r.Fare, &r.City, &r.EstimatedTime, &r.Status, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan route: %w", err)
		}
		routes = append(routes, r)
	}
	return routes, nil
}

func (repository *RouteRepository) FindByStartAndDestination(ctx context.Context, start, dest string) ([]models.Route, error) {
	const query = `
		SELECT id::text, route_name, start_point, destination_point, COALESCE(fare, 0), city, COALESCE(estimated_time, 0), status, created_at, updated_at
		FROM routes
		WHERE status = 'Active'
		  AND (LOWER(start_point) LIKE LOWER($1) OR LOWER(destination_point) LIKE LOWER($2))
	`
	rows, err := repository.db.Query(ctx, query, "%"+start+"%", "%"+dest+"%")
	if err != nil {
		return nil, fmt.Errorf("query routes by points: %w", err)
	}
	defer rows.Close()

	var routes []models.Route
	for rows.Next() {
		var r models.Route
		if err := rows.Scan(&r.ID, &r.RouteName, &r.StartPoint, &r.DestinationPoint, &r.Fare, &r.City, &r.EstimatedTime, &r.Status, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan route by points: %w", err)
		}
		routes = append(routes, r)
	}
	return routes, nil
}

func (repository *RouteRepository) GetFaresByRouteID(ctx context.Context, routeID string) ([]models.Fare, error) {
	const query = `
		SELECT id::text, route_id::text, transport_type, amount, created_at, updated_at
		FROM fares
		WHERE route_id = $1
	`
	rows, err := repository.db.Query(ctx, query, routeID)
	if err != nil {
		return nil, fmt.Errorf("query fares by route ID: %w", err)
	}
	defer rows.Close()

	var fares []models.Fare
	for rows.Next() {
		var f models.Fare
		if err := rows.Scan(&f.ID, &f.RouteID, &f.TransportType, &f.Amount, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan fare: %w", err)
		}
		fares = append(fares, f)
	}
	return fares, nil
}
