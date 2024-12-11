package conf

type Config struct {
	Server serverConf
	Kafka  kafkaConf
}

type serverConf struct {
	Tcp  string
	Ws   string
	Etcd string
}

type kafkaConf struct {
	Brokers []string
	Topic   string
	GroupId string
}
