package services

import (
	"context"

	"stops/backend/internal/repositories"
)

type AdminService struct {
	admin *repositories.AdminRepository
}

func NewAdminService(admin *repositories.AdminRepository) *AdminService {
	return &AdminService{admin: admin}
}

func (service *AdminService) GetPendingReports(ctx context.Context) ([]repositories.PendingReport, error) {
	return service.admin.FindPendingReports(ctx)
}

func (service *AdminService) VerifyReport(ctx context.Context, reportID, adminUserID string) error {
	stopID, err := service.admin.VerifyReport(ctx, reportID)
	if err != nil {
		return err
	}

	metadata := map[string]interface{}{
		"report_id": reportID,
	}

	return service.admin.CreateAuditLog(ctx, adminUserID, "Verify Report", "transit_stops", &stopID, metadata)
}

func (service *AdminService) RejectReport(ctx context.Context, reportID, adminUserID string) error {
	stopID, err := service.admin.RejectReport(ctx, reportID)
	if err != nil {
		return err
	}

	metadata := map[string]interface{}{
		"report_id": reportID,
	}

	return service.admin.CreateAuditLog(ctx, adminUserID, "Reject Report", "transit_stops", &stopID, metadata)
}

func (service *AdminService) GetAuditLogs(ctx context.Context) ([]repositories.AuditLog, error) {
	return service.admin.FindAuditLogs(ctx)
}

func (service *AdminService) CreateStop(ctx context.Context, name, stopType string, lat, lng float64, area string, verified bool) (string, error) {
	return service.admin.CreateStop(ctx, name, stopType, lat, lng, area, verified)
}

func (service *AdminService) UpdateStop(ctx context.Context, stopID string, name, stopType *string, lat, lng *float64, area *string, verified *bool) error {
	return service.admin.UpdateStop(ctx, stopID, name, stopType, lat, lng, area, verified)
}

func (service *AdminService) DeleteStop(ctx context.Context, stopID string) error {
	return service.admin.DeleteStop(ctx, stopID)
}

func (service *AdminService) CreateRoute(ctx context.Context, routeName, startPoint, destPoint string, fare float64, city string, estimatedTime int) (string, error) {
	return service.admin.CreateRoute(ctx, routeName, startPoint, destPoint, fare, city, estimatedTime)
}

func (service *AdminService) UpdateRoute(ctx context.Context, routeID string, routeName, startPoint, destPoint *string, fare *float64, city *string, estimatedTime *int, status *string) error {
	return service.admin.UpdateRoute(ctx, routeID, routeName, startPoint, destPoint, fare, city, estimatedTime, status)
}

func (service *AdminService) DeleteRoute(ctx context.Context, routeID string) error {
	return service.admin.DeleteRoute(ctx, routeID)
}

func (service *AdminService) CreateFare(ctx context.Context, routeID, transportType string, amount float64) (string, error) {
	return service.admin.CreateFare(ctx, routeID, transportType, amount)
}

func (service *AdminService) UpdateFare(ctx context.Context, fareID string, amount *float64) error {
	return service.admin.UpdateFare(ctx, fareID, amount)
}

func (service *AdminService) DeleteFare(ctx context.Context, fareID string) error {
	return service.admin.DeleteFare(ctx, fareID)
}

func (service *AdminService) GetAnalytics(ctx context.Context) (map[string]interface{}, error) {
	return service.admin.GetAnalytics(ctx)
}
