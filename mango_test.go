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

	PurchaserOnly int `access:"purchaserOnly"`
}

const PURCHASER_ONLY = iota

var conf MangoConfig = DefaultConfig().Extend(map[int]string{
	PURCHASER_ONLY: "purchaserOnly",
}).EnableLogs()

var have mockModel = mockModel{
	ID:       "foo",
	Username: "bar",
	Password: "baz",

	CreateOnly:    1,
	UpdateOnly:    1,
	PubOnly:       1,
	PrivOnly:      1,
	PurchaserOnly: 1,
}

func mockEqual(a mockModel, b mockModel) bool {
	if a.ID == b.ID &&
		a.Username == b.Username &&
		a.Password == b.Password &&
		a.CreateOnly == b.CreateOnly &&
		a.UpdateOnly == b.UpdateOnly &&
		a.PubOnly == b.PubOnly &&
		a.PrivOnly == b.PrivOnly &&
		a.PurchaserOnly == b.PurchaserOnly {
		return true
	}

	return false
}

func TestTrimCreate(t *testing.T) {
	result := Trim(have, CREATE, conf)

	want := mockModel{
		ID:       "",
		Username: "bar",
		Password: "baz",

		CreateOnly:    1,
		UpdateOnly:    0,
		PubOnly:       0,
		PrivOnly:      0,
		PurchaserOnly: 0,
	}

	if mockEqual(result, want) {
		t.Log("success")
		return
	}

	t.Error("mock values are not equal")
}

func TestTrimUpdate(t *testing.T) {
	result := Trim(have, UPDATE, conf)

	want := mockModel{
		ID:       "",
		Username: "bar",
		Password: "",

		CreateOnly:    0,
		UpdateOnly:    1,
		PubOnly:       0,
		PrivOnly:      0,
		PurchaserOnly: 0,
	}

	if mockEqual(result, want) {
		t.Log("success")
		return
	}

	t.Error("mock values are not equal")
}

func TestTrimPub(t *testing.T) {
	result := Trim(have, PUB, conf)

	want := mockModel{
		ID:       "",
		Username: "bar",
		Password: "",

		CreateOnly:    0,
		UpdateOnly:    0,
		PubOnly:       1,
		PrivOnly:      0,
		PurchaserOnly: 0,
	}

	if mockEqual(result, want) {
		t.Log("success")
		return
	}

	t.Error("mock values are not equal")
}

func TestTrimPriv(t *testing.T) {
	result := Trim(have, PRIV, conf)

	want := mockModel{
		ID:       "",
		Username: "bar",
		Password: "baz",

		CreateOnly:    0,
		UpdateOnly:    0,
		PubOnly:       0,
		PrivOnly:      1,
		PurchaserOnly: 0,
	}

	if mockEqual(result, want) {
		t.Log("success")
		return
	}

	t.Error("mock values are not equal")
}

func TestTrimCustom(t *testing.T) {
	result := Trim(have, PURCHASER_ONLY, conf)

	want := mockModel{
		ID:       "",
		Username: "",
		Password: "",

		CreateOnly:    0,
		UpdateOnly:    0,
		PubOnly:       0,
		PrivOnly:      0,
		PurchaserOnly: 1,
	}

	if mockEqual(result, want) {
		t.Log("success")
		return
	}

	t.Error("mock values are not equal")
}
