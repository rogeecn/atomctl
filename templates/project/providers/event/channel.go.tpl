package event

import "go.ipao.vip/atom/contracts"

const (
	Go    contracts.Channel = "go"
	Kafka contracts.Channel = "kafka"
	Redis contracts.Channel = "redis"
	Sql   contracts.Channel = "sql"
)

type DefaultPublishTo struct{}

func (d *DefaultPublishTo) PublishTo() (contracts.Channel, string) {
	return Go, "event:processed"
}

type DefaultChannel struct{}

func (d *DefaultChannel) Channel() contracts.Channel { return Go }

// kafka
type KafkaChannel struct{}

func (k *KafkaChannel) Channel() contracts.Channel { return Kafka }

// kafka
type RedisChannel struct{}

func (k *RedisChannel) Channel() contracts.Channel { return Redis }
