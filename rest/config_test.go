package rest

import (
	"fmt"
	"os"
	"testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
)

func TestEnvconfigExtractsAllConfigs(t *testing.T) {
	type envTest struct {
		key   string
		value string
	}

	corsEnabledVal := true
	corsOriginsVal := "http://test123"
	restSysVal := Nop
	restHostVal := "localhost"
	restPort := 44444
	homeEnabled := false

	var expectedCfg Config
	expectedCfg.Cors.Enabled = corsEnabledVal
	expectedCfg.Cors.AllowOrigins = []string{corsOriginsVal}
	expectedCfg.Host = restHostVal
	expectedCfg.Port = restPort
	expectedCfg.System = restSysVal
	expectedCfg.HomeRouteEnabled = homeEnabled
	expectedCfg.BufferSize = 1024

	envVarsTest := []envTest{
		{
			key:   "CORS_ENABLED",
			value: fmt.Sprintf("%t", corsEnabledVal),
		},
		{
			key:   "CORS_ALLOWED_ORIGINS",
			value: corsOriginsVal,
		},
		{
			key:   "REST_SYSTEM",
			value: string(restSysVal),
		},
		{
			key:   "REST_HOST",
			value: restHostVal,
		},
		{
			key:   "REST_PORT",
			value: fmt.Sprintf("%d", restPort),
		},
		{
			key:   "REST_HOME_ROUTE_ENABLED",
			value: fmt.Sprintf("%t", homeEnabled),
		},
	}

	for _, envVar := range envVarsTest {
		err := os.Setenv(envVar.key, envVar.value)
		assert.NoError(t, err)
	}

	var cfg Config
	err := envconfig.Process("", &cfg)
	assert.NoError(t, err)

	assert.Equal(t, expectedCfg, cfg)
}
