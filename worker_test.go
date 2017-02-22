package tjmsync

import (
	"testing"
	"time"
)

func TestWorker_Start(t *testing.T) {
	JobLogger = &JobLoggerDefault{}
	jobPool := &JobPool{
		Jobs: []*Job{{
			Name: "Test",
			Command: Command{
				Path: "date",
			},
		}},
	}
	w := NewWorker(jobPool)
	w.Start()
	time.Sleep(3 * time.Second)
}

func TestWork_Start(t *testing.T) {
	JobLogger = &JobLoggerDefault{}
	w := &Work{Job: &Job{
		Name: "Test",
		Command: Command{
			Path: "date",
		},
	}}
	w.Start()
}
