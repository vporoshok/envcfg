package envcfg

type readConfig struct {
	prefix     string
	useDefault bool
}

type Option interface {
	apply(*readConfig)
}

type funcOption func(*readConfig)

func (fo funcOption) apply(cfg *readConfig) {
	fo(cfg)
}

func WithPrefix(prefix string) Option {

	return funcOption(func(cfg *readConfig) {
		cfg.prefix = prefix + "_"
	})
}

func WithDefault() Option {

	return funcOption(func(cfg *readConfig) {
		cfg.useDefault = true
	})
}
