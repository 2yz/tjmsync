package manager

import (
	"errors"
	"gopkg.in/redis.v5"
	"time"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"fmt"
	"bytes"
	"context"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/volume"
	"github.com/docker/go-connections/tlsconfig"
)

const JOB_PREFIX = "tjmsync_job_"

type Job struct {
	Host            string
	Version         string
	Image           string
	Name            string
	Command         Command
	Env             []string
	Volumes         []string
	Interval        TimeDuration
	Status          JOB_STATUS
	LastStartedAt   time.Time
	LastFinishedAt  time.Time
	LastSucceededAt time.Time
	LastExitCode    int
	ctx             context.Context
	cancel          context.CancelFunc
}

type JOB_STATUS string

const (
	JOB_STATUS_IDLE = "idle"
	JOB_STATUS_QUEUED = "queued"
	JOB_STATUS_RUNNING = "running"
)

func (j *Job) Run() {
	l := &JobLog{JobHost: j.Host, JobVersion:j.Version, JobImage:j.Image, JobName: j.Name, JobCommand: j.Command, JobEnv: j.Env, JobVolumes: j.Volumes}
	j.pre(l)
	j.exec(l)
	j.post(l)
}

func (j *Job) pre(l *JobLog) {
	l.StartedAt = time.Now()
	j.LastStartedAt = l.StartedAt
	j.Status = JOB_STATUS_RUNNING
}

func (j *Job) exec(l *JobLog) {
	defer func() {
		if r := recover(); r != nil {
			var buf bytes.Buffer
			fmt.Fprint(&buf, r)
			l.Error = buf.String()
			l.ExitCode = 1
		}
	}()

	cli, err := NewDockerClient(j.Host, j.Version)
	if err != nil {
		panic(err)
	}

	c, err := cli.ContainerInspect(j.ctx, JOB_PREFIX + j.Name)
	if err == nil {
		if err := cli.ContainerRemove(j.ctx, c.ID, types.ContainerRemoveOptions{Force: true});
			err != nil {
			panic(err)
		}
	}

	_, err = cli.ImagePull(j.ctx, j.Image, types.ImagePullOptions{})
	if err != nil {
		log.Println(err)
	}

	resp, err := cli.ContainerCreate(j.ctx, &container.Config{
		Image: j.Image,
		Cmd: j.Command.GetCmd(),
		Env: j.Env,
	}, &container.HostConfig{
		Binds: j.Volumes,
	}, nil, JOB_PREFIX + j.Name)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(j.ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	if _, err = cli.ContainerWait(j.ctx, resp.ID); err != nil {
		panic(err)
	}

	c, err = cli.ContainerInspect(j.ctx, JOB_PREFIX + j.Name)
	if err != nil {
		panic(err)
	}
	l.ExitCode = c.State.ExitCode

	out, err := cli.ContainerLogs(j.ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(out); err != nil {
		panic(err)
	}
	l.Logs = buf.String()

	if err = cli.ContainerRemove(j.ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
	 	panic(err)
	}
}

func (j *Job) post(l *JobLog) {
	l.FinishedAt = time.Now()
	l.Duration = TimeDuration{Duration:l.FinishedAt.Sub(l.StartedAt)}
	l.Log()
	j.LastFinishedAt = l.FinishedAt
	j.LastExitCode = l.ExitCode
	if l.ExitCode == 0 {
		j.LastSucceededAt = l.FinishedAt
	}
	j.Status = JOB_STATUS_IDLE
}

//func (j *Job) convertMount() ([]mount.Mount) {
//	var mounts []mount.Mount
//	for _, m := range j.Volumes {
//		mounts = append(mounts, m.mount)
//	}
//	return mounts
//}

func NewDockerClient(host string, version string) (*client.Client, error) {
	var c *http.Client
	if dockerCertPath := os.Getenv("DOCKER_CERT_PATH"); dockerCertPath != "" {
		options := tlsconfig.Options{
			CAFile:             filepath.Join(dockerCertPath, "ca.pem"),
			CertFile:           filepath.Join(dockerCertPath, "cert.pem"),
			KeyFile:            filepath.Join(dockerCertPath, "key.pem"),
			InsecureSkipVerify: os.Getenv("DOCKER_TLS_VERIFY") == "",
		}
		tlsc, err := tlsconfig.Client(options)
		if err != nil {
			return nil, err
		}

		c = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsc,
			},
		}
	}

	if host == "" {
		host = os.Getenv("DOCKER_HOST")
	}
	if host == "" {
		host = client.DefaultDockerHost
	}
	if version == "" {
		version = os.Getenv("DOCKER_API_VERSION")
	}
	if version == "" {
		version = client.DefaultVersion
	}

	return client.NewClient(host, version, c, nil)
}

func (j *Job) IsIdle() bool {
	return j.Status == JOB_STATUS_IDLE || j.Status == ""
}

type JobLog struct {
	JobHost    string
	JobVersion string
	JobImage   string
	JobName    string
	JobCommand Command
	JobEnv     []string
	JobVolumes []string
	ExitCode   int
	Error      string
	Logs       string
	StartedAt  time.Time
	FinishedAt time.Time
	Duration   TimeDuration
}

func (l *JobLog) Log() {
	data, err := json.Marshal(l)
	if err != nil {
		log.Println("JobLog Error:", err)
		return
	}
	err = JobLogger.Log(string(data))
	if err != nil {
		log.Println("JobLog:", string(data))
		log.Println("JobLog Error:", err)
		return
	}
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
	c := redis.NewClient(opt)
	return &JobLoggerRedis{client: c, key: key}
}

func (l *JobLoggerRedis) Log(val string) error {
	return l.client.LPush(l.key, val).Err()
}

func SetJobLogger(c ConfigLog) {
	if c.Type == CONFIG_LOG_TYPE_REDIS {
		JobLogger = NewJobLoggerRedis(&redis.Options{
			Addr: c.RedisAddr, Password: c.RedisPassword, DB: c.RedisDB,
		}, c.RedisKey)
	}
}

type TimeDuration struct {
	Duration time.Duration
}

func (d *TimeDuration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func (d *TimeDuration) MarshalText() ([]byte, error) {
	return []byte(d.Duration.String()), nil
}

type DockerMount struct {
	str   string
	mount mount.Mount
}

func ParseDockerMount(str string) DockerMount {
	m := DockerMount{}
	m.UnmarshalText([]byte(str))
	return m
}

func (m *DockerMount) UnmarshalText(text []byte) error {
	m.str = string(text)
	mountPoint, err := volume.ParseMountRaw(m.str, "")
	if err != nil {
		return err
	}
	m.mount = mountPoint.Spec
	return nil
}

func (m *DockerMount) MarshalText() ([]byte, error) {
	return []byte(m.str), nil
}