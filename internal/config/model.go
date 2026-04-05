package config

import "time"

type AppConfig struct {
	Server   Server   `mapstructure:"server" validate:"required"`
	Secrets  Secrets  `mapstructure:"secrets" validate:"required"`
	DBConfig DBConfig `mapstructure:"db" validate:"required"`
}

type Server struct {
	Address  string `mapstructure:"address" validate:"required"`
	TimeZone string `mapstructure:"time-zone" validate:"required"`
}

type DBConfig struct {
	Host                string        `mapstructure:"host"`
	Port                string        `mapstructure:"port"`
	User                string        `mapstructure:"user"`
	Name                string        `mapstructure:"name"`
	SSLMode             string        `mapstructure:"ssl-mode"`
	MaxOpenConns        *int          `mapstructure:"max-open-conns"`
	MaxIdleConns        *int          `mapstructure:"max-idle-conns"`
	ConnMaxLifetimeHour time.Duration `mapstructure:"conn-max-lifetime-hour"`
}

type Secrets struct {
	DBPassword   string `envconfig:"DB_PASSWORD" required:"true"`
	JwtSecretKey string `envconfig:"JWT_KEY"     required:"true"`
}
