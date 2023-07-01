package customerserver

type Config struct {
	HttpAddr          string `toml:"http_addr"`
	GrpcAddr          string `toml:"grpc_addr"`
	GrpcClAddr        string `toml:"grpc_cl_addr"`
	LogLevel          string `toml:"log_level"`
	DatabaseURL       string `toml:"database_url"`
	BrokerAddr        string `toml:"broker_addr"`
	OrderCreatedTopic string `toml:"order_created"`
}

func NewConfig() *Config {
	return &Config{
		HttpAddr:          ":8080",
		GrpcAddr:          ":8081",
		GrpcClAddr:        ":8081",
		LogLevel:          "debug",
		BrokerAddr:        "9092",
		OrderCreatedTopic: "order_created",
	}
}
