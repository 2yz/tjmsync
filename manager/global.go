package manager

import "github.com/yezersky/grpool"

var global struct {
	Config     Config
	WorkerPool *grpool.WorkerPool
	Manager    *Manager
}

func SetGlobalConfig(config *Config) {
	global.Config = *config
}

func SetGlobalWorkerPool(wp *grpool.WorkerPool) {
	global.WorkerPool = wp
}

func SetGlobalManager(m *Manager) {
	global.Manager = m
}
