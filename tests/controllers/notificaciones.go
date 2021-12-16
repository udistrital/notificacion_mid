package controllers

import (
	"flag"
	"os"
	"testing"
)

var parameters struct {
	Endpoint1 string
	Endpoint2 string
	Endpoint3 string
	Endpoint4 string
	Endpoint5 string
}

func TestMain(m *testing.M) {
	parameters.Endpoint1 = os.Getenv("Endpoint1")
	parameters.Endpoint2 = os.Getenv("Endpoint2")
	parameters.Endpoint3 = os.Getenv("Endpoint3")
	parameters.Endpoint4 = os.Getenv("Endpoint4")
	parameters.Endpoint5 = os.Getenv("Endpoint5")
	flag.Parse()
	os.Exit(m.Run())
}

func TestEndPointPostOneNotif(t *testing.T) {
	t.Log("AAAAAAAAAAAAAAAA")
	t.Log(parameters.Endpoint1)
	t.Log("AAAAAAAAAAAAAAAA")
}

func TestEndPointSubscribe(t *testing.T) {
	t.Log("AAAAAAAAAAAAAAAA")
	t.Log(parameters.Endpoint2)
	t.Log("AAAAAAAAAAAAAAAA")
}

func TestEndPointGetTopics(t *testing.T) {
	t.Log("AAAAAAAAAAAAAAAA")
	t.Log(parameters.Endpoint3)
	t.Log("AAAAAAAAAAAAAAAA")
}

func TestEndPointCreateTopic(t *testing.T) {
	t.Log("AAAAAAAAAAAAAAAA")
	t.Log(parameters.Endpoint4)
	t.Log("AAAAAAAAAAAAAAAA")
}

func TestEndPointVerifSus(t *testing.T) {
	t.Log("AAAAAAAAAAAAAAAA")
	t.Log(parameters.Endpoint5)
	t.Log("AAAAAAAAAAAAAAAA")
}
