package main

import (
	"os"

	"github.com/kangkyu/microservices/payment/config"
	"github.com/kangkyu/microservices/payment/internal/adapters/db"
	"github.com/kangkyu/microservices/payment/internal/adapters/grpc"
	"github.com/kangkyu/microservices/payment/internal/application/core/api"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(customLogger{
		formatter: logrus.JSONFormatter{FieldMap: logrus.FieldMap{
			"msg": "message",
		}},
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

type customLogger struct {
	formatter logrus.JSONFormatter
}

func (l customLogger) Format(entry *logrus.Entry) ([]byte, error) {
	return l.formatter.Format(entry)
}

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		logrus.Fatalf("Failed to connect to database. Error: %v", err)
	}

	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
