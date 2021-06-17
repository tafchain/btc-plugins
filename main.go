package main

import (
	"fmt"
	"math/rand"
	"time"
	"vbh/btc-plugins/api/routers"
	"vbh/btc-plugins/blockchain/omni"
	"vbh/btc-plugins/blockchain/tx"
	"vbh/btc-plugins/log"
	"vbh/btc-plugins/persistence"
)

func main() {
	rand.Seed(time.Now().Unix())
	log.InitLogger(conf.Log.Path, "btc-plugins", "debug")
	db.MongoInit(conf.Db.Mongo.Conn, conf.Db.Mongo.Name)
	omni.Init(conf.Omni.Host, conf.Omni.User, conf.Omni.Pass, conf.Omni.PropertyId)
	if conf.Sync.State == 1 {
		tx.Init(conf.Sync.Period, conf.Sync.StartHeight)
	}
	startSever()
}

func startSever() {
	r := routers.Init()
	e := r.Run(fmt.Sprintf(":%d", conf.HttpPort))
	if e != nil {
		log.Panic(e)
		panic(e)
	}
}
