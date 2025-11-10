package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func LoadConfig(configPath string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("[config.LoadConfig]: Warning: .env file not found")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("[config.LoadConfig]: failed to read config file: %w", err)
	}

	var AppParams Config
	if err := yaml.Unmarshal(data, &AppParams); err != nil {
		return nil, fmt.Errorf("[config.LoadConfig]: failed to parse config file: %w", err)
	}

	overrideFromEnv(&AppParams)

	if err := validate(&AppParams); err != nil {
		return nil, fmt.Errorf("[config.LoadConfig]: config validation failed: %w", err)
	}

	return &AppParams, nil
}

func overrideFromEnv(AppParams *Config) {
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		AppParams.Database.Host = dbHost
	}

	if dbPassword := os.Getenv("POSTGRES_PASSWORD"); dbPassword != "" {
		AppParams.Database.Password = dbPassword
	}

	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		AppParams.Redis.Host = redisHost
	}

	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		AppParams.Redis.Password = redisPassword
	}

	if env := os.Getenv("APP_ENVIRONMENT"); env != "" {
		AppParams.App.Environment = env
		// Set GIN_MODE based on environment
		if env == EnvironmentProduction {
			AppParams.App.GinMode = GinModeRelease
		}
	}

	if ginMode := os.Getenv("GIN_MODE"); ginMode != "" {
		AppParams.App.GinMode = ginMode
	}
}

func validate(AppParams *Config) error {
	if AppParams.App.Name == "" {
		return fmt.Errorf("[config.validate]: app.name is required")
	}
	if AppParams.App.Environment != EnvironmentDevelopment && AppParams.App.Environment != EnvironmentProduction {
		return fmt.Errorf("[config.validate]: app.environment must be '%s' or '%s'", EnvironmentDevelopment, EnvironmentProduction)
	}

	if AppParams.Server.Port == "" {
		return fmt.Errorf("[config.validate]: server.port is required")
	}
	if AppParams.Database.Host == "" {
		return fmt.Errorf("[config.validate]: database.host is required")
	}
	if AppParams.Database.User == "" {
		return fmt.Errorf("[config.validate]: database.user is required")
	}
	if AppParams.Database.Password == "" {
		return fmt.Errorf("[config.validate]: database.password is required (set DB_PASSWORD in .env)")
	}
	if AppParams.Database.Database == "" {
		return fmt.Errorf("[config.validate]: database.database is required")
	}

	if AppParams.Redis.Host == "" {
		return fmt.Errorf("[config.validate]: redis.host is required")
	}

	if AppParams.Log.Directory == "" {
		return fmt.Errorf("[config.validate]: log.directory is required")
	}

	if AppParams.Auth.HMACAlgorithm != "sha1" && AppParams.Auth.HMACAlgorithm != "sha256" {
		return fmt.Errorf("[config.validate]: auth.hmac_algorithm must be 'sha1' or 'sha256'")
	}

	return nil
}

func MustLoad(configPath string) *Config {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("[config.MustLoad]: failed to load config: %v", err))
	}
	return cfg
}
