package conf

type Config struct {
	Server serverConf
}

type serverConf struct {
	Tcp  string `yml:"tcp"`
	Etcd string `yml:"etcd"`
}
