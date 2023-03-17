package payment

import (
	"context"

	"github.com/huseyinbabal/microservices-proto/golang/payment"
	"github.com/kangkyu/microservices/order/internal/application/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	payment payment.PaymentClient
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	// data model for connection configurations
	var opts []grpc.DialOption

	// disabling TLS, to explain concepts simple
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// connect to service
	conn, err := grpc.Dial(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}

	// initialize new payment stub instance
	client := payment.NewPaymentClient(conn)

	return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	_, err := a.payment.Create(context.Background(), &payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})

	return err
}
