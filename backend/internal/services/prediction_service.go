package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type PredictionService struct {
	// Add repository references if needed
}

type PredictionResult struct {
	Type           string    `json:"type"`
	PredictedValue string    `json:"predictedValue"`
	Confidence     float64   `json:"confidence"`
	GeneratedAt    time.Time `json:"generatedAt"`
}

func NewPredictionService() *PredictionService {
	return &PredictionService{}
}

func (service *PredictionService) GetPredictions(ctx context.Context, routeID string) ([]PredictionResult, error) {
	// Seed random number generator based on time and route ID (so we get consistent but varying results)
	seed := time.Now().UnixNano()
	if routeID != "" {
		// simple hash of route ID to seed
		var h int64
		for _, c := range routeID {
			h = h*31 + int64(c)
		}
		seed += h
	}
	r := rand.New(rand.NewSource(seed))

	// Heuristic 1: ETA
	etaMinutes := r.Intn(25) + 5 // between 5 and 30 minutes
	etaVal := fmt.Sprintf("%d minutes", etaMinutes)

	// Heuristic 2: Taxi Availability
	taxiCount := r.Intn(8) // between 0 and 7
	taxiVal := "Busy - No taxis nearby"
	if taxiCount > 0 {
		taxiVal = fmt.Sprintf("Available - %d taxis nearby", taxiCount)
	}

	// Heuristic 3: Congestion
	congestionLevels := []string{"Low", "Medium", "High"}
	congestionVal := congestionLevels[r.Intn(len(congestionLevels))]

	return []PredictionResult{
		{
			Type:           "ETA",
			PredictedValue: etaVal,
			Confidence:     0.85,
			GeneratedAt:    time.Now(),
		},
		{
			Type:           "Taxi",
			PredictedValue: taxiVal,
			Confidence:     0.78,
			GeneratedAt:    time.Now(),
		},
		{
			Type:           "Congestion",
			PredictedValue: congestionVal,
			Confidence:     0.90,
			GeneratedAt:    time.Now(),
		},
	}, nil
}
