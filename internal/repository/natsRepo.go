package repository

import (
	"fmt"
	"github.com/kalunik/testCurrencyBalance/config"
	"github.com/nats-io/stan.go"
)

type NatsRepository interface {
	Publish(data []byte) error
	Subscribe(chan []byte) (stan.Subscription, error)
}

type natsRepo struct {
	natsWriter stan.Conn
	natsReader stan.Conn
	conf       config.AppConfig
}

func (n natsRepo) Publish(data []byte) error {
	err := n.natsWriter.Publish(n.conf.Nats.Subject, data)
	if err != nil {
		return fmt.Errorf("can't publish: %w", err)
	}
	return nil
}

func (n natsRepo) Subscribe(dataChan chan []byte) (stan.Subscription, error) {
	sub, err := n.natsReader.Subscribe(n.conf.Nats.Subject,
		func(msg *stan.Msg) { dataChan <- msg.Data },
		stan.DurableName(n.conf.Nats.DurableName))
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func NewNatsRepository(nsWriter stan.Conn, nsReader stan.Conn, conf config.AppConfig) NatsRepository {
	return &natsRepo{natsWriter: nsWriter, natsReader: nsReader, conf: conf}
}
