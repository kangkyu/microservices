package ports

import (
	"github.com/kangkyu/microservices/order/internal/application/domain"
)

type DBPort interface {
	Get(id string) (domain.Order, error)
	Save(*domain.Order) error // save Order domain into database
}
