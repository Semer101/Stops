package models

import "time"

type Route struct {
	ID               string    `json:"id"`
	RouteName        string    `json:"routeName"`
	StartPoint       string    `json:"startPoint"`
	DestinationPoint string    `json:"destinationPoint"`
	Fare             float64   `json:"fare"`
	City             string    `json:"city"`
	EstimatedTime    int       `json:"estimatedTime"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type Fare struct {
	ID            string    `json:"id"`
	RouteID       string    `json:"routeId"`
	TransportType string    `json:"transportType"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
