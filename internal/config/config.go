package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `yaml:"app"`
	DB     DBConfig     `yaml:"db"`
	Secure SecureConfig `yaml:"secure"`
}

type AppConfig struct {
	Env   string `yaml:"env"`
	Port  int    `yaml:"port"`
	Debug bool   `yaml:"debug"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"` // sslmode는 disable, require, verify-ca, verify-full 를 설정 가능
}

type SecureConfig struct {
	CORSAllowOrigins []string `yaml:"cors_allow_origins"`
	JWTSecret        string   `yaml:"jwt_secret"`
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
