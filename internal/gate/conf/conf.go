package conf

type Config struct {
	Server serverConf
	Kafka  kafkaConf
	Redis  redisConf
}

type serverConf struct {
	Tcp         string
	Ws          string
	RpcRegistry string
}

type kafkaConf struct {
	Brokers []string
	Topic   string
	GroupId string
}

type redisConf struct {
	Addr     string `yml:"addr"`
	Password string `yml:"password"`
	DB       int    `yml:"db"`
}
