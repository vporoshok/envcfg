package envcfg_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vporoshok/envcfg"
)

func TestRead(t *testing.T) {
	os.Setenv("PREFIX42_DEBUG", "true")
	os.Setenv("PREFIX42_NUMBER", "42")
	os.Setenv("HOST_IP", "1.1.1.1")
	os.Setenv("PREFIX42_VH_NAME", "someVH")

	cfg := struct {
		Debug   bool
		Default string `default:"foo"`
		Number  uint32
		HostIP  string
		VH      string `envcfg:"VH_NAME"`
		Exclude string `envcfg:"-"`
		private string
	}{
		Exclude: "exclude",
		private: "bar",
	}

	require.NoError(t, envcfg.Read(&cfg, envcfg.WithDefault(), envcfg.WithPrefix("PREFIX42")))
	assert.Equal(t, true, cfg.Debug)
	assert.Equal(t, "foo", cfg.Default)
	assert.EqualValues(t, 42, cfg.Number)
	assert.Equal(t, "1.1.1.1", cfg.HostIP)
	assert.Equal(t, "someVH", cfg.VH)
	assert.Equal(t, "exclude", cfg.Exclude)
	assert.Equal(t, "bar", cfg.private)
}
