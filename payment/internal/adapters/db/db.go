package db

import (
	"context"
	"fmt"
	"time"

	"github.com/kangkyu/microservices/payment/internal/application/core/domain"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Payment struct {
	ID         int64     `db:"id"`
	CustomerID int64     `db:"customer_id"`
	Status     string    `db:"status"`
	OrderID    int64     `db:"order_id"`
	TotalPrice float32   `db:"total_price"`
	CreatedAt  time.Time `db:"created_at"`
}

type Adapter struct {
	db *sqlx.DB
}

func (a Adapter) Get(ctx context.Context, id string) (domain.Payment, error) {
	var paymentEntity Payment
	stmt := "SELECT id, customer_id, status, order_id, total_price, created_at FROM payments WHERE id = $1 LIMIT 1"
	err := a.db.GetContext(ctx, &paymentEntity, stmt, id)
	if err != nil {
		return domain.Payment{}, err
	}
	payment := domain.Payment{
		ID:         int64(paymentEntity.ID),
		CustomerID: paymentEntity.CustomerID,
		Status:     paymentEntity.Status,
		OrderID:    paymentEntity.OrderID,
		TotalPrice: paymentEntity.TotalPrice,
		CreatedAt:  paymentEntity.CreatedAt.UnixNano(),
	}

	return payment, nil
}

func (a Adapter) Save(ctx context.Context, payment *domain.Payment) error {
	paymentModel := Payment{
		CustomerID: payment.CustomerID,
		Status:     payment.Status,
		OrderID:    payment.OrderID,
		TotalPrice: payment.TotalPrice,
	}
	var paymentID int64
	stmt := "INSERT INTO payments (customer_id, status, order_id, total_price) VALUES (:customer_id, :status, :order_id, :total_price) RETURNING id"
	err := a.db.QueryRowContext(ctx, stmt, paymentModel).Scan(&paymentID)
	if err == nil {
		payment.ID = paymentID
	}

	return err
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, openErr := sqlx.Open("postgres", dataSourceUrl)
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}

	return &Adapter{db: db}, nil
}
