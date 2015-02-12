package tempconsul

import (
	"github.com/hashicorp/consul/api"
	"testing"
)

func TestServer(t *testing.T) {
	server := Server{}
	err := server.Start()
	defer server.Term()
	if err != nil {
		t.Fatal(err)
	}

	client, _ := api.NewClient(api.DefaultConfig())
	kv := client.KV()

	putpair := &api.KVPair{Key: "test", Value: []byte("cool")}
	_, err = kv.Put(putpair, nil)
	if err != nil {
		t.Fatal(err)
	}
	getpair, _, err := kv.Get("test", nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(getpair.Value) != "cool" {
		t.Errorf("got: %#v", string(getpair.Value))
	}
}
