package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/kangkyu/microservices/order/internal/application/domain"
)

type Order struct {
	ID         int64     `db:"id"`
	CustomerID int64     `db:"customer_id"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
}

type OrderItem struct {
	ProductCode string  `db:"product_code"`
	UnitPrice   float32 `db:"unit_price"`
	Quantity    int32   `db:"quantity"`
	OrderID     int64   `db:"order_id"`
}

type Adapter struct {
	db *sqlx.DB
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, openErr := sqlx.Open("postgres", dataSourceUrl)
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}

	return &Adapter{db: db}, nil
}

func (a Adapter) Get(id int64) (domain.Order, error) {
	var orderEntity Order
	stmt := "SELECT id, customer_id, status, created_at FROM orders WHERE id = $1 LIMIT 1"
	err := a.db.Get(&orderEntity, stmt, id)
	if err != nil {
		return domain.Order{}, err
	}
	var orderEntityOrderItems []OrderItem
	stmt = "SELECT product_code, unit_price, quantity FROM order_items WHERE order_id = $1"
	err = a.db.Select(&orderEntityOrderItems, stmt, id)
	if err != nil {
		return domain.Order{}, err
	}

	var orderItems []domain.OrderItem
	for _, orderItem := range orderEntityOrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	order := domain.Order{
		ID:         int64(orderEntity.ID),
		CustomerID: orderEntity.CustomerID,
		Status:     orderEntity.Status,
		OrderItems: orderItems,
		CreatedAt:  orderEntity.CreatedAt.UnixNano(),
	}

	return order, nil
}

func (a Adapter) Save(order *domain.Order) error {
	var orderItems []OrderItem
	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
			OrderID:     order.ID,
		})
	}
	orderModel := Order{
		CustomerID: order.CustomerID,
		Status:     order.Status,
	}
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}
	err = tx.Exec("INSERT INTO orders (customer_id, status) VALUES ($1, $2)", orderModel.CustomerID, orderModel.Status)
	if err != nil {
		return err
	}
	err = tx.NamedExec("INSERT INTO order_items (product_code, unit_price, quantity, order_id) VALUES (:product_code, :unit_price, :quantity, :order_id)", orderItems)
	if err != nil {
		return err
	}
	err = tx.Commit()

	return err
}
