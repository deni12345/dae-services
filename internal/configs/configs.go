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

type Config struct {
	Port      int       `yaml:"port" env:"PORT" env-default:"8080"`
	Firestore Firestore `yaml:"fire_store" env-prefix:"FIRESTORE_"`
}

type Source interface {
	Apply(*Config) error
}

type yamlSource struct {
	path string
}

func (y *yamlSource) Apply(cfg *Config) error {
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

func (e *envSource) Apply(cfg *Config) error {
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

func Load(cfg *Config, sources ...Source) error {
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
	Values         = Config{}
)

func LoadWithRetry(ctx context.Context, sources ...Source) error {
	var err error
	var cfg Config

	for i := 0; i < defaultRetries; i++ {
		sources = append(sources, WithEnv())
		err = Load(&cfg, sources...)

		if err == nil {
			Values = cfg
			return nil
		}

		if i == defaultRetries-1 {
			return fmt.Errorf("load config with entry after %d attempts: %w after ", i+1, err)
		}

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("load config with entry canceled: %w after ", err)
		case <-timer.C:
		}

		delay *= time.Duration(i)
		if delay > maxDelay {
			delay = maxDelay
		}
	}
	return err
}
