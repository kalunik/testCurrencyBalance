package consumer

import (
	"fmt"
	"github.com/kalunik/testCurrencyBalance/config"
	"github.com/kalunik/testCurrencyBalance/pkg/log"
	"github.com/nats-io/stan.go"
)

func NewNatsConnection(conf config.AppConfig, logger log.Logger) (stan.Conn, error) {
	natsUrl := fmt.Sprintf("nats://%s%s", conf.Nats.Host, conf.Nats.Port)
	sc, err := stan.Connect(conf.Nats.Cluster, conf.Nats.Client, stan.NatsURL(natsUrl))
	if err != nil {
		return nil, fmt.Errorf("can't connect: %w", err)
	}
	return sc, nil
}
