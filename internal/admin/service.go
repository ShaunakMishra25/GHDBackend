package admin

import (
	"context"

	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type DashboardResponse struct {
	TodayOrderCount      int64              `json:"today_order_count"`
	TodayRevenue         float64            `json:"today_revenue"`
	PendingOrderCount    int64              `json:"pending_order_count"`
	DispatchedOrderCount int64              `json:"dispatched_order_count"`
	TotalUsers           int64              `json:"total_users"`
	TotalOrders          int64              `json:"total_orders"`
	ActiveProducts       int64              `json:"active_products"`
	StatusCounts         map[string]int64   `json:"status_counts"`
	TodayOrders          []OrderBrief       `json:"today_orders"`
}

type OrderBrief struct {
	ID          int64   `json:"id"`
	Status      string  `json:"status"`
	Total       float64 `json:"total"`
	UserName    string  `json:"user_name"`
	UserPhone   string  `json:"user_phone"`
	AddressText string  `json:"address_text"`
	CreatedAt   string  `json:"created_at"`
}

func (s *Service) GetDashboard(ctx context.Context) (*DashboardResponse, *apperrors.AppError) {
	stats, err := s.repo.GetDashboardStats(ctx)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to fetch dashboard stats", "डैशबोर्ड लोड करने में विफल", err)
	}

	statusCounts, err := s.repo.GetOrderStatusCounts(ctx)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to fetch status counts", "स्थिति गणना लोड करने में विफल", err)
	}

	todayOrders, err := s.repo.GetTodayOrders(ctx, 10, 0)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to fetch today orders", "आज के ऑर्डर लोड करने में विफल", err)
	}

	revenue := float64(0)
	if stats.TodayRevenue.Valid {
		v, _ := stats.TodayRevenue.Float64Value()
		revenue = v.Float64
	}

	sc := make(map[string]int64)
	for _, row := range statusCounts {
		sc[string(row.Status)] = row.Count
	}

	orders := make([]OrderBrief, len(todayOrders))
	for i, o := range todayOrders {
		total := float64(0)
		if o.Total.Valid {
			v, _ := o.Total.Float64Value()
			total = v.Float64
		}
		orders[i] = OrderBrief{
			ID:          o.ID,
			Status:      string(o.Status),
			Total:       total,
			UserName:    o.UserName,
			UserPhone:   o.UserPhone,
			AddressText: o.AddressText,
			CreatedAt:   o.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &DashboardResponse{
		TodayOrderCount:      stats.TodayOrderCount,
		TodayRevenue:         revenue,
		PendingOrderCount:    stats.PendingOrderCount,
		DispatchedOrderCount: stats.DispatchedOrderCount,
		TotalUsers:           stats.TotalUsers,
		TotalOrders:          stats.TotalOrders,
		ActiveProducts:       stats.ActiveProducts,
		StatusCounts:         sc,
		TodayOrders:          orders,
	}, nil
}
