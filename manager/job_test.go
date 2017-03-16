package manager

import (
	"gopkg.in/redis.v5"
	"testing"
	"context"
)

func TestJobLoggerRedis_Log(t *testing.T) {
	JobLogger = NewJobLoggerRedis(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	}, "test")
	JobLogger.Log("{\"test\":\"test\"}")
}

func TestJob_Run(t *testing.T) {
	JobLogger = &JobLoggerDefault{}
	ctx, cancel := context.WithCancel(context.Background())
	j := &Job{
		//Host: "tcp://vm1.mirrors.tongji.edu.cn:2376",
		Version: "1.24",
		Image: "yezersky/tjmbase",
		Name: "test",
		Command: ParseCommand("echo hello world"),
		ctx: ctx,
		cancel: cancel,
	}
	go j.Run()
}