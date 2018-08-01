package envcfg_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vporoshok/envcfg"
)

func ExampleDefault() {
	type Config struct {
		S string `default:"foo"`
		N int    `default:"42"`
		B bool   `default:"true"`
	}

	cfg := Config{}

	envcfg.Default(&cfg)

	fmt.Println(cfg)
}

func TestDefault(t *testing.T) {
	type Config struct {
		S string `default:"foo"`
		N int    `default:"42"`
		B bool   `default:"true"`
	}

	cfg := Config{}

	require.NoError(t, envcfg.Default(&cfg))
	assert.Equal(t, "foo", cfg.S)
	assert.Equal(t, 42, cfg.N)
	assert.Equal(t, true, cfg.B)
}
