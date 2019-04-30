package redis_handler

import (
	"github.com/adjust/rmq"
	"github.com/rajatgpt1521/cachingSystem/service/pkg/cache_handler"
	"github.com/rs/zerolog/log"
	"time"
)

var (
	notificationQueue rmq.Queue
)

func AddMessage(str string)(string) {
	if !notificationQueue.Publish(str) {
		log.Error().Msg("unable to add msg in queue")
		return "unable to add msg in queue"
	}
	return "Successfully notified"
}

func init() {
	connection := rmq.OpenConnection("my service", "tcp", "localhost:6379", 1)

	notificationQueue = connection.OpenQueue("tasks")
	notificationQueue.StartConsuming(30, time.Second)

	notificationQueue.AddConsumerFunc("test", func(delivery rmq.Delivery) {

		p := delivery.Payload()
		delivery.Ack()
		if p == "reload" {
			cache_handler.LoadFromDB()
		}
	})
}
