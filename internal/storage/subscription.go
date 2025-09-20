package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/AndreySirin/-Effective-Mobile-/internal/entity"
	"github.com/google/uuid"
	"time"
)

type SubscriptionStorage interface {
	CreateSubs(ctx context.Context, subs *entity.Subscription) (uuid.UUID, error)
	ReadSubs(ctx context.Context, subsID uuid.UUID) (*entity.Subscription, error)
	UpdateSubs(ctx context.Context, subsID uuid.UUID, subs *entity.Subscription) error
	DeleteSubs(ctx context.Context, subsID uuid.UUID) error
	ListSubs(ctx context.Context, time time.Time) ([]entity.Subscription, error)
	TotalCost(ctx context.Context, t entity.TotalCost) (int, error)
}

var ErrNotFound = errors.New("subscription not found")

func (s *Storage) CreateSubs(ctx context.Context, subs *entity.Subscription) (uuid.UUID, error) {
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO subscription(serviceName, price, userID, startDate, endDate)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING subscriptionId`,
		subs.ServiceName, subs.Price, subs.UserId, subs.StartDate, subs.EndDate,
	).Scan(&subs.SubsID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create subscription: %w", err)
	}
	return subs.SubsID, nil
}

func (s *Storage) ReadSubs(ctx context.Context, subsID uuid.UUID) (*entity.Subscription, error) {
	var subs entity.Subscription
	err := s.db.QueryRowContext(ctx,
		`SELECT serviceName, price, userID, startDate, endDate
		 FROM subscription
		 WHERE subscriptionId = $1`, subsID).
		Scan(&subs.ServiceName, &subs.Price, &subs.UserId, &subs.StartDate, &subs.EndDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query subscription: %w", err)
	}
	return &subs, nil
}

func (s *Storage) UpdateSubs(ctx context.Context, subsID uuid.UUID, subs *entity.Subscription) error {

	r, err := s.db.ExecContext(ctx, `UPDATE subscription
	SET serviceName=$1, price=$2, userID=$3, startDate=$4, endDate=$5
	WHERE subscriptionId=$6`,
		subs.ServiceName,
		subs.Price,
		subs.UserId,
		subs.StartDate,
		subs.EndDate,
		subsID)

	if err != nil {
		return fmt.Errorf("update subscription: %w", err)
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Storage) DeleteSubs(ctx context.Context, subsID uuid.UUID) error {

	r, err := s.db.ExecContext(ctx, `DELETE FROM subscription WHERE subscriptionId = $1`, subsID)
	if err != nil {
		return fmt.Errorf("deleting a subscription: %w", err)
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Storage) ListSubs(ctx context.Context, pointOfReference time.Time) ([]entity.Subscription, error) {
	subs := make([]entity.Subscription, 0, 10)
	rows, err := s.db.QueryContext(ctx,
		`SELECT
       subscriptionId,
       serviceName,
       price,
       userID,
       startDate,
       endDate
FROM subscription
WHERE startDate<$1
ORDER BY startDate DESC 
LIMIT 10
`, pointOfReference)
	if err != nil {
		return nil, fmt.Errorf("listing subscriptions: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var sub entity.Subscription
		err = rows.Scan(&sub.SubsID, &sub.ServiceName, &sub.Price, &sub.UserId, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		subs = append(subs, sub)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("listing rows: %w", err)
	}
	return subs, nil
}

func (s *Storage) TotalCost(ctx context.Context, t entity.TotalCost) (int, error) {
	var sum int
	err := s.db.QueryRowContext(ctx, `
WITH months AS (
    SELECT generate_series(
        $3::date, 
        $4::date, 
        interval '1 month'
    )::date AS month_start
),
active_subs AS (
    SELECT serviceName, price, startDate, endDate
    FROM subscription
    WHERE userID = $1
      AND serviceName = $2
      AND startDate <= $4
      AND endDate >= $3
),
month_subs AS (
    SELECT m.month_start, a.price
    FROM months m
    JOIN active_subs a 
      ON m.month_start BETWEEN a.startDate AND a.endDate
),
unique_months AS (
    SELECT month_start, MAX(price) AS price
    FROM month_subs
    GROUP BY month_start
)
SELECT COALESCE(SUM(price), 0) AS total_cost
FROM unique_months
`, t.UserId, t.ServiceName, t.Date1, t.Date2).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("getting total cost: %w", err)
	}
	return sum, nil
}
