package models

type TransitStopType string

const (
	TransitStopTypeBusStop   TransitStopType = "bus_stop"
	TransitStopTypeTaxiStand TransitStopType = "taxi_stand"
)

type TransitStop struct {
	ID                 string
	Name               string
	Type               TransitStopType
	Latitude           float64
	Longitude          float64
	Area               *string
	Verified           bool
	Votes              int
	AvailabilityStatus string
	DistanceMeters     float64
	WalkingTimeMinutes int
}

type PublicTransitStop struct {
	ID                 string          `json:"id"`
	Name               string          `json:"name"`
	Type               TransitStopType `json:"type"`
	Latitude           float64         `json:"latitude"`
	Longitude          float64         `json:"longitude"`
	Area               *string         `json:"area"`
	Verified           bool            `json:"verified"`
	Votes              int             `json:"votes"`
	AvailabilityStatus string          `json:"availabilityStatus"`
	DistanceMeters     float64         `json:"distanceMeters"`
	WalkingTimeMinutes int             `json:"walkingTimeMinutes"`
}

func (stop TransitStop) Public() PublicTransitStop {
	return PublicTransitStop{
		ID:                 stop.ID,
		Name:               stop.Name,
		Type:               stop.Type,
		Latitude:           stop.Latitude,
		Longitude:          stop.Longitude,
		Area:               stop.Area,
		Verified:           stop.Verified,
		Votes:              stop.Votes,
		AvailabilityStatus: stop.AvailabilityStatus,
		DistanceMeters:     stop.DistanceMeters,
		WalkingTimeMinutes: stop.WalkingTimeMinutes,
	}
}
