package envcfg_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
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
		S string        `default:"foo"`
		N int           `default:"42"`
		B bool          `default:"true"`
		D time.Duration `default:"1m"`
		F float64       `default:"36.6"`
		L []int         `default:"1,2,3,4"`
	}

	cfg := Config{}

	require.NoError(t, envcfg.Default(&cfg))
	assert.Equal(t, "foo", cfg.S)
	assert.Equal(t, 42, cfg.N)
	assert.Equal(t, true, cfg.B)
	assert.Equal(t, time.Minute, cfg.D)
	assert.Equal(t, 36.6, cfg.F)
	assert.EqualValues(t, []int{1, 2, 3, 4}, cfg.L)
}

func TestDefaultInvalidType(t *testing.T) {
	require.EqualError(t, envcfg.Default(struct{}{}), envcfg.InvalidObjectType.Error())
}

func ExampleRead() {
	cfg := struct {
		Debug   bool
		Default string `default:"foo"`
		Number  uint32
		HostIP  string
		VH      string `envcfg:"VH_NAME"`
		Exclude string `envcfg:"-"`
	}{}

	err := envcfg.Read(&cfg, envcfg.WithDefault(), envcfg.WithPrefix("PREFIX42_"))
	if err != nil {

		log.Fatal(err)
	}

	log.Print(cfg)
}

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

	require.NoError(t, envcfg.Read(&cfg, envcfg.WithDefault(), envcfg.WithPrefix("PREFIX42_")))
	assert.Equal(t, true, cfg.Debug)
	assert.Equal(t, "foo", cfg.Default)
	assert.EqualValues(t, 42, cfg.Number)
	assert.Equal(t, "1.1.1.1", cfg.HostIP)
	assert.Equal(t, "someVH", cfg.VH)
	assert.Equal(t, "exclude", cfg.Exclude)
	assert.Equal(t, "bar", cfg.private)
}

func TestReadInvalidObjectType(t *testing.T) {
	require.EqualError(t, envcfg.Read(struct{}{}), envcfg.InvalidObjectType.Error())
}

func TestReadInvalidFieldType(t *testing.T) {
	os.Setenv("PREFIX42_TEST", "true")
	err := envcfg.Read(&struct{ Test interface{} }{}, envcfg.WithPrefix("PREFIX42_"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "interface")
	assert.Equal(t, envcfg.InvalidFieldType, errors.Cause(err))
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

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := envcfg.SplitWords(c.src)
			assert.EqualValues(t, c.res, res)
		})
	}
}
