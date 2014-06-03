package tempconsul

import (
	"github.com/armon/consul-kv"
	"testing"
)

func TestServer(t *testing.T) {
	server := Server{}
	err := server.Start()
	defer server.Term()
	if err != nil {
		t.Fatal(err)
	}

	client, _ := consulkv.NewClient(consulkv.DefaultConfig())
	err = client.Put("test", []byte("cool"), 0)
	if err != nil {
		t.Fatal(err)
	}
	_, pair, err := client.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	if string(pair.Value) != "cool" {
		t.Errorf("got: %#v", string(pair.Value))
	}
}
