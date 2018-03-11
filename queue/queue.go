package queue

import (
	"sync"
)

type Queue struct {
	QH       QueueHandler
	ErrCount *uint64

	Topic string
	Chan  chan []byte
}

type QueueHandler interface {
	setup() error
	add(string, chan []byte, *uint64)
}

func New(name string) *Queue {
	var register = map[string]QueueHandler{
		"nsq":   new(NSQ),
		"kafka": new(Kafka),
	}

	return &Queue{
		QH: register[name],
	}
}

func (q *Queue) Run() error {
	var (
		wg  sync.WaitGroup
		err error
	)

	err = q.QH.setup()
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		topic := q.Topic
		q.QH.add(topic, q.Chan, q.ErrCount)
	}()

	wg.Wait()
	return nil
}

func (q *Queue) End() {
	close(q.Chan)
}