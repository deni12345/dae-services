package configs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

type Firestore struct {
	SecretKey string `yaml:"secret_key" env:"SECRET_KEY" env-default:""`
	ProjectID string `yaml:"project_id" env:"PROJECT_ID" env-default:""`
}

type Value struct {
	Environment string    `yaml:"environment" env:"ENVIRONMENT" env-default:"dev"`
	LogLevel    string    `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
	ServiceName string    `yaml:"service_name" env:"SERVICE_NAME" env-default:"dae-core-service"`
	GRPCAddress string    `yaml:"grpc_address" env:"GRPC_ADDRESS" env-default:"50051"`
	Firestore   Firestore `yaml:"firestore" env-prefix:"FIRESTORE_"`
	Redis       struct {
		Addr     string `yaml:"addr" env:"ADDR" env-default:"localhost:6379"`
		Password string `yaml:"password" env:"PASSWORD" env-default:"password"`
	}
	PageSize int32  `yaml:"page_size" env:"PAGE_SIZE" env-default:"20"`
	OtelCol  string `yaml:"otelcol" env:"OTELCOL" env-default:"tempo:4317"`
	Insecure bool   `yaml:"insecure" env:"INSECURE" env-default:"true"`
}

type Source interface {
	Apply(*Value) error
}

type yamlSource struct {
	path string
}

func (y *yamlSource) Apply(cfg *Value) error {
	data, err := os.ReadFile(y.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read yaml %s: %w", y.path, err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}
	return nil
}

type envSource struct{}

func (e *envSource) Apply(cfg *Value) error {
	return cleanenv.ReadEnv(cfg)
}

func WithYamlFile(path string) Source {
	if path == "" {
		return nil
	}
	return &yamlSource{
		path: path,
	}
}

func WithEnv() Source {
	return &envSource{}
}

func Load(cfg *Value, sources ...Source) error {
	var errs error

	for _, s := range sources {
		if s == nil {
			continue
		}
		if err := s.Apply(cfg); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

var (
	defaultRetries = 3
	delay          = 100 * time.Millisecond
	maxDelay       = 1000 * time.Millisecond
)

func LoadWithRetry(ctx context.Context, sources ...Source) (Value, error) {
	var err error
	var cfg Value

	for i := 0; i < defaultRetries; i++ {
		sources = append(sources, WithEnv())
		err = Load(&cfg, sources...)

		if err == nil {
			return cfg, nil
		}

		if i == defaultRetries-1 {
			return cfg, fmt.Errorf("load config with entry after %d attempts: %w after ", i+1, err)
		}

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return cfg, fmt.Errorf("load config with entry canceled: %w after ", err)
		case <-timer.C:
		}

		delay *= time.Duration(i)
		if delay > maxDelay {
			delay = maxDelay
		}
	}
	return cfg, err
}
