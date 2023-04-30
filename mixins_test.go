package main

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func Test_Env(t *testing.T) {
	home, _ := os.UserHomeDir()
	var tests = []struct {
		Prefix, Key string
	}{
		{"%env=", "HOME"},
		{"%custom=", "HOME"},
	}
	for _, testCase := range tests {
		name := fmt.Sprintf("Testing_with_prefix_ prefix_%s", testCase.Prefix)
		t.Setenv("MIXIN_PREFIX_ENV_VAR", testCase.Prefix)
		t.Run(name, func(t *testing.T) {
			var mix Mixin = Mixin(fmt.Sprintf("%s%s", testCase.Prefix, testCase.Key))
			val, err := mix.Read(context.Background())
			if err != nil {
				t.Errorf("Unable to read mixin '%v' : %v", mix, err)
			} else if val != home {
				t.Errorf("Expected: '%s', Got: '%s'", home, val)
			}
		})
	}
}

func Test_Secret(t *testing.T) {
	var contents = "contents"
	_ = os.MkdirAll("tests/secret", 0700)
	_ = os.WriteFile("tests/secret/key", []byte(contents), 0600)
	var tests = []struct {
		Prefix, Key string
	}{
		{"%secret=", "secret/key"},
		{"%custom=", "secret/key"},
	}
	t.Setenv("SECRETS_ROOT_DIR", "./tests")
	t.Setenv("MIXIN_PREFIX_ENV_VAR", "%env")
	for _, testCase := range tests {
		name := fmt.Sprintf("Testing_with_prefix_%s", testCase.Prefix)
		_ = os.Setenv("MIXIN_PREFIX_SECRET", testCase.Prefix)
		t.Run(name, func(t *testing.T) {
			var mix Mixin = Mixin(fmt.Sprintf("%s%s", testCase.Prefix, testCase.Key))
			val, err := mix.Read(context.Background())
			if err != nil {
				t.Errorf("Unable to read mixin '%v' : %v", mix, err)
			} else if val != contents {
				t.Errorf("Expected: '%s', Got: '%s'", contents, val)
			}
		})
	}
}

func Test_MissingEnvKey(t *testing.T) {
	var mix Mixin = "%env=NON_EXISTENT_ENV_VAR"
	_, err := mix.Read(context.Background())
	if err == nil {
		t.Error("Should have had an error here")
	} else if err.Error() != "mixins: environment variable named 'NON_EXISTENT_ENV_VAR' does not exists" {
		t.Error("Unexpected error")
	}
}

func Test_NonExistentSecret(t *testing.T) {
	var mix Mixin = "%secret=path_do_not/exists"
	t.Setenv("SECRETS_ROOT_DIR", "./tests")
	t.Setenv("MIXIN_PREFIX_SECRET", "%secret=")
	_, err := mix.Read(context.Background())
	if err == nil {
		t.Error("Should have had an error here")
	} else if err.Error() != "mixins: unable to read secret 'path_do_not/exists'" {
		t.Error(fmt.Sprintf("Unexpected error: %s", err))
	}
}

func Test_UndefinedMixin(t *testing.T) {
	var mix Mixin = "%undefined=KEY"
	_, err := mix.Read(context.Background())
	if err == nil {
		t.Error("Should have had an error here")
	} else if err.Error() != "unknown mixin type" {
		t.Error("Unexpected error")
	}
}
