package config

import (
	"os"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	// 设置环境变量
	os.Setenv("MYSQL_HOST", "localhost")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_USER", "root")
	os.Setenv("MYSQL_PASSWORD", "root1234")
	os.Setenv("MYSQL_DBNAME", "course")
	os.Setenv("PASSWORD_SALT", "test-salt")
	os.Setenv("JWT_SIGNING_KEY", "test-key")
	os.Setenv("HTTP_PORT", "8000")
	os.Setenv("ACCESS_TOKEN_TTL", "2h")
	os.Setenv("REFRESH_TOKEN_TTL", "720h")
	os.Setenv("CACHE_TTL", "60s")
	os.Setenv("SMTP_HOST", "mail.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_FROM", "test@example.com")
	os.Setenv("APP_ENV", "local")

	defer func() {
		os.Unsetenv("MYSQL_HOST")
		os.Unsetenv("MYSQL_PORT")
		os.Unsetenv("MYSQL_USER")
		os.Unsetenv("MYSQL_PASSWORD")
		os.Unsetenv("MYSQL_DBNAME")
		os.Unsetenv("PASSWORD_SALT")
		os.Unsetenv("JWT_SIGNING_KEY")
		os.Unsetenv("HTTP_PORT")
		os.Unsetenv("ACCESS_TOKEN_TTL")
		os.Unsetenv("REFRESH_TOKEN_TTL")
		os.Unsetenv("CACHE_TTL")
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("SMTP_FROM")
		os.Unsetenv("APP_ENV")
	}()

	cfg, err := Init(".")
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if cfg.Environment != "local" {
		t.Errorf("Environment = %v, want local", cfg.Environment)
	}

	if cfg.HTTP.Port != "8000" {
		t.Errorf("HTTP.Port = %v, want 8000", cfg.HTTP.Port)
	}

	if cfg.MySQL.Host != "localhost" {
		t.Errorf("MySQL.Host = %v, want localhost", cfg.MySQL.Host)
	}

	if cfg.Auth.PasswordSalt != "test-salt" {
		t.Errorf("Auth.PasswordSalt = %v, want test-salt", cfg.Auth.PasswordSalt)
	}

	if cfg.Auth.JWT.SigningKey != "test-key" {
		t.Errorf("Auth.JWT.SigningKey = %v, want test-key", cfg.Auth.JWT.SigningKey)
	}

	if cfg.CacheTTL != 60*time.Second {
		t.Errorf("CacheTTL = %v, want 60s", cfg.CacheTTL)
	}
}
