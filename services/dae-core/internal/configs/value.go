package configs

// Value holds dae-core service configuration
type Value struct {
	Environment        string `yaml:"environment" env:"ENVIRONMENT" env-default:"dev"`
	LogLevel           string `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
	ServiceName        string `yaml:"service_name" env:"SERVICE_NAME" env-default:"dae-core-service"`
	GRPCAddress        string `yaml:"grpc_address" env:"GRPC_ADDRESS" env-default:"50051"`
	FirestoreProjectID string `yaml:"firestore_project_id" env:"FIRESTORE_PROJECT_ID" env-default:""`
	RedisAddr          string `yaml:"redis_addr" env:"REDIS_ADDR" env-default:"localhost:6379"`
	RedisPassword      string `yaml:"redis_password" env:"REDIS_PASSWORD" env-default:"password"`
	PageSize           int32  `yaml:"page_size" env:"PAGE_SIZE" env-default:"20"`
	OtelCol            string `yaml:"otelcol" env:"OTELCOL" env-default:"tempo:4317"`
	Insecure           bool   `yaml:"insecure" env:"INSECURE" env-default:"true"`

	// Observability toggles
	EnableTracing bool `yaml:"enable_tracing" env:"ENABLE_TRACING" env-default:"true"`
	EnableMetrics bool `yaml:"enable_metrics" env:"ENABLE_METRICS" env-default:"true"`
	EnableLogging bool `yaml:"enable_logging" env:"ENABLE_LOGGING" env-default:"true"`
}
