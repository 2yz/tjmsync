package tjmsync

import (
	"errors"
	"gopkg.in/redis.v5"
	"time"
)

type Job struct {
	Name            string
	Command         Command
	Env             []string
	Interval        duration
	Status          JOB_STATUS
	LastStartedAt   time.Time
	LastFinishedAt  time.Time
	LastSucceededAt time.Time
	LastExitStatus  JOB_EXIT_STATUS
}

type JOB_STATUS string

const (
	JOB_STATUS_IDLE = "idle"
	JOB_STATUS_QUEUED = "queued"
	JOB_STATUS_RUNNING = "running"
)

type JOB_EXIT_STATUS string

const (
	JOB_EXIT_STATUS_SUCCESS = "success"
	JOB_EXIT_STATUS_FAIL = "fail"
)

func (j *Job) IsIdle() bool {
	return j.Status == JOB_STATUS_IDLE || j.Status == ""
}

type duration struct {
	Duration time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type JobPool struct {
	Jobs []*Job
}

func (jp *JobPool) ScheduleNextJob() *Job {
	var nextJob *Job = nil
	var nextTime time.Time = time.Now()
	for _, j := range jp.Jobs {
		if !j.IsIdle() {
			continue
		}
		t := j.LastStartedAt.Add(j.Interval.Duration)
		if t.Before(nextTime) {
			nextJob = j
			nextTime = t
		}
	}
	return nextJob
}

var JobLogger JobLoggerInterface = &JobLoggerDefault{}

type JobLoggerInterface interface {
	Log(string) error
}

type JobLoggerDefault struct{}

func (l *JobLoggerDefault) Log(val string) error {
	return errors.New("Please Set Job Logger First!")
}

type JobLoggerRedis struct {
	client *redis.Client
	key    string
}

func NewJobLoggerRedis(opt *redis.Options, key string) *JobLoggerRedis {
	client := redis.NewClient(opt)
	return &JobLoggerRedis{client: client, key: key}
}

func (l *JobLoggerRedis) Log(val string) error {
	return l.client.LPush(l.key, val).Err()
}

type JobResult struct {
	JobName    string
	JobCommand Command
	JobEnv     []string
	ExitStatus JOB_EXIT_STATUS
	Error      string
	Stdout     string
	Stderr     string
	Duration   time.Duration
}

func InitJobLogger() {
	l := global.Config.Log
	if l.Type == CONFIG_LOG_TYPE_REDIS {
		JobLogger = NewJobLoggerRedis(&redis.Options{
			Addr: l.RedisAddr, Password: l.RedisPassword, DB: l.RedisDB,
		}, l.RedisKey)
	}
}
