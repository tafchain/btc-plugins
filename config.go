package main

import "github.com/spf13/viper"

func init() {
	v := viper.New()
	v.SetConfigName("btc")
	v.AddConfigPath("conf/")
	v.SetConfigType("yaml")

	if e := v.ReadInConfig(); e != nil {
		panic(e)
	}

	if e := v.Unmarshal(&conf); e != nil {
		panic(e)
	}

}

var conf *Config

type Config struct {
	Log      Log  `yaml:"log"`
	HttpPort int  `yaml:"httpPort"`
	Db       Db   `yaml:"db"`
	Omni     Omni `yaml:"omni"`
	Sync     Sync `yaml:"sync"`
}

type Log struct {
	Path string `yaml:"path"`
}

type Sync struct {
	Period      int   `yaml:"period"`
	StartHeight int64 `yaml:"startHeight"`
	State       int   `yaml:"state"`
}

type Db struct {
	Mongo MongoDb `yaml:"mongo"`
}

type MongoDb struct {
	Conn string `yaml:"conn"`
	Name string `yaml:"name"`
}

type Omni struct {
	Host       string `yaml:"host"`
	User       string `yaml:"user"`
	Pass       string `yaml:"pass"`
	PropertyId uint32 `yaml:"propertyId"`
}
