package producer

import (
	"fmt"
	"github.com/kalunik/testCurrencyBalance/config"
	"github.com/kalunik/testCurrencyBalance/pkg/log"
	"github.com/nats-io/stan.go"
)

func NewNatsConnection(conf config.AppConfig, logger log.Logger) (stan.Conn, error) {
	conn, err := stan.Connect(conf.Nats.Cluster, conf.Nats.DurableName)
	if err != nil {
		return nil, fmt.Errorf("can't connect: %w", err)
	}
	return conn, nil
}
