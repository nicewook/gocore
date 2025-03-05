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
	CORSAllowOrigins []string  `mapstructure:"cors_allow_origins"`
	JWT              JWTConfig `mapstructure:"jwt"`
}

type JWTConfig struct {
	PrivateKey           string       `mapstructure:"private_key"`            // RSA 개인키 (서명용)
	PublicKey            string       `mapstructure:"public_key"`             // RSA 공개키 (검증용)
	AccessExpirationMin  int          `mapstructure:"access_expiration_min"`  // Access Token 만료 시간 (분)
	RefreshExpirationDay int          `mapstructure:"refresh_expiration_day"` // Refresh Token 만료 시간 (일)
	Cookie               CookieConfig `mapstructure:"cookie"`                 // 쿠키 관련 설정
}

type CookieConfig struct {
	Secure   bool   `mapstructure:"secure"`    // HTTPS 전용 여부
	HTTPOnly bool   `mapstructure:"http_only"` // JavaScript 접근 불가 여부
	SameSite string `mapstructure:"same_site"` // SameSite 정책 (Strict, Lax, None)
	Domain   string `mapstructure:"domain"`    // 쿠키 도메인
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
