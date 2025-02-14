package seq

import (
	"context"
	"github.com/redis/go-redis/v9"
	"ppim/internal/logic/store"
	"testing"
	"time"
)

func TestRedisSequence_Next(t *testing.T) {
	r := redis.NewClient(&redis.Options{
		Addr:            "127.0.0.1:6379",
		DB:              0,
		MinIdleConns:    2,
		MaxActiveConns:  10,
		ConnMaxLifetime: 5 * time.Minute,
	})

	s1 := NewStepSequence(10, WithDefaultSeqStorage(r, new(store.Conversation)))
	for range 123 {
		go func() {
			seq, err := s1.Next(context.TODO(), "single|20906144@19726487")
			if err != nil {
				t.Error(err)
				return
			}
			t.Logf("s1 - seq=%d \n", seq)
		}()
	}

	//s2 := NewStepSequence(r, 10)
	for range 13 {
		go func() {
			seq, err := s1.Next(context.TODO(), "single|20906144@19726487")
			if err != nil {
				t.Error(err)
				return
			}
			t.Logf("s2 - seq=%d \n", seq)
		}()
	}

	select {}
}
