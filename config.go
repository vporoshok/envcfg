package envcfg

type readConfig struct {
	prefix     string
	useDefault bool
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
		cfg.prefix = prefix + "_"
	})
}

// WithDefault add read default values from tags before read environment
func WithDefault() Option {

	return funcOption(func(cfg *readConfig) {
		cfg.useDefault = true
	})
}
