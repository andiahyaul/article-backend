package article

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

const maxLimit = 100

type Handler struct{ repo *Repository }

func NewHandler(repo *Repository) *Handler { return &Handler{repo: repo} }

func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /article/{$}", h.create)
	mux.HandleFunc("GET /article/{limit}/{offset}", h.list)
	mux.HandleFunc("GET /article/{id}", h.get)
	mux.HandleFunc("PUT /article/{id}", h.update)
	mux.HandleFunc("PATCH /article/{id}", h.update)
	mux.HandleFunc("POST /article/{id}", h.update)
	mux.HandleFunc("DELETE /article/{id}", h.delete)
	// POST /article/{id} is taken by update, so POST-delete gets its own path
	mux.HandleFunc("POST /article/{id}/delete", h.delete)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]any{"error": msg})
}

func (h *Handler) fail(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrNotFound) {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	log.Printf("internal error: %v", err)
	writeErr(w, http.StatusInternalServerError, "internal server error")
}

func pathID(r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	return id, err == nil && id > 0
}

func decode(w http.ResponseWriter, r *http.Request) (Article, bool) {
	var a Article
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&a); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON body")
		return a, false
	}
	if errs := a.Validate(); len(errs) > 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "fields": errs})
		return a, false
	}
	return a, true
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	a, ok := decode(w, r)
	if !ok {
		return
	}
	id, err := h.repo.Create(a)
	if err != nil {
		h.fail(w, err)
		return
	}
	w.Header().Set("Location", "/article/"+strconv.FormatInt(id, 10))
	writeJSON(w, http.StatusCreated, map[string]any{})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	limit, e1 := strconv.Atoi(r.PathValue("limit"))
	offset, e2 := strconv.Atoi(r.PathValue("offset"))
	if e1 != nil || e2 != nil || limit < 1 || offset < 0 {
		writeErr(w, http.StatusBadRequest, "limit must be >= 1 and offset >= 0")
		return
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	status := r.URL.Query().Get("status")
	if status != "" && !validStatus[status] {
		writeErr(w, http.StatusBadRequest, "status must be one of: publish, draft, thrash")
		return
	}
	items, total, err := h.repo.List(limit, offset, status)
	if err != nil {
		h.fail(w, err)
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}
	a, err := h.repo.Get(id)
	if err != nil {
		h.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}
	a, ok := decode(w, r)
	if !ok {
		return
	}
	if err := h.repo.Update(id, a); err != nil {
		h.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.repo.Delete(id); err != nil {
		h.fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}
