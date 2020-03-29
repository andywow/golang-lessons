package config

import (
	"flag"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// DBConfig database config
type DBConfig struct {
	Database    string `mapstructure:"name"`
	Host        string `mapstructure:"host"`
	MaxConn     string `mapstructure:"max_connections"`
	MaxIdleConn string `mapstructure:"max_idle_connections"`
	Password    string `mapstructure:"password"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
}

// RabbitMQConfig config for rabbit mq
type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Queue    string `mapstructure:"queue"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Retries  int    `mapstructure:"retries"`
	User     string `mapstructure:"user"`
}

// Config application config
type Config struct {
	DB          DBConfig       `mapstructure:"db"`
	GRPCListen  string         `mapstructure:"grpc_listen"`
	HTTPListen  string         `mapstructure:"http_listen"`
	LogFile     string         `mapstructure:"log_file"`
	LogLevel    string         `mapstructure:"log_level"`
	LogStdout   bool           `mapstructure:"log_console"`
	RabbitMQ    RabbitMQConfig `mapstructure:"rabbitmq"`
	StorageType string         `mapstructure:"storage_type"`
}

func init() {
	flag.String("configfile", "config.yaml", "config file path")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

// ParseConfig parse config
func ParseConfig() (*Config, error) {
	viper.SetConfigFile(viper.GetString("configfile"))
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "could not read config")
	}

	// default values
	cfg := Config{
		LogLevel:  "info",
		LogStdout: true,
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "could not parse config")
	}

	return &cfg, nil

}
