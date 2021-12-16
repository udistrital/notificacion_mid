package helpers

import (
	"flag"
	"os"
	"testing"
)

var parameters struct {
	Endpoint1 string
}

func TestMain(m *testing.M) {
	parameters.Endpoint1 = os.Getenv("Endpoint1")
	flag.Parse()
	os.Exit(m.Run())
}
