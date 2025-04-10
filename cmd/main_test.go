package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestMain(t *testing.T) {
	os.Setenv("BEETLE_ENV", "test") // Forces load of configs/test.yaml

	go main()
	is := is.New(t)

	_, err := http.NewRequest("GET", "localhost:8080", nil)
	is.NoErr(err)
}
