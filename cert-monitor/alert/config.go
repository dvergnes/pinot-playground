package alert

type KafkaConfig struct {
	Topic string `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}