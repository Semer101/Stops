package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	db *pgxpool.Pool
}

type PendingReport struct {
	ReportID    string    `json:"reportId"`
	StopID      string    `json:"stopId"`
	StopName    string    `json:"stopName"`
	StopType    string    `json:"stopType"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Area        *string   `json:"area"`
	NetVotes    int       `json:"netVotes"`
	CreatedBy   string    `json:"createdBy"`
	CreatorName string    `json:"creatorName"`
	CreatedAt   time.Time `json:"createdAt"`
}

type AuditLog struct {
	ID          string    `json:"id"`
	AdminUserID string    `json:"adminUserId"`
	AdminName   string    `json:"adminName"`
	Action      string    `json:"action"`
	EntityType  string    `json:"entityType"`
	EntityID    *string   `json:"entityId"`
	Metadata    string    `json:"metadata"` // JSON string representation
	CreatedAt   time.Time `json:"createdAt"`
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{db: db}
}

func (repository *AdminRepository) FindPendingReports(ctx context.Context) ([]PendingReport, error) {
	const query = `
		SELECT
			r.id::text,
			s.id::text,
			s.name,
			s.type,
			ST_Y(s.location::geometry) AS latitude,
			ST_X(s.location::geometry) AS longitude,
			s.area,
			r.net_votes,
			r.created_by::text,
			u.name,
			r.created_at
		FROM crowdsource_reports r
		JOIN transit_stops s ON r.stop_id = s.id
		JOIN users u ON r.created_by = u.id
		WHERE r.status = 'pending'
		ORDER BY r.created_at DESC
	`
	rows, err := repository.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query pending reports: %w", err)
	}
	defer rows.Close()

	var reports []PendingReport
	for rows.Next() {
		var pr PendingReport
		if err := rows.Scan(
			&pr.ReportID,
			&pr.StopID,
			&pr.StopName,
			&pr.StopType,
			&pr.Latitude,
			&pr.Longitude,
			&pr.Area,
			&pr.NetVotes,
			&pr.CreatedBy,
			&pr.CreatorName,
			&pr.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan pending report: %w", err)
		}
		reports = append(reports, pr)
	}
	return reports, nil
}

func (repository *AdminRepository) VerifyReport(ctx context.Context, reportID string) (string, error) {
	tx, err := repository.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin verify transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const updateReportQuery = `
		UPDATE crowdsource_reports
		SET status = 'verified', updated_at = NOW()
		WHERE id = $1
		RETURNING stop_id::text
	`
	var stopID string
	err = tx.QueryRow(ctx, updateReportQuery, reportID).Scan(&stopID)
	if err != nil {
		return "", fmt.Errorf("update report status to verified: %w", err)
	}

	const updateStopQuery = `
		UPDATE transit_stops
		SET verified = TRUE, updated_at = NOW()
		WHERE id = $1
	`
	_, err = tx.Exec(ctx, updateStopQuery, stopID)
	if err != nil {
		return "", fmt.Errorf("update stop status to verified: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit verify transaction: %w", err)
	}

	return stopID, nil
}

func (repository *AdminRepository) RejectReport(ctx context.Context, reportID string) (string, error) {
	tx, err := repository.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin reject transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const updateReportQuery = `
		UPDATE crowdsource_reports
		SET status = 'rejected', updated_at = NOW()
		WHERE id = $1
		RETURNING stop_id::text
	`
	var stopID string
	err = tx.QueryRow(ctx, updateReportQuery, reportID).Scan(&stopID)
	if err != nil {
		return "", fmt.Errorf("update report status to rejected: %w", err)
	}

	const deleteStopQuery = `
		DELETE FROM transit_stops
		WHERE id = $1
	`
	_, err = tx.Exec(ctx, deleteStopQuery, stopID)
	if err != nil {
		return "", fmt.Errorf("delete rejected transit stop: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit reject transaction: %w", err)
	}

	return stopID, nil
}

func (repository *AdminRepository) CreateAuditLog(ctx context.Context, adminUserID, action, entityType string, entityID *string, metadata map[string]interface{}) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("marshal audit metadata: %w", err)
	}

	const query = `
		INSERT INTO admin_audit_logs (admin_user_id, action, entity_type, entity_id, metadata)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = repository.db.Exec(ctx, query, adminUserID, action, entityType, entityID, metadataBytes)
	if err != nil {
		return fmt.Errorf("insert admin audit log: %w", err)
	}
	return nil
}

