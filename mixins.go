package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v8"
	"os"
	"path/filepath"
	"strings"
)

type Engine struct {
	SecretsRootDir string `env:"SECRETS_ROOT_DIR"     envDefault:"/run/secrets"`
	PrefixSecret   string `env:"MIXIN_PREFIX_SECRET"  envDefault:"%secret="`
	PrefixEnv      string `env:"MIXIN_PREFIX_ENV_VAR" envDefault:"%env="`
}

type Mixin string

var _engine *Engine

func getEngine() *Engine {
	_engine = new(Engine)
	_ = env.Parse(_engine)
	return _engine
}

func getKey(reference, prefix string) string { return strings.Replace(reference, prefix, "", 1) }

func (en *Engine) getEnvValue(ctx context.Context, reference string) (string, error) {
	key := getKey(reference, en.PrefixEnv)
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", errors.New(fmt.Sprintf("mixins: environment variable named '%s' does not exists", key))
	}
	return value, nil
}

func (en *Engine) getSecretValue(ctx context.Context, reference string) (string, error) {
	key := getKey(reference, en.PrefixSecret)
	secretFile := filepath.Join(en.SecretsRootDir, key)
	value, err := os.ReadFile(secretFile)
	if err != nil {
		return "", errors.New(fmt.Sprintf("mixins: unable to read secret '%s'", secretFile))
	}
	return string(value), nil
}

func (m Mixin) Read(ctx context.Context) (string, error) {
	en := getEngine()
	if strings.HasPrefix(string(m), en.PrefixEnv) {
		return en.getEnvValue(ctx, string(m))
	}
	if strings.HasPrefix(string(m), en.PrefixSecret) {
		return en.getSecretValue(ctx, string(m))
	}
	return "", errors.New("unknown mixin type")
}
