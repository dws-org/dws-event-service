package configs

type Server struct {
	GinMode string `mapstructure:"gin_mode"`
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
}
