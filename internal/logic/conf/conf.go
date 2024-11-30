package conf

type Config struct {
	Server serverConf
	Mongo  mongoConf
	Redis  redisConf
	Log    logConf
}

type serverConf struct {
	Api  string `yml:"api"`
	Rpc  string `yml:"rpc"`
	Etcd string `yml:"etcd"`
}

type mongoConf struct {
	Addr     string `yml:"addr"`
	Database string `yml:"database"`
	User     string `yml:"user"`
	Password string `yml:"password"`
}

type redisConf struct {
	Addr     string `yml:"addr"`
	Password string `yml:"password"`
	DB       int    `yml:"db"`
}

type logConf struct {
	Path string `yml:"path"`
}
