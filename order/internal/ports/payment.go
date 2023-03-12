package ports

import "github.com/kangkyu/microservices/order/internal/application/domain"

type PaymentPort interface {
	Charge(*domain.Order) error
}
