# microservices

Use the "gRPC Microservices in Go" [book](https://www.manning.com/books/grpc-microservices-in-go) but use [sqlx](https://github.com/jmoiron/sqlx) and postgresql instead of mysql and [gorm](https://github.com/go-gorm/gorm).

### How to use

* `order` service

Need [dbmate](https://github.com/amacneil/dbmate) and [grpcurl](https://github.com/fullstorydev/grpcurl)
```
brew install dbmate
brew install grpcurl
```

Make sure that postgresql db is running
```
brew install postgresql@14
brew services start postgresql@14
```

Run the code written in Go (database url is your local url)
```
$ cd order
$ DATABASE_URL=postgres://kangkyu:postgres@localhost:5432/order_db?sslmode=disable dbmate up
$ DATA_SOURCE_URL=postgres://kangkyu:postgres@localhost:5432/order_db?sslmode=disable APPLICATION_PORT=3000 ENV=development go run cmd/main.go
```

And in another terminal window,
```
$ grpcurl -d '{"user_id": 123, "order_items": [{"product_code": "prod", "quantity": 4, "unit_price": 12}]}' -plaintext localhost:3000 Order/Create
{
  "orderId": "13"
}
```