func (repository *AdminRepository) FindAuditLogs(ctx context.Context) ([]AuditLog, error) {
	const query = `
		SELECT
			l.id::text,
			l.admin_user_id::text,
			u.name,
			l.action,
			l.entity_type,
			l.entity_id::text,
			l.metadata::text,
			l.created_at
		FROM admin_audit_logs l
		JOIN users u ON l.admin_user_id = u.id
		ORDER BY l.created_at DESC
	`
	rows, err := repository.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var al AuditLog
		if err := rows.Scan(
			&al.ID,
			&al.AdminUserID,
			&al.AdminName,
			&al.Action,
			&al.EntityType,
			&al.EntityID,
			&al.Metadata,
			&al.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan audit log: %w", err)
		}
		logs = append(logs, al)
	}
	return logs, nil
}

func (repository *AdminRepository) CreateStop(ctx context.Context, name, stopType string, lat, lng float64, area string, verified bool) (string, error) {
	const query = `
		INSERT INTO transit_stops (name, type, location, area, verified)
		VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326), $5, $6)
		RETURNING id::text
	`
	var stopID string
	err := repository.db.QueryRow(ctx, query, name, stopType, lng, lat, area, verified).Scan(&stopID)
	if err != nil {
		return "", fmt.Errorf("create stop: %w", err)
	}
	return stopID, nil
}

func (repository *AdminRepository) UpdateStop(ctx context.Context, stopID string, name, stopType *string, lat, lng *float64, area *string, verified *bool) error {
	const query = `
		UPDATE transit_stops
		SET name = COALESCE($2, name),
		    type = COALESCE($3, type),
		    location = COALESCE(ST_SetSRID(ST_MakePoint($4, $5), 4326), location),
		    area = COALESCE($6, area),
		    verified = COALESCE($7, verified),
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := repository.db.Exec(ctx, query, stopID, name, stopType, lng, lat, area, verified)
	if err != nil {
		return fmt.Errorf("update stop: %w", err)
	}
	return nil
}

func (repository *AdminRepository) DeleteStop(ctx context.Context, stopID string) error {
	const query = `DELETE FROM transit_stops WHERE id = $1`
	result, err := repository.db.Exec(ctx, query, stopID)
	if err != nil {
		return fmt.Errorf("delete stop: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("stop not found")
	}
	return nil
}

func (repository *AdminRepository) CreateRoute(ctx context.Context, routeName, startPoint, destPoint string, fare float64, city string, estimatedTime int) (string, error) {
	const query = `
		INSERT INTO routes (route_name, start_point, destination_point, fare, city, estimated_time)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text
	`
	var routeID string
	err := repository.db.QueryRow(ctx, query, routeName, startPoint, destPoint, fare, city, estimatedTime).Scan(&routeID)
	if err != nil {
		return "", fmt.Errorf("create route: %w", err)
	}
	return routeID, nil
}

func (repository *AdminRepository) UpdateRoute(ctx context.Context, routeID string, routeName, startPoint, destPoint *string, fare *float64, city *string, estimatedTime *int, status *string) error {
	const query = `
		UPDATE routes
		SET route_name = COALESCE($2, route_name),
		    start_point = COALESCE($3, start_point),
		    destination_point = COALESCE($4, destination_point),
		    fare = COALESCE($5, fare),
		    city = COALESCE($6, city),
		    estimated_time = COALESCE($7, estimated_time),
		    status = COALESCE($8, status),
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := repository.db.Exec(ctx, query, routeID, routeName, startPoint, destPoint, fare, city, estimatedTime, status)
	if err != nil {
		return fmt.Errorf("update route: %w", err)
	}
	return nil
}

