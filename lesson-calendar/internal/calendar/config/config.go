package config

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
