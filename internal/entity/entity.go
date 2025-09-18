package entity

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	SubsID      uuid.UUID `json:"subsId"`
	ServiceName string    `json:"serviceName"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"userId"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
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

func SubsToDataBase(req SubsRequest) (Subscription, error) {

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return Subscription{}, fmt.Errorf("error parsing start date: %v", err)
	}
	endDate, err := time.Parse("01-2006", req.EndDate)
	if err != nil {
		return Subscription{}, fmt.Errorf("error parsing end date: %v", err)
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return Subscription{}, fmt.Errorf("error parsing user id: %v", err)
	}
	subs := Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserId:      userId,
		StartDate:   startDate,
		EndDate:     endDate,
	}
	return subs, nil
}

func TotalCostToDataBase(req TotalCostRequest) (TotalCost, error) {
	date1, err := time.Parse("01-2006", req.Date1)
	if err != nil {
		return TotalCost{}, fmt.Errorf("error parsing start date: %v", err)
	}
	date2, err := time.Parse("01-2006", req.Date2)
	if err != nil {
		return TotalCost{}, fmt.Errorf("error parsing start date: %v", err)
	}
	if date2.Before(date1) {
		return TotalCost{}, fmt.Errorf("invalid date range: t2 before t1")
	}
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return TotalCost{}, fmt.Errorf("error parsing user id: %v", err)
	}
	T := TotalCost{
		ServiceName: req.ServiceName,
		UserId:      userId,
		Date1:       date1,
		Date2:       date2,
	}
	return T, nil
}
