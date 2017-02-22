package tjmsync

import (
	"gopkg.in/redis.v5"
	"testing"
)

func TestJobLoggerRedis_Log(t *testing.T) {
	JobLogger = NewJobLoggerRedis(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	}, "test")
	JobLogger.Log("{\"test\":\"test\"}")
}
