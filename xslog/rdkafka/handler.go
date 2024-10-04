package rdkafka

import "C"
import (
	"context"
	"encoding/json"
	librdkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/orinchen/xlib/xtool"
	slogcommon "github.com/samber/slog-common"
	"log/slog"
	"time"
)

type Option struct {
	slog.HandlerOptions
	// optional: customize Kafka event builder
	Converter Converter
	// optional: fetch attributes from context
	AttrFromContext []func(ctx context.Context) []slog.Attr

	Kafka *KafkaConfig
}

type KafkaConfig struct {
	Brokers                string
	RequestTimeout         int
	SocketTimeout          int
	MetadataRequestTimeout int
	Retries                int
}

func (o Option) NewKafkaHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	if o.Converter == nil {
		o.Converter = DefaultConverter
	}

	if o.AttrFromContext == nil {
		o.AttrFromContext = []func(ctx context.Context) []slog.Attr{}
	}

	if o.Kafka == nil {
		panic("kafka config is nil")
	}

	kafka, _ := librdkafka.NewProducer(&librdkafka.ConfigMap{
		"bootstrap.servers":           o.Kafka.Brokers,
		"acks":                        "all",
		"compression.type":            "snappy",
		"request.timeout.ms":          o.Kafka.RequestTimeout * 1000, // default 5000ms
		"retries":                     o.Kafka.Retries,
		"debug":                       "broker, topic, metadata", //all, generic, b
		"message.max.bytes":           3000000,                   //default 1000000,
		"socket.keepalive.enable":     "true",
		"socket.timeout.ms":           o.Kafka.SocketTimeout * 1000, // Timeout for network
		"socket.max.fails":            4,                            // default 3
		"heartbeat.interval.ms":       1000,
		"metadata.max.age.ms":         o.Kafka.MetadataRequestTimeout * 1000,
		"metadata.request.timeout.ms": 30000,
	})
	return &KafkaHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
		kafka:  kafka,
	}
}

var _ slog.Handler = (*KafkaHandler)(nil)

type KafkaHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string
	kafka  *librdkafka.Producer
}

func (h *KafkaHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *KafkaHandler) Handle(ctx context.Context, record slog.Record) error {
	fromContext := slogcommon.ContextExtractor(ctx, h.option.AttrFromContext)
	payload := h.option.Converter(h.option.AddSource, h.option.ReplaceAttr, append(h.attrs, fromContext...), h.groups, &record)

	return h.publish(record.Time, payload)
}

func (h *KafkaHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &KafkaHandler{
		option: h.option,
		attrs:  slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *KafkaHandler) WithGroup(name string) slog.Handler {
	return &KafkaHandler{
		option: h.option,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}

func (h *KafkaHandler) publish(timestamp time.Time, payload map[string]interface{}) error {
	key, err := timestamp.MarshalBinary()
	if err != nil {
		return err
	}

	values, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return h.kafka.Produce(&librdkafka.Message{
		TopicPartition: librdkafka.TopicPartition{
			Topic:     xtool.P(""),
			Partition: librdkafka.PartitionAny,
		},
		Key:   key,
		Value: values,
	}, nil)
}
