package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driif/golang-test-task/internal/server/config/env"
	"github.com/stretchr/testify/assert"
)

var (
	pwd, _ = os.Getwd()
)

func TestDotEnvOverride(t *testing.T) {
	assert.Equal(t, "", os.Getenv("IS_THIS_A_TEST_ENV"))

	orgPsqlUser := os.Getenv("PSQL_USER")

	env.DotEnvTryLoad(
		filepath.Join(pwd, "/testdata/.env1.local"),
		func(k string, v string) error { t.Setenv(k, v); return nil })

	assert.Equal(t, "yes", os.Getenv("IS_THIS_A_TEST_ENV"))
	assert.Equal(t, "dotenv_override_psql_user", os.Getenv("PSQL_USER"))
	assert.Equal(t, orgPsqlUser, os.Getenv("ORIGINAL_PSQL_USER"))

	// override works as expected?
	env.DotEnvTryLoad(
		filepath.Join(pwd, "/testdata/.env2.local"),
		func(k string, v string) error { t.Setenv(k, v); return nil })

	assert.Equal(t, "yes still", os.Getenv("IS_THIS_A_TEST_ENV"))
	assert.NotEqual(t, "dotenv_override_psql_user", os.Getenv("PSQL_USER"))
	assert.Equal(t, orgPsqlUser, os.Getenv("PSQL_USER"), "Reset to original does not work!")
}

func TestNoopEnvNotFound(t *testing.T) {
	assert.NotPanics(t, assert.PanicTestFunc(func() {

		env.DotEnvTryLoad(
			filepath.Join(pwd, "/testdata/.env.does.not.exist"),
			func(k string, v string) error { t.Setenv(k, v); return nil })

	}), "does not panic on file inexistance")
}

func TestEmptyEnv(t *testing.T) {
	assert.NotPanics(t, assert.PanicTestFunc(func() {

		env.DotEnvTryLoad(
			filepath.Join(pwd, "/testdata/.env.local.sample"),
			func(k string, v string) error { t.Setenv(k, v); return nil })

	}), "does not panic on file inexistance")

	assert.Empty(t, os.Getenv("EMPTY_VARIABLE_INIT"), "should be empty")
}

func TestPanicsOnEnvMalform(t *testing.T) {
	assert.Panics(t, assert.PanicTestFunc(func() {

		env.DotEnvTryLoad(
			filepath.Join(pwd, "/testdata/.env.local.malformed"),
			func(k string, v string) error { t.Setenv(k, v); return nil })

	}), "does panic on file malform")
}
