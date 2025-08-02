package server

import (
	"encoding/json"
	"github.com/AndreySirin/-Effective-Mobile-/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func (s *Server) CreateSubs(w http.ResponseWriter, r *http.Request) {
	var req entity.SubsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subs, err := entity.ToDataBase(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := s.storage.CreateSubs(r.Context(), &subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) ReadSubs(w http.ResponseWriter, r *http.Request) {
	subsID := chi.URLParam(r, "id")
	id, err := uuid.Parse(subsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id == uuid.Nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	subs, err := s.storage.ReadSubs(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(subs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateSubs(w http.ResponseWriter, r *http.Request) {
	var req entity.SubsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subs, err := entity.ToDataBase(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subsID := chi.URLParam(r, "id")
	id, err := uuid.Parse(subsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id == uuid.Nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = s.storage.UpdateSubs(r.Context(), id, &subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) DeleteSubs(w http.ResponseWriter, r *http.Request) {
	subsID := chi.URLParam(r, "id")
	id, err := uuid.Parse(subsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.storage.DeleteSubs(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) ListSubs(w http.ResponseWriter, r *http.Request) {
	subs, err := s.storage.ListSubs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(subs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
