package configs

import (
	"github.com/spf13/viper"
	"log"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type MongoConfig struct {
	Uri    string `mapstructure:"uri"`
	DBName string `mapstructure:"dbName"`
}

type DBConfig struct {
	Mongo MongoConfig `mapstructure:"mongodb"`
}

type KafkaSASLConfig struct {
	Enable       bool   `mapstructure:"enable"`
	Algorithm    string `mapstructure:"algorithm"`
	SASLUserName string `mapstructure:"user"`
	SASLPassword string `mapstructure:"password"`
}

type KafkaTlsConfig struct {
	Enable             bool `mapstructure:"enable"`
	InsecureSkipVerify bool `mapstructure:"enable"`
}

type KafkaProducerConfig struct {
	Idempotent bool `mapstructure:"idempotent"`
	RetryMax   int  `mapstructure:"retryMax"`
}

type KafkaConsumerConfig struct {
	Assignor          string `mapstructure:"assignor"`
	Oldest            bool   `mapstructure:"oldest"`
	Verbose           int    `mapstructure:"verbose"`
	ConsumerGroupName string `mapstructure:"consumerGroupName"`
	ConsumerTopics    string `mapstructure:"topics"`
}

type KafkaConfig struct {
	// mode: producer|consumer|both
	Mode string `mapstructure:"mode"`
	// brokers: comma separated list
	Brokers  string              `mapstructure:"brokers"`
	SASL     KafkaSASLConfig     `mapstructure:"sasl"`
	TLS      KafkaTlsConfig      `mapstructure:"tls"`
	Producer KafkaProducerConfig `mapstructure:"producer"`
	Consumer KafkaConsumerConfig `mapstructure:"consumer"`
}

type AppConfig struct {
	AppName   string       `mapstructure:"appName"`
	APIServer ServerConfig `mapstructure:"server"`
	Database  DBConfig     `mapstructure:"database"`
	Kafka     KafkaConfig  `mapstructure:"kafka"`
}

var MainConfig AppConfig

func Initialize() error {
	//log.Println("[CONFIG] trying to load configuration file")
	viper.SetConfigType("json")
	viper.AddConfigPath("./")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	err := viper.Unmarshal(&MainConfig)
	if err != nil {
		return err
	}

	log.Println("[INIT] configuration >> loaded")
	return nil
}
