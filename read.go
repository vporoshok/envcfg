package envcfg

import (
	"os"
	"reflect"
	"strings"
)

func Read(v interface{}, opts ...Option) error {
	cfg := readConfig{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}

	s := reflect.ValueOf(v)
	if err := checkType(s); err != nil {

		return err
	}

	if cfg.useDefault {
		if err := applyDefault(s); err != nil {

			return err
		}
	}

	m, err := getTagMap(s.Elem().Type(), "envcfg")
	if err != nil {

		return err
	}

	for k, v := range m {
		if v == "-" {
			delete(m, k)

			continue
		}
		if len(v) == 0 {
			v = strings.ToUpper(strings.Join(SplitWords(k), "_"))
		}
		m[k] = os.Getenv(cfg.prefix + v)
		if len(m[k]) == 0 && len(cfg.prefix) > 0 {
			m[k] = os.Getenv(v)
		}
		if len(m[k]) == 0 {
			delete(m, k)

			continue
		}
	}

	return applyMap(s, m)
}
