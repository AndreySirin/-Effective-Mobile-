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
}

type SubsRequest struct {
	ServiceName string `json:"serviceName"`
	Price       int    `json:"price"`
	UserId      string `json:"userId"`
	StartDate   string `json:"startDate"`
}

func ToDataBase(req SubsRequest) (Subscription, error) {

	date, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return Subscription{}, fmt.Errorf("error parsing start date: %v", err)
	}
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return Subscription{}, fmt.Errorf("error parsing user id: %v", err)
	}
	subs := Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserId:      userId,
		StartDate:   date,
	}
	return subs, nil
}
