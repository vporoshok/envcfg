package envcfg_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/vporoshok/envcfg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleDefault() {
	type StdConfig struct {
		S string `default:"foo"`
		R string `default:"buz"`
	}
	type Config struct {
		StdConfig
		N int  `default:"42"`
		B bool `default:"true"`
	}

	cfg := Config{}

	envcfg.Default(&cfg, map[string]string{
		"R": "baz",
	})

	fmt.Println(cfg)
}

func TestDefault(t *testing.T) {
	type StdConfig struct {
		S string `default:"foo"`
		N int    `default:"42"`
		B bool   `default:"true"`
	}
	type Config struct {
		StdConfig
		D time.Duration `default:"1m"`
		F float64       `default:"36.6"`
		L []int         `default:"1,2,3,4"`
	}

	cfg := Config{}

	require.NoError(t, envcfg.Default(&cfg, map[string]string{"S": "bar"}))
	assert.Equal(t, "bar", cfg.S)
	assert.Equal(t, 42, cfg.N)
	assert.Equal(t, true, cfg.B)
	assert.Equal(t, time.Minute, cfg.D)
	assert.Equal(t, 36.6, cfg.F)
	assert.EqualValues(t, []int{1, 2, 3, 4}, cfg.L)
}

func ExampleRead() {
	type StdConfig struct {
		Debug   bool
		Default string `default:"foo"`
		HostIP  string
	}
	cfg := struct {
		StdConfig
		Number  uint32
		VH      string `envcfg:"VH_NAME"`
		Exclude string `envcfg:"-"`
	}{}

	err := envcfg.Read(&cfg, envcfg.WithDefault(map[string]string{"Debug": "true"}), envcfg.WithPrefix("PREFIX42_"))
	if err != nil {

		log.Fatal(err)
	}

	log.Print(cfg)
}

func TestRead(t *testing.T) {
	os.Setenv("PREFIX42_DEBUG", "true")
	os.Setenv("PREFIX42_NUMBER", "42")
	os.Setenv("HOST_IP", "1.1.1.1")
	os.Setenv("PREFIX42_VH_NAME", "someVH VH")

	type StdConfig struct {
		Debug   bool
		Default string `default:"foo"`
		HostIP  string
	}
	cfg := struct {
		StdConfig
		Number  uint32
		VH      string `envcfg:"VH_NAME"`
		Exclude string `envcfg:"-"`
		private string
	}{
		Exclude: "exclude",
		private: "bar",
	}

	require.NoError(t, envcfg.Read(&cfg,
		envcfg.WithDefault(map[string]string{"Exclude": "default"}),
		envcfg.WithPrefix("PREFIX42_"),
	))
	assert.Equal(t, true, cfg.Debug)
	assert.Equal(t, "foo", cfg.Default)
	assert.EqualValues(t, 42, cfg.Number)
	assert.Equal(t, "1.1.1.1", cfg.HostIP)
	assert.Equal(t, "someVH VH", cfg.VH)
	assert.Equal(t, "default", cfg.Exclude)
	assert.Equal(t, "bar", cfg.private)
}

func TestSplitWords(t *testing.T) {
	cases := []struct {
		name string
		src  string
		res  []string
	}{
		{
			name: "camelCase",
			src:  "someCamelCase",
			res:  []string{"some", "Camel", "Case"},
		},
		{
			name: "snake_case",
			src:  "some_snake_case",
			res:  []string{"some", "snake", "case"},
		},
		{
			name: "abbr",
			src:  "JSONFileAndSOMEMore",
			res:  []string{"JSON", "File", "And", "SOME", "More"},
		},
		{
			name: "numbers",
			src:  "JSON42File42An42dSO42MEMore",
			res:  []string{"JSON42", "File42", "An42d", "SO42ME", "More"},
		},
		{
			name: "mix",
			src:  "JSON42_File42An42dSO42ME_More",
			res:  []string{"JSON42", "File42", "An42d", "SO42ME", "More"},
		},
	}

	// nolint:scopelint
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := envcfg.SplitWords(c.src)
			assert.EqualValues(t, c.res, res)
		})
	}
}
