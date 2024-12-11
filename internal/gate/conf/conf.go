package conf

type Config struct {
	Server serverConf
	Kafka  kafkaConf
}

type serverConf struct {
	Tcp  string `yml:"tcp"`
	Ws   string `yml:"ws"`
	Etcd string `yml:"etcd"`
}

type kafkaConf struct {
	Brokers []string `yml:"brokers"`
	Topic   string   `yml:"topic"`
	GroupID string   `yml:"group"`
}
