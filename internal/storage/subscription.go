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
	lg := s.lg.With("module", "storage", "method", "CreateSubs")

	endDateStr := "nil"
	if subs.EndDate != nil {
		endDateStr = subs.EndDate.Format("2006-01")
	}

	lg.Info("creating subscription in database",
		"user_id", subs.UserId,
		"service_name", subs.ServiceName,
		"price", subs.Price,
		"start_date", subs.StartDate.Format("2006-01"),
		"end_date", endDateStr,
	)

	end := interface{}(nil)
	if subs.EndDate != nil {
		end = *subs.EndDate
	}

	err := s.db.QueryRowContext(ctx,
		`INSERT INTO subscription(serviceName, price, userID, startDate, endDate)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING subscriptionId`,
		subs.ServiceName, subs.Price, subs.UserId, subs.StartDate, end,
	).Scan(&subs.SubsID)

	if err != nil {
		lg.Error("failed to create subscription in database", "err", err)
		return uuid.Nil, fmt.Errorf("create subscription: %w", err)
	}

	lg.Info("subscription created successfully", "subscription_id", subs.SubsID)
	return subs.SubsID, nil
}

func (s *Storage) ReadSubs(ctx context.Context, subsID uuid.UUID) (*entity.Subscription, error) {
	lg := s.lg.With("module", "storage", "method", "ReadSubs")
	lg.Info("reading subscription from database", "subscription_id", subsID)

	var subs entity.Subscription
	var startDate time.Time
	var endDate sql.NullTime

	err := s.db.QueryRowContext(ctx,
		`SELECT serviceName, price, userID, startDate, endDate
		 FROM subscription
		 WHERE subscriptionId = $1`, subsID).
		Scan(&subs.ServiceName, &subs.Price, &subs.UserId, &startDate, &endDate)
	subs.SubsID = subsID

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			lg.Info("subscription not found", "subscription_id", subsID)
			return nil, ErrNotFound
		}
		lg.Error("failed to query subscription", "subscription_id", subsID, "err", err)
		return nil, fmt.Errorf("query subscription: %w", err)
	}

	subs.StartDate = startDate
	if endDate.Valid {
		subs.EndDate = &endDate.Time
	} else {
		subs.EndDate = nil
	}

	startDateStr := subs.StartDate.Format("2006-01")
	endDateStr := "nil"
	if subs.EndDate != nil {
		endDateStr = subs.EndDate.Format("2006-01")
	}

	lg.Info("subscription retrieved successfully",
		"subscription_id", subsID,
		"user_id", subs.UserId,
		"service_name", subs.ServiceName,
		"price", subs.Price,
		"start_date", startDateStr,
		"end_date", endDateStr,
	)

	return &subs, nil
}

func (s *Storage) UpdateSubs(ctx context.Context, subsID uuid.UUID, subs *entity.Subscription) error {
	lg := s.lg.With("module", "storage", "method", "UpdateSubs")

	endDateStr := "nil"
	if subs.EndDate != nil {
		endDateStr = subs.EndDate.Format("2006-01")
	}

	lg.Info("updating subscription in database",
		"subscription_id", subsID,
		"user_id", subs.UserId,
		"service_name", subs.ServiceName,
		"price", subs.Price,
		"start_date", subs.StartDate.Format("2006-01"),
		"end_date", endDateStr,
	)

	var end interface{}
	if subs.EndDate != nil {
		end = *subs.EndDate
	} else {
		end = nil
	}

	r, err := s.db.ExecContext(ctx, `UPDATE subscription
	SET serviceName=$1, price=$2, userID=$3, startDate=$4, endDate=$5
	WHERE subscriptionId=$6`,
		subs.ServiceName,
		subs.Price,
		subs.UserId,
		subs.StartDate,
		end,
		subsID,
	)

	if err != nil {
		lg.Error("failed to execute update query", "err", err)
		return fmt.Errorf("update subscription: %w", err)
	}

	rows, err := r.RowsAffected()
	if err != nil {
		lg.Error("failed to get rows affected", "err", err)
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		lg.Info("no subscription updated, not found", "subscription_id", subsID)
		return ErrNotFound
	}

	lg.Info("subscription updated successfully", "subscription_id", subsID)
	return nil
}

func (s *Storage) DeleteSubs(ctx context.Context, subsID uuid.UUID) error {
	lg := s.lg.With("module", "storage", "method", "DeleteSubs")
	lg.Info("deleting subscription from database", "subscription_id", subsID)

	r, err := s.db.ExecContext(ctx, `DELETE FROM subscription WHERE subscriptionId = $1`, subsID)
	if err != nil {
		lg.Error("failed to execute delete query", "subscription_id", subsID, "err", err)
		return fmt.Errorf("deleting a subscription: %w", err)
	}

	rows, err := r.RowsAffected()
	if err != nil {
		lg.Error("failed to get rows affected", "subscription_id", subsID, "err", err)
		return fmt.Errorf("checking rows affected: %w", err)
	}

	if rows == 0 {
		lg.Info("no subscription deleted, not found", "subscription_id", subsID)
		return ErrNotFound
	}

	lg.Info("subscription deleted successfully", "subscription_id", subsID)
	return nil
}

func (s *Storage) ListSubs(ctx context.Context, pointOfReference time.Time) ([]entity.Subscription, error) {
	subs := make([]entity.Subscription, 0, 10)

	rows, err := s.db.QueryContext(ctx, `
        SELECT
            subscriptionId,
            serviceName,
            price,
            userID,
            startDate,
            endDate
        FROM subscription
        WHERE startDate < $1
        ORDER BY startDate DESC
        LIMIT 10
    `, pointOfReference)
	if err != nil {
		return nil, fmt.Errorf("listing subscriptions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sub entity.Subscription
		var end sql.NullTime // используем NullTime для nullable поля

		err = rows.Scan(
			&sub.SubsID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserId,
			&sub.StartDate,
			&end,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription row: %w", err)
		}

		// Конвертируем NullTime в *time.Time
		if end.Valid {
			sub.EndDate = &end.Time
		} else {
			sub.EndDate = nil
		}

		subs = append(subs, sub)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return subs, nil
}

func (s *Storage) TotalCost(ctx context.Context, t entity.TotalCost) (int, error) {
	lg := s.lg.With("module", "storage", "method", "TotalCost")
	lg.Info("calculating total cost for user",
		"user_id", t.UserId,
		"service_name", t.ServiceName,
		"date1", t.Date1.Format("2006-01"),
		"date2", t.Date2.Format("2006-01"),
	)

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
	      AND (endDate IS NULL OR endDate >= $3)
	),
	month_subs AS (
	    SELECT m.month_start, a.price
	    FROM months m
	    JOIN active_subs a 
	      ON m.month_start >= a.startDate
	     AND (a.endDate IS NULL OR m.month_start <= a.endDate)
	),
	unique_months AS (
	    SELECT month_start, MAX(price) AS price
	    FROM month_subs
	    GROUP BY month_start
	)
	SELECT COALESCE(SUM(price), 0) AS total_cost
	FROM unique_months;
	`, t.UserId, t.ServiceName, t.Date1, t.Date2).Scan(&sum)

	if err != nil {
		lg.Error("failed to calculate total cost", "err", err)
		return 0, fmt.Errorf("getting total cost: %w", err)
	}

	lg.Info("total cost calculated successfully",
		"user_id", t.UserId,
		"service_name", t.ServiceName,
		"total_cost", sum,
	)

	return sum, nil
}
