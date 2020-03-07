package config

// Config application config
type Config struct {
	DBHost        string `mapstructure:"db_host"`
	DBMaxConn     string `mapstructure:"db_max_connections"`
	DBMaxIdleConn string `mapstructure:"db_max_idle_connections"`
	DBName        string `mapstructure:"db_name"`
	DBPassword    string `mapstructure:"db_password"`
	DBPort        int    `mapstructure:"db_port"`
	DBUser        string `mapstructure:"db_user"`
	GRPCListen    string `mapstructure:"grpc_listen"`
	HTTPListen    string `mapstructure:"http_listen"`
	LogFile       string `mapstructure:"log_file"`
	LogLevel      string `mapstructure:"log_level"`
	LogStdout     bool   `mapstructure:"log_console"`
	StorageType   string `mapstructure:"storage_type"`
}