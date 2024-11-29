package conf

type Config struct {
	Grpc grpcConf
}

type grpcConf struct {
	Addr string `yml:"addr"`
}