func (repository *AdminRepository) DeleteRoute(ctx context.Context, routeID string) error {
	const query = `DELETE FROM routes WHERE id = $1`
	result, err := repository.db.Exec(ctx, query, routeID)
	if err != nil {
		return fmt.Errorf("delete route: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("route not found")
	}
	return nil
}

func (repository *AdminRepository) CreateFare(ctx context.Context, routeID, transportType string, amount float64) (string, error) {
	const query = `
		INSERT INTO fares (route_id, transport_type, amount)
		VALUES ($1, $2, $3)
		RETURNING id::text
	`
	var fareID string
	err := repository.db.QueryRow(ctx, query, routeID, transportType, amount).Scan(&fareID)
	if err != nil {
		return "", fmt.Errorf("create fare: %w", err)
	}
	return fareID, nil
}

func (repository *AdminRepository) UpdateFare(ctx context.Context, fareID string, amount *float64) error {
	const query = `
		UPDATE fares
		SET amount = COALESCE($2, amount),
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := repository.db.Exec(ctx, query, fareID, amount)
	if err != nil {
		return fmt.Errorf("update fare: %w", err)
	}
	return nil
}

func (repository *AdminRepository) DeleteFare(ctx context.Context, fareID string) error {
	const query = `DELETE FROM fares WHERE id = $1`
	result, err := repository.db.Exec(ctx, query, fareID)
	if err != nil {
		return fmt.Errorf("delete fare: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("fare not found")
	}
	return nil
}

func (repository *AdminRepository) GetAnalytics(ctx context.Context) (map[string]interface{}, error) {
	analytics := make(map[string]interface{})

	// Total users
	var totalUsers int
	err := repository.db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&totalUsers)
	if err != nil {
		return nil, fmt.Errorf("count users: %w", err)
	}
	analytics["total_users"] = totalUsers

	// Total stops
	var totalStops int
	err = repository.db.QueryRow(ctx, "SELECT COUNT(*) FROM transit_stops").Scan(&totalStops)
	if err != nil {
		return nil, fmt.Errorf("count stops: %w", err)
	}
	analytics["total_stops"] = totalStops

	// Verified stops
	var verifiedStops int
	err = repository.db.QueryRow(ctx, "SELECT COUNT(*) FROM transit_stops WHERE verified = true").Scan(&verifiedStops)
	if err != nil {
		return nil, fmt.Errorf("count verified stops: %w", err)
	}
	analytics["verified_stops"] = verifiedStops

	// Total routes
	var totalRoutes int
	err = repository.db.QueryRow(ctx, "SELECT COUNT(*) FROM routes WHERE status = 'Active'").Scan(&totalRoutes)
	if err != nil {
		return nil, fmt.Errorf("count routes: %w", err)
	}
	analytics["total_routes"] = totalRoutes

	// Pending reports
	var pendingReports int
	err = repository.db.QueryRow(ctx, "SELECT COUNT(*) FROM crowdsource_reports WHERE status = 'pending'").Scan(&pendingReports)
	if err != nil {
		return nil, fmt.Errorf("count pending reports: %w", err)
	}
	analytics["pending_reports"] = pendingReports

	// Total trips
	var totalTrips int
	err = repository.db.QueryRow(ctx, "SELECT COUNT(*) FROM trips").Scan(&totalTrips)
	if err != nil {
		return nil, fmt.Errorf("count trips: %w", err)
	}
	analytics["total_trips"] = totalTrips

	// Top contributors
	const topContributorsQuery = `
		SELECT u.name, u.contribution_score
		FROM users u
		WHERE u.contribution_score > 0
		ORDER BY u.contribution_score DESC
		LIMIT 5
	`
	rows, err := repository.db.Query(ctx, topContributorsQuery)
	if err != nil {
		return nil, fmt.Errorf("query top contributors: %w", err)
	}
	defer rows.Close()

	var contributors []map[string]interface{}
	for rows.Next() {
		var name string
		var score int
		if err := rows.Scan(&name, &score); err != nil {
			return nil, fmt.Errorf("scan contributor: %w", err)
		}
		contributors = append(contributors, map[string]interface{}{
			"name":  name,
			"score": score,
		})
	}
	analytics["top_contributors"] = contributors

	return analytics, nil
}
