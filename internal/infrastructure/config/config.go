package config

import "time"

type Config struct {
	App      AppConfig      `yaml:"app"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Log      LogConfig      `yaml:"log"`
	Auth     AuthConfig     `yaml:"auth"`
}

// AppConfig - App params
type AppConfig struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string `yaml:"version"`
}

// ServerConfig - http server params
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// DatabaseConfig - PostgreSQL params
type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	SSLMode         string `yaml:"ssl_mode"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime_minutes"`
}

// RedisConfig - Redis params
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

// LogConfig - logger params
type LogConfig struct {
	Directory        string `yaml:"directory"`
	InfoFile         string `yaml:"info_file"`
	ErrorFile        string `yaml:"error_file"`
	WarnFile         string `yaml:"warn_file"`
	DebugFile        string `yaml:"debug_file"`
	GormFile         string `yaml:"gorm_file"`
	MaxSizeMegabytes int    `yaml:"max_size_megabytes"`
	MaxBackups       int    `yaml:"max_backups"`
	MaxAgeDays       int    `yaml:"max_age_days"`
	Compress         bool   `yaml:"compress"`
	LocalTime        bool   `yaml:"local_time"`
}

// AuthConfig - auth params
type AuthConfig struct {
	HMACAlgorithm string `yaml:"hmac_algorithm"`
}
