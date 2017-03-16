package manager

import (
	"github.com/yezersky/grpool"
	"errors"
	"time"
	"context"
	"sync"
)

type MANAGER_STATE string

const (
	MANAGER_STATE_CREATED = "created"
	MANAGER_STATE_RUNNING = "running"
	MANAGER_STATE_STOPPED = "stopped"
)

type Manager struct {
	jobs   []*Job
	pool   *grpool.WorkerPool
	state  MANAGER_STATE
	ctx    context.Context
	cancel context.CancelFunc
	lock   sync.Mutex
}

func NewManager(configs []ConfigJob, pool *grpool.WorkerPool) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	var jobs []*Job
	for _, c := range configs {
		childCtx, childCancel := context.WithCancel(ctx)
		jobs = append(jobs, &Job{
			Host: c.Host,
			Version: c.Version,
			Image: c.Image,
			Name: c.Name,
			Command: c.Command,
			Env: c.Env,
			Volumes: c.Volumes,
			Interval: c.Interval,
			ctx: childCtx,
			cancel: childCancel,
		})
	}

	return &Manager{
		jobs: jobs,
		pool: pool,
		state: MANAGER_STATE_CREATED,
		ctx: ctx,
		cancel: cancel,
	}
}

func (m *Manager) Start() error {
	m.lock.Lock()
	switch m.state {
	case MANAGER_STATE_RUNNING:
		m.lock.Unlock()
		return errors.New("manager is running")
	case MANAGER_STATE_STOPPED:
		m.lock.Unlock()
		return errors.New("manager cannot be restarted")
	}
	m.state = MANAGER_STATE_RUNNING
	m.lock.Unlock()

	go func() {
		tick := time.Tick(time.Second)
		for {
			select {
			case <-m.ctx.Done():
				return
			case <-tick:
				m.schedule()
			}
		}
	}()
	return nil
}

func (m *Manager) Stop() error {
	m.lock.Lock()
	switch m.state {
	case MANAGER_STATE_CREATED:
		m.lock.Unlock()
		return errors.New("manager is not running")
	case MANAGER_STATE_STOPPED:
		m.lock.Unlock()
		return errors.New("manager is stopped")
	}
	m.state = MANAGER_STATE_STOPPED
	m.lock.Unlock()

	m.cancel()

	return nil
}

func (m *Manager) schedule() {
	j := m.getNextJob()
	if j == nil {
		return
	}
	j.Status = JOB_STATUS_QUEUED
	err := m.pool.Queue(j)
	if err != nil {
		j.Status = JOB_STATUS_IDLE
		return
	}
}

func (m *Manager) getNextJob() *Job {
	var nextJob *Job = nil
	var nextTime time.Time = time.Now()
	jobs := m.jobs
	for _, j := range jobs {
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

func (m *Manager) GetJobs() []*Job {
	return m.jobs
}