package db

import (
	"gopkg.in/mgo.v2"
)

type MongoSession struct {
	db *mgo.Database
}

var mgoS = new(MongoSession)

func MongoInit(conn, dbName string) {
	session, e := mgo.Dial(conn)
	if e != nil {
		panic(e)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	e = session.Ping()
	if e != nil {
		panic(e)
	}
	mgoS.db = session.DB(dbName)

	InitMongoIndex()
}

func C(name string) *mgo.Collection {
	return mgoS.db.C(name)
}

func Close() {
	mgoS.db.Session.Close()
}

func InitMongoIndex() {
	e := mgoS.db.C("btc_raw_transaction").EnsureIndexKey("txid")
	if e != nil {
		panic(e)
	}
}
