package models

import "time"

type SavedPlace struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	PlaceType string    `json:"placeType"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Trip struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	RouteID     *string   `json:"routeId"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	Duration    *int      `json:"duration"`
	Fare        *float64  `json:"fare"`
	Timestamp   time.Time `json:"timestamp"`
}

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID                     string
	Name                   string
	Email                  *string
	Phone                  *string
	PasswordHash           string
	Role                   UserRole
	ContributionScore      int
	HomeLocation           *string
	WorkLocation           *string
	PreferredTransportType *string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type PublicUser struct {
	ID                     string    `json:"id"`
	Name                   string    `json:"name"`
	Email                  *string   `json:"email"`
	Phone                  *string   `json:"phone"`
	Role                   UserRole  `json:"role"`
	ContributionScore      int       `json:"contributionScore"`
	HomeLocation           *string   `json:"homeLocation"`
	WorkLocation           *string   `json:"workLocation"`
	PreferredTransportType *string   `json:"preferredTransportType"`
	CreatedAt              time.Time `json:"createdAt"`
}

func (user User) Public() PublicUser {
	return PublicUser{
		ID:                     user.ID,
		Name:                   user.Name,
		Email:                  user.Email,
		Phone:                  user.Phone,
		Role:                   user.Role,
		ContributionScore:      user.ContributionScore,
		HomeLocation:           user.HomeLocation,
		WorkLocation:           user.WorkLocation,
		PreferredTransportType: user.PreferredTransportType,
		CreatedAt:              user.CreatedAt,
	}
}
