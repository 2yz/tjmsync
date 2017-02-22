package lib

var global struct {
	Config  Config
	JobPool JobPool
}

func InitGlobalConfig(config *Config) {
	global.Config = *config
}

func InitGlobalJobPool(jobs []Job) {
	l := len(jobs)
	global.JobPool.Jobs = make([]*Job, l)
	for i := 0; i < l; i++ {
		j := jobs[i]
		global.JobPool.Jobs[i] = &j
	}
}

func GetGlobalJobPool() *JobPool {
	return &global.JobPool
}
