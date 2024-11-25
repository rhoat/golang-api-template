package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	gotel "github.com/rhoat/go-exercise/pkg/otel"
	"github.com/rhoat/go-exercise/pkg/system"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

func (cfg Config) LogLevel() zapcore.Level {
	switch strings.TrimSpace(strings.ToLower(cfg.General.LogLevel)) {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

type ServerConfig struct {
	Port              string `mapstructure:"port" yaml:"port"`
	ReadTimeout       int    `mapstructure:"ReadTimeout" yaml:"ReadTimeout"`
	ReadHeaderTimeout int    `mapstructure:"ReadHeaderTimeout" yaml:"ReadHeaderTimeout"`
	WriteTimeout      int    `mapstructure:"WriteTimeout" yaml:"WriteTimeout"`
	IdleTimeout       int    `mapstructure:"IdleTimeout" yaml:"IdleTimeout"`
}

type GeneralConfig struct {
	LogLevel string `mapstructure:"logLevel" yaml:"logLevel"`
}

type OtelConfig struct {
	Destination gotel.Destination `mapstructure:"destination" yaml:"destination"`
}

type Config struct {
	ServerConfig ServerConfig  `mapstructure:"serverConfig" yaml:"serverConfig"`
	General      GeneralConfig `mapstructure:"general" yaml:"general"`
	Otel         OtelConfig    `mapstructure:"otel" yaml:"otel"`
}

func LoadConfig() (*Config, error) {
	viper.SetEnvPrefix(system.ApplicationName)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, err
		}
	}
	viper.AutomaticEnv()
	var config Config
	if err := viper.Unmarshal(&config, func(dc *mapstructure.DecoderConfig) {
		dc.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			otelDestinationDecodeHook, // Our custom decode hook
		)
	}); err != nil {
		return nil, err
	}
	return config.Validate()
}

func (cfg Config) Validate() (*Config, error) {
	if cfg.ServerConfig.Port == "" {
		return &cfg, fmt.Errorf("port given is not valid port:%s", cfg.ServerConfig.Port)
	}
	return &cfg, nil
}
