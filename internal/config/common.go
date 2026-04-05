package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/mstoykov/envconfig"
	"github.com/ory/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Load() AppConfig {
	viper.SetConfigFile("config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("load config: read config.yaml error: %v", err)
	}
	viper.AutomaticEnv()
	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("load config: unmarshal error: %v", err)
	}
	// load secrets
	if err := loadSecrets(&cfg.Secrets, "config/secret.env", "SECRET"); err != nil {
		log.Fatalf("load config: load secrets error: %v", err)
	}
	log.Printf("config loaded: address=%s db=%s@%s:%s/%s",
		cfg.Server.Address,
		cfg.DBConfig.User,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.Name,
	)
	return cfg
}
func loadSecrets(secrets interface{}, secretPath string, prefix string) error {
	if err := godotenv.Load(secretPath); err != nil {
		log.Fatalf("load secrets: read %s error: %v", secretPath, err)
	}
	if err := envconfig.Process(prefix, secrets); err != nil {
		log.Fatalf("load secrets: unmarshal error: %v", err)
	}
	return nil
}

type DBConnectConfig struct {
	Name                string
	SSLMode             string
	MaxOpenConns        *int
	MaxIdleConns        *int
	ConnMaxLifetimeHour time.Duration
	Host                string
	Port                string
	User                string
	Password            string
	TimeZone            string
}

func Connect(config *DBConnectConfig) (*gorm.DB, error) {
	if config.Host == "" {
		return nil, errors.New("host is required")
	}
	if config.Port == "" {
		return nil, errors.New("port is required")
	}
	if config.User == "" {
		return nil, errors.New("user is required")
	}
	if config.Password == "" {
		return nil, errors.New("password is required")
	}
	if config.Name == "" {
		return nil, errors.New("entity name is required")
	}
	if config.SSLMode == "" {
		config.SSLMode = "disable"
	}

	if config.MaxOpenConns == nil {
		config.MaxOpenConns = new(int)
		*config.MaxOpenConns = 20
	}
	if config.MaxIdleConns == nil {
		config.MaxIdleConns = new(int)
		*config.MaxIdleConns = 10
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
		config.SSLMode,
	)

	if config.TimeZone != "" {
		dsn = fmt.Sprintf("%s TimeZone=%s", dsn, config.TimeZone)
	}

	gDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Join(err, errors.New("can't initialize database session"))
	}

	postgresDb, err := gDB.DB()
	if err != nil {
		return nil, errors.Join(err, errors.New("can't get database"))
	}

	postgresDb.SetMaxOpenConns(*config.MaxOpenConns)
	postgresDb.SetMaxIdleConns(*config.MaxIdleConns)
	postgresDb.SetConnMaxLifetime(time.Hour * config.ConnMaxLifetimeHour)

	if err := postgresDb.Ping(); err != nil {
		return nil, errors.Join(err, errors.New("database ping error"))
	}

	return gDB, nil
}
