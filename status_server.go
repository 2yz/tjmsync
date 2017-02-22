package tjmsync

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
	data, err := json.Marshal(GetGlobalJobPool())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
