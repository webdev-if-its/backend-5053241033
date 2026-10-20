package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// TODO: lihat SOAL.md untuk kontrak lengkap. Ganti setiap "panic" dengan
// implementasi yang benar; tambahkan import (encoding/json, errors, log,
// strconv, ...) sendiri kalau dibutuhkan.

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
func writeError(w http.ResponseWriter, err error) {
	panic("belum diimplementasikan")
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

	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		http.Error(w, "body bukan JSON yang valid", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(in.Judul) == "" {
		http.Error(w, "input tidak boleh kosong", http.StatusBadRequest)
		return
	}

	task, err := a.store.Add(Task{
		Judul:   in.Judul,
		Selesai: in.Selesai,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

// GET /tasks/{id}
func (a *API) getTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id harus berupa angka", http.StatusBadRequest)
		return
	}

	task, err := a.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// PUT /tasks/{id}
func (a *API) updateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id harus berupa angka", http.StatusBadRequest)
		return
	}

	var in taskInput
	err = json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		http.Error(w, "body bukan JSON yang valid", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(in.Judul) == "" {
		http.Error(w, "input tidak boleh kosong", http.StatusBadRequest)
		return
	}

	task, err := a.store.Update(Task{
		ID:      id,
		Judul:   in.Judul,
		Selesai: in.Selesai,
	})
	if err != nil {
		if err == ErrTugasTidakDitemukan {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err == ErrInputKosong {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// DELETE /tasks/{id}
func (a *API) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id harus berupa angka", http.StatusBadRequest)
		return
	}

	err = a.store.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
