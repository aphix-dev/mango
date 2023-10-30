package mango

import "testing"

type mockModel struct {
	ID       string
	Username string `access:"pub,priv,create,update"`
	Password string `access:"priv,create"`

	CreateOnly int `access:"create"`
	UpdateOnly int `access:"update"`
	PubOnly    int `access:"pub"`
	PrivOnly   int `access:"priv"`
}

var have mockModel = mockModel{
	ID:       "foo",
	Username: "bar",
	Password: "baz",
}

func TestTrimCreate(t *testing.T) {
	Trim(have, CREATE)
}

func TestTrimUpdate(t *testing.T) {
	Trim(have, UPDATE)
}

func TestTrimPub(t *testing.T) {
	Trim(have, PUB)
}

func TestTrimPriv(t *testing.T) {
	Trim(have, PRIV)
}
