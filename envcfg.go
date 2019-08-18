package envcfg

import (
	"os"
	"strings"

	"github.com/vporoshok/casey"
	"github.com/vporoshok/reflector"
)

type readConfig struct {
	prefix     string
	useDefault bool
	overrides  map[string]string
}

// Option to configure environment read
type Option interface {
	apply(*readConfig)
}

type funcOption func(*readConfig)

func (fo funcOption) apply(cfg *readConfig) {
	fo(cfg)
}

// WithPrefix add prefix to read variable
//
// If PREFIX_VARNAME is not setted, try to read VARNAME without prefix
func WithPrefix(prefix string) Option {

	return funcOption(func(cfg *readConfig) {
		cfg.prefix = prefix
	})
}

// WithDefault add read default values from tags before read environment
func WithDefault(overrides map[string]string) Option {

	return funcOption(func(cfg *readConfig) {
		cfg.useDefault = true
		cfg.overrides = overrides
	})
}

// Read environment into struct
func Read(v interface{}, opts ...Option) error {
	cfg := readConfig{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	r := reflector.New(v)
	if cfg.useDefault {
		if err := applyDefault(r, cfg.overrides); err != nil {

			return err
		}
	}

	m := r.ExtractTags("envcfg", reflector.WithoutMinus())
	for k, v := range m {
		if len(v) == 0 {
			v = strings.ToUpper(casey.Camel(k).SNAKE())
		}
		var ok bool
		m[k], ok = os.LookupEnv(cfg.prefix + v)
		if !ok && len(cfg.prefix) > 0 {
			m[k], ok = os.LookupEnv(v)
		}
		if !ok {
			delete(m, k)

			continue
		}
	}

	return r.Apply(m)
}

// Default read and parse values from struct tag `default:"some value"`
func Default(v interface{}, overrides map[string]string) error {
	r := reflector.New(v)

	return applyDefault(r, overrides)
}

func applyDefault(r reflector.Reflector, overrides map[string]string) error {
	m := r.ExtractTags("default", reflector.WithoutEmpty())
	for k, v := range overrides {
		m[k] = v
	}

	return r.Apply(m)
}
