package redis

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
)

func setup() *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return s
}

// TestDirectPing test the redis for tests is available
func TestDirectPing(t *testing.T) {
	s := setup()
	defer s.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     s.Addr(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := rdb.Ping().Result()
	fmt.Println(pong, err)
}

func TestPing(t *testing.T) {
	s := setup()
	defer s.Close()

}
