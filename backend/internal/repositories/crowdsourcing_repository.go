package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CrowdsourcingRepository struct {
	db *pgxpool.Pool
}

type CrowdsourceReport struct {
	ID        string    `json:"id"`
	StopID    string    `json:"stopId"`
	Status    string    `json:"status"`
	NetVotes  int       `json:"netVotes"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewCrowdsourcingRepository(db *pgxpool.Pool) *CrowdsourcingRepository {
	return &CrowdsourcingRepository{db: db}
}

func (repository *CrowdsourcingRepository) GetPendingSubmissionsCount(ctx context.Context, userID string) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM crowdsource_reports
		WHERE created_by = $1
		  AND status = 'pending'
		  AND created_at >= NOW() - INTERVAL '24 hours'
	`
	var count int
	err := repository.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count pending submissions: %w", err)
	}
	return count, nil
}

func (repository *CrowdsourcingRepository) CreateStopReport(ctx context.Context, name, stopType string, lat, lng float64, area string, userID string) (string, error) {
	tx, err := repository.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin create stop report transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Insert into transit_stops
	const insertStopQuery = `
		INSERT INTO transit_stops (name, type, location, area, verified, votes, availability_status, created_by)
		VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326), $5, FALSE, 0, 'Unknown', $6)
		RETURNING id::text
	`
	var stopID string
	err = tx.QueryRow(ctx, insertStopQuery, name, stopType, lng, lat, area, userID).Scan(&stopID)
	if err != nil {
		return "", fmt.Errorf("insert transit stop: %w", err)
	}

	// 2. Insert into crowdsource_reports
	const insertReportQuery = `
		INSERT INTO crowdsource_reports (stop_id, status, net_votes, created_by)
		VALUES ($1, 'pending', 0, $2)
		RETURNING id::text
	`
	var reportID string
	err = tx.QueryRow(ctx, insertReportQuery, stopID, userID).Scan(&reportID)
	if err != nil {
		return "", fmt.Errorf("insert crowdsource report: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit create stop report: %w", err)
	}

	return reportID, nil
}

func (repository *CrowdsourcingRepository) VoteReport(ctx context.Context, reportID, userID string, voteValue int) (int, error) {
	tx, err := repository.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin vote transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Insert or update the vote in crowdsource_votes
	const upsertVoteQuery = `
		INSERT INTO crowdsource_votes (report_id, user_id, vote_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (report_id, user_id)
		DO UPDATE SET vote_value = EXCLUDED.vote_value, updated_at = NOW()
	`
	_, err = tx.Exec(ctx, upsertVoteQuery, reportID, userID, voteValue)
	if err != nil {
		return 0, fmt.Errorf("upsert vote: %w", err)
	}

	// 2. Calculate the new net_votes for this report
	const sumVotesQuery = `
		SELECT COALESCE(SUM(vote_value), 0)
		FROM crowdsource_votes
		WHERE report_id = $1
	`
	var netVotes int
	err = tx.QueryRow(ctx, sumVotesQuery, reportID).Scan(&netVotes)
	if err != nil {
		return 0, fmt.Errorf("calculate net votes: %w", err)
	}

	// 3. Update the crowdsource_reports table
	const updateReportQuery = `
		UPDATE crowdsource_reports
		SET net_votes = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING stop_id::text
	`
	var stopID string
	err = tx.QueryRow(ctx, updateReportQuery, netVotes, reportID).Scan(&stopID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("report not found")
		}
		return 0, fmt.Errorf("update report net votes: %w", err)
	}

	// 4. Update the transit_stops votes count
	const updateStopQuery = `
		UPDATE transit_stops
		SET votes = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.Exec(ctx, updateStopQuery, netVotes, stopID)
	if err != nil {
		return 0, fmt.Errorf("update transit stop votes: %w", err)
	}

	// 5. If net votes >= 5, auto-verify stop and report
	if netVotes >= 5 {
		const verifyStopQuery = `
			UPDATE transit_stops
			SET verified = TRUE, updated_at = NOW()
			WHERE id = $1
		`
		_, err = tx.Exec(ctx, verifyStopQuery, stopID)
		if err != nil {
			return 0, fmt.Errorf("auto-verify stop: %w", err)
		}

		const verifyReportQuery = `
			UPDATE crowdsource_reports
			SET status = 'verified', updated_at = NOW()
			WHERE id = $1
		`
		_, err = tx.Exec(ctx, verifyReportQuery, reportID)
		if err != nil {
			return 0, fmt.Errorf("auto-verify report: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit vote transaction: %w", err)
	}

	return netVotes, nil
}
