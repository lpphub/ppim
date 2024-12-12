package global

import (
	"context"
	"fmt"
	"github.com/lpphub/golib/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDb() {
	logOpts := options.
		Logger().
		SetSink(&mongoLog{}).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", Conf.Mongo.Addr)).
		SetAuth(options.Credential{Username: Conf.Mongo.User, Password: Conf.Mongo.Password}).
		SetLoggerOptions(logOpts)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic("init mongo error: " + err.Error())
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic("ping mongo error: " + err.Error())
	}

	Mongo = client.Database(Conf.Mongo.Database)
}

type mongoLog struct{}

func (*mongoLog) Info(_ int, message string, keysAndValues ...interface{}) {
	logger.Log().Info().Fields(keysAndValues).Msg(message)
}

func (*mongoLog) Error(err error, message string, keysAndValues ...interface{}) {
	logger.Log().Err(err).Fields(keysAndValues).Msg(message)
}
