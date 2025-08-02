package storage

import (
	"context"
	"fmt"
	"github.com/AndreySirin/-Effective-Mobile-/internal/entity"
	"github.com/google/uuid"
)

func (s *Storage) CreateSubs(ctx context.Context, subs *entity.Subscription) (uuid.UUID, error) {
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO subscription(serviceName, price, userID, expiresDate)
		 VALUES ($1, $2, $3, $4)
		 RETURNING subscriptionId`,
		subs.ServiceName, subs.Price, subs.UserId, subs.ExpiresDate,
	).Scan(&subs.SubsID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create subscription: %w", err)
	}
	return subs.SubsID, nil
}

func (s *Storage) ReadSubs(ctx context.Context, subsID uuid.UUID) (*entity.Subscription, error) {
	var subs entity.Subscription
	err := s.db.QueryRowContext(ctx,
		`SELECT serviceName, price, userID, startDate,expiresDate 
		FROM subscription 
		WHERE subscriptionId = $1`, subsID).Scan(&subs.ServiceName, &subs.Price, &subs.UserId, &subs.StartDate, &subs.ExpiresDate)
	if err != nil {
		return nil, fmt.Errorf("error for query:%v", err)
	}
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
		return fmt.Errorf("subscription does not exist")
	}
	r, err := tx.Exec(`UPDATE subscription
	SET serviceName=$1, price=$2, userID=$3, expiresDate=$4 
	WHERE subscriptionId=$5`,
		subs.ServiceName,
		subs.Price,
		subs.UserId,
		subs.ExpiresDate,
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
	rows, err := s.db.QueryContext(ctx, `SELECT subscriptionId,serviceName,price,userID,startDate,expiresDate FROM subscription`)
	if err != nil {
		return nil, fmt.Errorf("listing subscriptions: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var sub entity.Subscription
		err = rows.Scan(&sub.SubsID, &sub.ServiceName, &sub.Price, &sub.UserId, &sub.StartDate, &sub.ExpiresDate)
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
