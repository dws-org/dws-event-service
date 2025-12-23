package configs

// RabbitMQ represents the configuration for RabbitMQ connection
type RabbitMQ struct {
	// Host is the RabbitMQ server hostname or IP address
	Host string `mapstructure:"host"`

	// Port is the RabbitMQ server port (default: 5672 for AMQP)
	Port string `mapstructure:"port"`

	// ManagementPort is the RabbitMQ management UI port (default: 15672)
	ManagementPort string `mapstructure:"management_port"`

	// Username for RabbitMQ authentication
	Username string `mapstructure:"username"`

	// Password for RabbitMQ authentication
	Password string `mapstructure:"password"`

	// VirtualHost is the RabbitMQ virtual host (default: "/")
	VirtualHost string `mapstructure:"virtual_host"`

	// Enabled determines if RabbitMQ integration is enabled
	Enabled bool `mapstructure:"enabled"`
}
