package tjmsync

import (
	"bytes"
	"encoding/json"
	"log"
	"time"
)

type Worker struct {
	JobPool    *JobPool
	JobChannel chan *Job
}

func NewWorker(jobPool *JobPool) *Worker {
	return &Worker{
		JobPool:    jobPool,
		JobChannel: make(chan *Job, 1),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			j := w.JobPool.ScheduleNextJob()
			if j != nil {
				j.Status = JOB_STATUS_QUEUED
				w.JobChannel <- j
			}
			time.Sleep(time.Second)
		}
	}()
	go func() {
		for j := range w.JobChannel {
			work := &Work{Job: j}
			work.Start()
		}
	}()
}

type Work struct {
	Job        *Job
	StartedAt  time.Time
	FinishedAt time.Time
	ExitStatus JOB_EXIT_STATUS
	Error      error
	ErrStr     string
	Stdout     bytes.Buffer
	Stderr     bytes.Buffer
}

func (w *Work) Start() {
	w.pre_exec()
	w.exec()
	w.post_exec()
}

func (w *Work) pre_exec() {
	w.StartedAt = time.Now()
	w.Job.LastStartedAt = w.StartedAt
	w.Job.Status = JOB_STATUS_RUNNING
}

func (w *Work) exec() {
	// TODO check job nil
	cmd := w.Job.Command.GetCmd()
	cmd.Env = w.Job.Env
	cmd.Stdout = &w.Stdout
	cmd.Stderr = &w.Stderr
	w.Error = cmd.Run()
}

func (w *Work) post_exec() {
	w.FinishedAt = time.Now()
	if w.Error == nil {
		w.ExitStatus = JOB_EXIT_STATUS_SUCCESS
	} else {
		w.ExitStatus = JOB_EXIT_STATUS_FAIL
		w.ErrStr = w.Error.Error()
	}

	w.log()

	w.Job.LastExitStatus = w.ExitStatus
	if w.ExitStatus == JOB_EXIT_STATUS_SUCCESS {
		w.Job.LastSucceededAt = w.FinishedAt
	}
	w.Job.LastFinishedAt = w.FinishedAt
	w.Job.Status = JOB_STATUS_IDLE
}

func (w *Work) log() {
	result := &JobResult{
		JobName:    w.Job.Name,
		JobCommand: w.Job.Command,
		JobEnv:     w.Job.Env,
		ExitStatus: w.ExitStatus,
		Error:      w.ErrStr,
		Stdout:     w.Stdout.String(),
		Stderr:     w.Stderr.String(),
		Duration:   w.FinishedAt.Sub(w.StartedAt),
	}
	data, err1 := json.Marshal(result)
	if err1 != nil {
		log.Println("Work log: ", err1)
		return
	}
	err2 := JobLogger.Log(string(data))
	if err2 != nil {
		log.Println("Work log: ", err2)
		log.Println("Work log result: ", string(data))
		return
	}
}
