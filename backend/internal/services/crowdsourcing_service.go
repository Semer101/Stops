package services

import (
	"context"
	"errors"

	"stops/backend/internal/repositories"
)

var (
	ErrSubmissionRateLimitExceeded = errors.New("rate limit exceeded: maximum 5 pending submissions per 24 hours")
	ErrInvalidVoteValue            = errors.New("vote value must be 1 or -1")
)

type CrowdsourcingService struct {
	repo *repositories.CrowdsourcingRepository
}

func NewCrowdsourcingService(repo *repositories.CrowdsourcingRepository) *CrowdsourcingService {
	return &CrowdsourcingService{repo: repo}
}

func (service *CrowdsourcingService) ReportStop(ctx context.Context, name, stopType string, lat, lng float64, area string, userID string) (string, error) {
	// Check the 24-hour rate limit
	count, err := service.repo.GetPendingSubmissionsCount(ctx, userID)
	if err != nil {
		return "", err
	}

	if count >= 5 {
		return "", ErrSubmissionRateLimitExceeded
	}

	// Create stop report
	return service.repo.CreateStopReport(ctx, name, stopType, lat, lng, area, userID)
}

func (service *CrowdsourcingService) VoteStop(ctx context.Context, reportID, userID string, voteValue int) (int, error) {
	if voteValue != 1 && voteValue != -1 {
		return 0, ErrInvalidVoteValue
	}

	return service.repo.VoteReport(ctx, reportID, userID, voteValue)
}
