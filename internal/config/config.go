package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	DB     DBConfig     `mapstructure:"db"`
	Secure SecureConfig `mapstructure:"secure"`
}

type AppConfig struct {
	Env      string `mapstructure:"env"`
	Port     int    `mapstructure:"port"`
	Debug    bool   `mapstructure:"debug"`
	LogLevel string `mapstructure:"log_level"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"` // sslmode는 disable, require, verify-ca, verify-full 를 설정 가능
}

type SecureConfig struct {
	CORSAllowOrigins []string `mapstructure:"cors_allow_origins"`
}

func LoadConfig(env string) (*Config, error) {

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath("./config")     // 실행 파일 기준 경로
	viper.AddConfigPath("../../config") // 테스트 환경에서 상대 경로로 접근 시 대비

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
