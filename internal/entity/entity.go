package entity

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type Subscription struct {
	SubsID      uuid.UUID  `json:"subsId"`
	ServiceName string     `json:"serviceName"`
	Price       int        `json:"price"`
	UserId      uuid.UUID  `json:"userId"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
}

type SubsRequest struct {
	ServiceName string `json:"serviceName"`
	Price       int    `json:"price"`
	UserId      string `json:"userId"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

type TotalCost struct {
	ServiceName string    `json:"serviceName"`
	UserId      uuid.UUID `json:"userId"`
	Date1       time.Time `json:"date_1"`
	Date2       time.Time `json:"date_2"`
}

type TotalCostRequest struct {
	ServiceName string `json:"serviceName"`
	UserId      string `json:"userId"`
	Date1       string `json:"date_1"`
	Date2       string `json:"date_2"`
}

func SubsToDataBase(lg *slog.Logger, req SubsRequest) (Subscription, error) {
	lg = lg.With("module", "converter")
	lg.Info("converting subscription request to database model", "user_id", req.UserId, "service_name", req.ServiceName)

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		lg.Error("failed to parse start date", "start_date", req.StartDate, "err", err)
		return Subscription{}, fmt.Errorf("error parsing start date: %v", err)
	}

	var endDatePtr *time.Time
	if req.EndDate != "" {
		endDate, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			lg.Error("failed to parse end date", "end_date", req.EndDate, "err", err)
			return Subscription{}, fmt.Errorf("error parsing end date: %v", err)
		}
		endDatePtr = &endDate
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		lg.Error("failed to parse user id", "user_id", req.UserId, "err", err)
		return Subscription{}, fmt.Errorf("error parsing user id: %v", err)
	}

	subs := Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserId:      userId,
		StartDate:   startDate,
		EndDate:     endDatePtr,
	}

	endDateStr := "nil"
	if subs.EndDate != nil {
		endDateStr = subs.EndDate.Format("2006-01")
	}

	lg.Info("subscription request converted successfully",
		"user_id", subs.UserId,
		"service_name", subs.ServiceName,
		"start_date", subs.StartDate.Format("2006-01"),
		"end_date", endDateStr,
	)

	return subs, nil
}

func TotalCostToDataBase(lg *slog.Logger, req TotalCostRequest) (TotalCost, error) {
	lg = lg.With("module", "converter")
	lg.Info("converting total cost request to database model",
		"user_id", req.UserId,
		"service_name", req.ServiceName,
		"date1", req.Date1,
		"date2", req.Date2,
	)

	date1, err := time.Parse("01-2006", req.Date1)
	if err != nil {
		lg.Error("failed to parse start date", "date1", req.Date1, "err", err)
		return TotalCost{}, fmt.Errorf("error parsing start date: %v", err)
	}

	date2, err := time.Parse("01-2006", req.Date2)
	if err != nil {
		lg.Error("failed to parse end date", "date2", req.Date2, "err", err)
		return TotalCost{}, fmt.Errorf("error parsing end date: %v", err)
	}

	if date2.Before(date1) {
		lg.Error("invalid date range: end date before start date", "date1", date1, "date2", date2)
		return TotalCost{}, fmt.Errorf("invalid date range: t2 before t1")
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		lg.Error("failed to parse user id", "user_id", req.UserId, "err", err)
		return TotalCost{}, fmt.Errorf("error parsing user id: %v", err)
	}

	T := TotalCost{
		ServiceName: req.ServiceName,
		UserId:      userId,
		Date1:       date1,
		Date2:       date2,
	}

	lg.Info("total cost request converted successfully",
		"user_id", T.UserId,
		"service_name", T.ServiceName,
		"date1", T.Date1.Format("2006-01"),
		"date2", T.Date2.Format("2006-01"),
	)
	return T, nil
}
