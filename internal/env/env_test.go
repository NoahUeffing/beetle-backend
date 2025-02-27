package env_test

import (
	"beetle/internal/env"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestGet(t *testing.T) {
	is := is.New(t)

	k := "TEST"
	os.Setenv(k, "fizz")
	is.Equal(env.Get(k, "buzz"), "fizz")
	os.Unsetenv(k)
}

func TestGet_Wrapped(t *testing.T) {
	is := is.New(t)

	k := "TEST"
	os.Setenv(k, "\"fizz\"")
	is.Equal(env.Get(k, "buzz"), "fizz")
	os.Unsetenv(k)
}

func TestGet_Fallback(t *testing.T) {
	is := is.New(t)

	k := "TEST"
	is.Equal(env.Get(k, "buzz"), "buzz")
}

func TestGetBool(t *testing.T) {
	is := is.New(t)

	k := "TEST"
	os.Setenv(k, "true")
	is.True(env.GetBool(k, false))
	os.Unsetenv(k)
}

func TestGetBool_Wrapped(t *testing.T) {
	is := is.New(t)

	k := "TEST"
	os.Setenv(k, "\"true\"")
	is.True(env.GetBool(k, false))
	os.Unsetenv(k)
}

func TestGetBool_Fallback(t *testing.T) {
	is := is.New(t)

	k := "TEST"
	is.Equal(env.GetBool(k, true), true)
}
