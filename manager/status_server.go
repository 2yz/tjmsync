package manager

import (
	"encoding/json"
	"log"
	"net/http"
)

type Status struct{}

type StatusServer struct{}

func (s *StatusServer) Serve() {
	http.HandleFunc("/status/job", s.job)
	log.Fatal(http.ListenAndServe(global.Config.StatusServer.GetAddr(), nil))
}

func (s *StatusServer) job(w http.ResponseWriter, r *http.Request) {
	m := global.Manager
	data, err := json.Marshal(m.GetJobs())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
