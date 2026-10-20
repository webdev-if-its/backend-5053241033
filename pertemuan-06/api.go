package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// taskInput adalah bentuk JSON yang boleh dikirim klien untuk POST/PUT.
// Sengaja TIDAK memuat ID -- ID bukan urusan klien.
type taskInput struct {
	Judul   string `json:"judul"`
	Selesai bool   `json:"selesai"`
}

// errorBody adalah bentuk JSON semua respons error: {"error":"pesan"}.
type errorBody struct {
	Error string `json:"error"`
}

// Error khusus lapisan HTTP (400 Bad Request).
var (
	errIDBukanAngka  = errors.New("id harus berupa angka")
	errBodyTidakJSON = errors.New("body bukan JSON yang valid")
)

// API menyatukan routing dan handler HTTP di atas sebuah TaskStore.
type API struct {
	store TaskStore
}

// NewAPI membuat API di atas store (interface, bukan MemoryStore).
func NewAPI(store TaskStore) *API {
	return &API{store: store}
}

// Routes mendaftarkan semua endpoint memakai pola Go 1.22 ("METHOD /path").
func (a *API) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", a.listTasks)
	mux.HandleFunc("POST /tasks", a.createTask)
	mux.HandleFunc("GET /tasks/{id}", a.getTask)
	mux.HandleFunc("PUT /tasks/{id}", a.updateTask)
	mux.HandleFunc("DELETE /tasks/{id}", a.deleteTask)
	return mux
}

// writeJSON menulis v sebagai JSON dengan status yang diminta.
// Urutan WAJIB: Header().Set(Content-Type) -> WriteHeader(status) -> Encode.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError menerjemahkan error aplikasi menjadi status + body errorBody.
// Memakai errors.Is supaya error yang dibungkus (%w) tetap dikenali.
func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTugasTidakDitemukan):
		writeJSON(w, http.StatusNotFound, errorBody{Error: ErrTugasTidakDitemukan.Error()})
	case errors.Is(err, ErrInputKosong):
		writeJSON(w, http.StatusBadRequest, errorBody{Error: ErrInputKosong.Error()})
	case errors.Is(err, errIDBukanAngka):
		writeJSON(w, http.StatusBadRequest, errorBody{Error: errIDBukanAngka.Error()})
	case errors.Is(err, errBodyTidakJSON):
		writeJSON(w, http.StatusBadRequest, errorBody{Error: errBodyTidakJSON.Error()})
	default:
		// Level 8 belum dikerjakan: pesan mentah masih dikirim ke klien.
		writeJSON(w, http.StatusInternalServerError, errorBody{Error: err.Error()})
	}
}

// GET /tasks
func (a *API) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks := a.store.List()
	if tasks == nil {
		tasks = []Task{}
	}
	writeJSON(w, http.StatusOK, tasks)
}

// POST /tasks
func (a *API) createTask(w http.ResponseWriter, r *http.Request) {
	var in taskInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, errBodyTidakJSON)
		return
	}
	if strings.TrimSpace(in.Judul) == "" {
		writeError(w, ErrInputKosong)
		return
	}

	task, err := a.store.Add(Task{
		Judul:   in.Judul,
		Selesai: in.Selesai,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

// GET /tasks/{id}
func (a *API) getTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, errIDBukanAngka)
		return
	}

	task, err := a.store.Get(id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

// PUT /tasks/{id}
func (a *API) updateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, errIDBukanAngka)
		return
	}

	var in taskInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, errBodyTidakJSON)
		return
	}
	if strings.TrimSpace(in.Judul) == "" {
		writeError(w, ErrInputKosong)
		return
	}

	task, err := a.store.Update(Task{
		ID:      id,
		Judul:   in.Judul,
		Selesai: in.Selesai,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

// DELETE /tasks/{id}
func (a *API) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, errIDBukanAngka)
		return
	}

	if err := a.store.Delete(id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
