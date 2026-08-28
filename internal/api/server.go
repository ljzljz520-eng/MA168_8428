package api

import (
	"bookstore/recommendation/internal/model"
	"bookstore/recommendation/internal/recommendation"
	"bookstore/recommendation/internal/workflow"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	Service *workflow.Service
	Logger  *log.Logger
}

func New(service *workflow.Service, logger *log.Logger) *Server {
	return &Server{Service: service, Logger: logger}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/records", s.records)
	mux.HandleFunc("/records/", s.record)
	mux.HandleFunc("/recommendations", s.recommendations)
	mux.HandleFunc("/reports", s.reports)
	return requestLog(s.Logger, mux)
}

func requestLog(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Printf("%s %s", r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		query := model.Query{StoreID: r.URL.Query().Get("store"), Text: r.URL.Query().Get("q"), Genre: r.URL.Query().Get("genre"), Status: model.Status(r.URL.Query().Get("status"))}
		page, err := s.Service.Catalog.Search(query)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, page)
	case http.MethodPost:
		var input model.RecordInput
		if err := readJSON(r, &input); err != nil {
			writeError(w, err)
			return
		}
		record, err := s.Service.Register(input, 1, actor(r))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, record)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) record(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/records/")
	if id == "" {
		writeErrorStatus(w, http.StatusBadRequest, "record id is required")
		return
	}
	switch r.Method {
	case http.MethodGet:
		record, err := s.Service.Catalog.Get(id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, record)
	case http.MethodPatch:
		var input model.RecordInput
		if err := readJSON(r, &input); err != nil {
			writeError(w, err)
			return
		}
		record, err := s.Service.Update(id, input, actor(r), 1)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, record)
	case http.MethodPost:
		action := r.URL.Query().Get("action")
		var record model.Record
		var err error
		switch action {
		case "submit":
			record, err = s.Service.Submit(id, actor(r), 1)
		case "publish":
			record, err = s.Service.Publish(id, actor(r), 1)
		case "archive":
			record, err = s.Service.Archive(id, actor(r), 1)
		case "approve":
			record, err = s.Service.Review(id, actor(r), true, r.URL.Query().Get("note"), 1)
		case "reject":
			record, err = s.Service.Review(id, actor(r), false, r.URL.Query().Get("note"), 1)
		default:
			writeErrorStatus(w, http.StatusBadRequest, "unknown action")
			return
		}
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, record)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) recommendations(w http.ResponseWriter, r *http.Request) {
	request := recommendation.Request{StoreID: r.URL.Query().Get("store"), Genre: r.URL.Query().Get("genre"), IncludePublished: r.URL.Query().Get("published") == "true"}
	collection, err := s.Service.Recommendations(request, 1)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, collection)
}

func (s *Server) reports(w http.ResponseWriter, r *http.Request) {
	report, err := s.Service.BuildReport(r.URL.Query().Get("record"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func actor(r *http.Request) string {
	if value := r.Header.Get("X-Actor"); value != "" {
		return value
	}
	return "web"
}

func readJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error) {
	writeErrorStatus(w, http.StatusBadRequest, err.Error())
}

func writeErrorStatus(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
