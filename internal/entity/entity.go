package entity

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	SubsID      uuid.UUID  `json:"subsId"`
	ServiceName string     `json:"serviceName"`
	Price       int        `json:"price"`
	UserId      uuid.UUID  `json:"userId"`
	StartDate   time.Time  `json:"startDate"`
	ExpiresDate *time.Time `json:"expiresDate"`
}
