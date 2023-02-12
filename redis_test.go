package asynm

import (
	"net"
	"testing"

	"github.com/redis/go-redis/v9"
)

func createExampleRedisClient() (*redisClient, error) {
	return newRedisClient(&redis.UniversalOptions{
		Addrs:    []string{net.JoinHostPort("127.0.0.1", "6379")},
		Password: "",
		PoolSize: 10,
	})
}

func Test_redisClient_HGet(t *testing.T) {
	r, err := createExampleRedisClient()
	if err != nil {
		t.Logf("createExampleRedisClient error: %s", err.Error())
		return
	}

	if err := r.HSet("a", "f1", "v1", "f2", "v2"); err != nil {
		t.Logf("HSet error: %v", err)
	}

	// get correct
	res := r.HGet("a", "f1")
	t.Logf("HGet a.f1: %s", res)

	res = r.HGet("a", "f3")
	if res == "" {
		t.Logf("HGet a.f3: [EMPTY]")
	} else {
		t.Logf("HGet a.f3: %s", res)
	}

	res = r.HGet("b", "f1")
	if res == "" {
		t.Logf("HGet b.f1: [EMPTY]")
	} else {
		t.Logf("HGet a.f1: %s", res)
	}
}
