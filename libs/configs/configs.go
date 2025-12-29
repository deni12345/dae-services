package configs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

// Source abstracts config loading from different sources
type Source interface {
	Apply(interface{}) error
}

type yamlSource struct {
	path string
}

func (y *yamlSource) Apply(cfg interface{}) error {
	abs, err := filepath.Abs(y.path)
	if err != nil {
		return fmt.Errorf("resolve config path %q: %w", y.path, err)
	}
	data, err := os.ReadFile(abs)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			slog.Debug("yaml config not found, skipping", "path", y.path, "error", err)
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

func (e *envSource) Apply(cfg interface{}) error {
	return cleanenv.ReadEnv(cfg)
}

// WithYamlFile creates a YAML file config source
func WithYamlFile(path string) Source {
	if path == "" {
		return nil
	}
	return &yamlSource{
		path: path,
	}
}

// WithEnv creates an environment variable config source
func WithEnv() Source {
	return &envSource{}
}

// Load applies multiple config sources to the provided config struct
func Load(cfg interface{}, sources ...Source) error {
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

// LoadWithRetry loads config with retry logic and exponential backoff
func LoadWithRetry(ctx context.Context, cfg interface{}, sources ...Source) error {
	var err error
	// Build sources here so we don't mutate the caller's slice
	srcs := append([]Source{WithEnv()}, sources...)
	retryDelay := delay

	for i := 0; i < defaultRetries; i++ {
		if err = Load(cfg, srcs...); err == nil {
			return nil
		}

		if i == defaultRetries-1 {
			return fmt.Errorf("load config failed after %d attempts: %w", i+1, err)
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("load config canceled: %w", ctx.Err())
		case <-time.After(retryDelay):
		}

		// exponential backoff
		retryDelay *= 2
		if retryDelay > maxDelay {
			retryDelay = maxDelay
		}
	}
	return nil
}

func LoadWithEnvOptions(ctx context.Context, cfg interface{}, envVar string, sources ...Source) error {
	if envVar == "" {
		envVar = "ENVIRONMENT"
	}

	envVal := os.Getenv(envVar)
	production := []string{"production", "prod", "prd"}
	for _, v := range production {
		if strings.EqualFold(envVal, v) {
			return LoadWithRetry(ctx, cfg)
		}
	}
	return LoadWithRetry(ctx, cfg, sources...)
}
