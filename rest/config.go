package rest

type System string

const (
	Gin System = "gin"
	Nop System = "nop"
)

type CorsConfig struct {
	Enabled      bool     `envconfig:"CORS_ENABLED" default:"true"`
	AllowOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS" default:"http://localhost:3000"`
}

type Config struct {
	HomeRouteEnabled bool   `envconfig:"REST_HOME_ROUTE_ENABLED" default:"true"`
	System           System `envconfig:"REST_SYSTEM" default:"gin"`
	Host             string `envconfig:"REST_HOST" default:"127.0.0.1"`
	Port             int    `envconfig:"REST_PORT" default:"12345"`
	Cors             CorsConfig
}
