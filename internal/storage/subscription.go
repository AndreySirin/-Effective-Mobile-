package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/AndreySirin/-Effective-Mobile-/internal/entity"
	"github.com/google/uuid"
)

type SubscriptionStorage interface {
	CreateSubs(ctx context.Context, subs *entity.Subscription) (uuid.UUID, error)
	ReadSubs(ctx context.Context, subsID uuid.UUID) (*entity.Subscription, error)
	UpdateSubs(ctx context.Context, subsID uuid.UUID, subs *entity.Subscription) error
	DeleteSubs(ctx context.Context, subsID uuid.UUID) error
	ListSubs(ctx context.Context) ([]entity.Subscription, error)
	TotalCost(ctx context.Context, t entity.TotalCost) (int, error)
}

var ErrNotFound = errors.New("subscription not found")

func (s *Storage) CreateSubs(ctx context.Context, subs *entity.Subscription) (uuid.UUID, error) {
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO subscription(serviceName, price, userID)
		 VALUES ($1, $2, $3)
		 RETURNING subscriptionId`,
		subs.ServiceName, subs.Price, subs.UserId,
	).Scan(&subs.SubsID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create subscription: %w", err)
	}
	return subs.SubsID, nil
}

func (s *Storage) ReadSubs(ctx context.Context, subsID uuid.UUID) (*entity.Subscription, error) {
	var subs entity.Subscription
	var exists bool

	err := s.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM subscription WHERE subscriptionId = $1)`, subsID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("read exists check: %w", err)
	}
	if !exists {
		return nil, ErrNotFound
	}

	err = s.db.QueryRowContext(ctx,
		`SELECT serviceName, price, userID, startDate
		FROM subscription 
		WHERE subscriptionId = $1`, subsID).Scan(&subs.ServiceName, &subs.Price, &subs.UserId, &subs.StartDate)
	if err != nil {
		return nil, fmt.Errorf("error for query:%v", err)
	}
	subs.SubsID = subsID
	return &subs, nil
}

func (s *Storage) UpdateSubs(ctx context.Context, subsID uuid.UUID, subs *entity.Subscription) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("creating transactions: %w", err)
	}
	var exists bool
	err = tx.QueryRow(`SELECT EXISTS (SELECT 1 FROM subscription WHERE subscriptionId = $1)`, subsID).Scan(&exists)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("checking for a subscription: %w", err)
	}
	if !exists {
		tx.Rollback()
		return ErrNotFound
	}
	r, err := tx.Exec(`UPDATE subscription
	SET serviceName=$1, price=$2, userID=$3
	WHERE subscriptionId=$4`,
		subs.ServiceName,
		subs.Price,
		subs.UserId,
		subsID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("subscription update request: %w", err)
	}
	rows, err := r.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("checking rows affected: %w", err)
	}
	if rows == 0 {
		tx.Rollback()
		return fmt.Errorf("subscription does not exist")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("making a commit: %w", err)
	}
	return nil
}

func (s *Storage) DeleteSubs(ctx context.Context, subsID uuid.UUID) error {

	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM subscription WHERE subscriptionId = $1)`, subsID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("checking for a subscription: %w", err)
	}
	if !exists {
		return ErrNotFound
	}

	r, err := s.db.ExecContext(ctx, `DELETE FROM subscription WHERE subscriptionId = $1`, subsID)
	if err != nil {
		return fmt.Errorf("deleting a subscription: %w", err)
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("subscription does not exist")
	}
	return nil
}

func (s *Storage) ListSubs(ctx context.Context) ([]entity.Subscription, error) {
	var subs []entity.Subscription
	rows, err := s.db.QueryContext(ctx, `SELECT subscriptionId,serviceName,price,userID,startDate FROM subscription`)
	if err != nil {
		return nil, fmt.Errorf("listing subscriptions: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var sub entity.Subscription
		err = rows.Scan(&sub.SubsID, &sub.ServiceName, &sub.Price, &sub.UserId, &sub.StartDate)
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
	SELECT sum(price) 
	FROM subscription
	where userID = $1
	AND serviceName = $2
	AND startDate BETWEEN $3 AND $4`, t.UserId, t.ServiceName, t.Date1, t.Date2).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("getting total cost: %w", err)
	}
	return sum, nil
}
