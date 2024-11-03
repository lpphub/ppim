package model

type User struct {
	UID    int64  `bson:"uid"`
	Name   string `bson:"name"`
	Mobile string `bson:"mobile"`
	Email  string `bson:"email"`
}
