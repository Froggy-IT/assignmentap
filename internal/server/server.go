package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"assignmentap/internal/model"
	"assignmentap/internal/store"
)

type Server struct {
	store      *store.Store[string, string]
	startTime  time.Time
	reqCounter int64
}

func NewServer(s *store.Store[string, string]) *Server {
	return &Server{
		store:     s,
		startTime: time.Now(),
	}
}

func (s *Server) RequestCount() int64 {
	return atomic.LoadInt64(&s.reqCounter)
}

func (s *Server) KeyCount() int {
	return s.store.Count()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&s.reqCounter, 1)

	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/data":
		s.postData(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/data":
		s.getAll(w)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/data/"):
		s.getOne(w, r)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/data/"):
		s.deleteOne(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/stats":
		s.stats(w)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) postData(w http.ResponseWriter, r *http.Request) {
	var d model.Data
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil || d.Key == "" || d.Value == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	s.store.Set(d.Key, d.Value)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

func (s *Server) getAll(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(s.store.Snapshot())
}

func (s *Server) getOne(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/data/")
	if val, ok := s.store.Get(key); ok {
		json.NewEncoder(w).Encode(map[string]string{key: val})
		return
	}
	http.NotFound(w, r)
}

func (s *Server) deleteOne(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/data/")
	if s.store.Delete(key) {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.NotFound(w, r)
}

func (s *Server) stats(w http.ResponseWriter) {
	uptime := int64(time.Since(s.startTime).Seconds())
	json.NewEncoder(w).Encode(model.Stats{
		Requests:      s.RequestCount(),
		Keys:          s.KeyCount(),
		UptimeSeconds: uptime,
	})
}
