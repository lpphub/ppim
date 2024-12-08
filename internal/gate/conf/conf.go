package conf

type Config struct {
	Server serverConf
	Kafka  kafkaConf
}

type serverConf struct {
	Tcp  string `yml:"tcp"`
	Etcd string `yml:"etcd"`
}

type kafkaConf struct {
	Brokers []string `yml:"brokers"`
}
