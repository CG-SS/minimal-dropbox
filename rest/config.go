package rest

type System string

const (
	Gin System = "gin"
	Nop System = "nop"
)

type CorsConfig struct {
	AllowOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS" default:"http://localhost:3000"`
}

type Config struct {
	System System `envconfig:"REST_SYSTEM" default:"gin"`
	Host   string `envconfig:"REST_HOST" default:"127.0.0.1"`
	Port   int    `envconfig:"REST_PORT" default:"12345"`
	Cors   CorsConfig
}
