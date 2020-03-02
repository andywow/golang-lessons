package config

// Config application config
type Config struct {
	GRPCListen  string `mapstructure:"grpc_listen"`
	HTTPListen  string `mapstructure:"http_listen"`
	LogFile     string `mapstructure:"log_file"`
	LogLevel    string `mapstructure:"log_level"`
	LogStdout   bool   `mapstructure:"log_console"`
	StorageType string `mapstructure:"storage_type"`
}
