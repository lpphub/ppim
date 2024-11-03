package conf

type Config struct {
	Mongo mongoConf
	Redis redisConf
	Log   logConf
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
