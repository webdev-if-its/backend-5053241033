package main

import (
	"cmp"
	"errors"
	"fmt"
	"strings"
	"time"
)

// TODO: lihat SOAL.md untuk kontrak lengkap tiap fungsi/method di bawah.
// Ganti setiap "panic" dengan implementasi yang benar. Tambahkan import
// (mis. "strings") sendiri kalau memang dibutuhkan.

// ErrTugasTidakDitemukan dikembalikan ketika ID tugas tidak ada.
var ErrTugasTidakDitemukan = errors.New("tugas tidak ditemukan")

// ErrInputKosong dikembalikan ketika judul kosong (atau hanya spasi).
var ErrInputKosong = errors.New("input tidak boleh kosong")

// ErrDaftarKosong dikembalikan oleh Max ketika slice-nya kosong.
var ErrDaftarKosong = errors.New("daftar kosong")

// Timestamps dipakai lewat embedding (composition) di Task.
type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Touch menyetel UpdatedAt ke waktu sekarang. (Level 9)
func (ts *Timestamps) Touch() {
	panic("belum diimplementasikan")
}

// Task merepresentasikan satu tugas.
type Task struct {
	ID      int
	Judul   string
	Selesai bool
	Timestamps
}

// NewTask membuat Task baru dari judul. (Level 1)
func NewTask(judul string) (Task, error) {
	judul = strings.TrimSpace(judul)
	if judul == "" {
		return Task{}, ErrInputKosong
	}

	t := Task{
		Judul: judul,
		Timestamps: Timestamps{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	return t, nil
}

// MarkDone menandai tugas selesai. (Level 2)
func (t *Task) MarkDone() {
	t.Selesai = true
}

// Rename mengganti judul tugas. (Level 3)
func (t *Task) Rename(judul string) error {
	judul = strings.TrimSpace(judul)
	if judul == "" {
		return ErrInputKosong
	}

	t.Judul = judul

	return nil
}

// String membuat Task memenuhi fmt.Stringer. (Level 10)
func (t Task) String() string {
	panic("belum diimplementasikan")
}

// TaskStore adalah kontrak penyimpanan tugas. JANGAN diubah -- yang kalian
// tulis adalah implementasinya (MemoryStore).
type TaskStore interface {
	Add(t Task) (Task, error)
	Get(id int) (Task, error)
	List() []Task
	Delete(id int) error
}

// MemoryStore menyimpan tugas di memori.
type MemoryStore struct {
	tasks  []Task
	nextID int
}

// Pastikan *MemoryStore memenuhi TaskStore -- kalau tidak, kode gagal
// dikompilasi di sini, bukan baru ketahuan saat dijalankan.
var _ TaskStore = (*MemoryStore)(nil)

// NewMemoryStore membuat store kosong.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks:  []Task{},
		nextID: 1,
	}
}

// Add menyimpan tugas dan memberinya ID baru. (Level 4)
func (m *MemoryStore) Add(t Task) (Task, error) {
	t.Judul = strings.TrimSpace(t.Judul)
	if t.Judul == "" {
		return Task{}, ErrInputKosong
	}

	t.ID = m.nextID
	m.nextID++

	m.tasks = append(m.tasks, t)
	return t, nil
}

// Get mencari tugas menurut ID. (Level 4)
func (m *MemoryStore) Get(id int) (Task, error) {
	for _, t := range m.tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return Task{}, ErrTugasTidakDitemukan
}

// List mengembalikan seluruh tugas. (Level 5)
func (m *MemoryStore) List() []Task {
	if len(m.tasks) == 0 {
		return []Task{}
	}

	list := make([]Task, len(m.tasks))
	copy(list, m.tasks)
	return list
}

// Delete menghapus tugas menurut ID. (Level 6)
func (m *MemoryStore) Delete(id int) error {
	for i, t := range m.tasks {
		if t.ID == id {
			m.tasks = append(
				m.tasks[:i],
				m.tasks[i+1:]...)
			return nil
		}
	}
	return ErrTugasTidakDitemukan
}

// Filter mengembalikan elemen xs yang lolos pred. (Level 7)
func Filter[T any](xs []T, pred func(T) bool) []T {
	result := make([]T, 0) // Inisialisasi slice non-nil (bukan nil)
	for _, v := range xs {
		if pred(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map mengubah tiap elemen xs dengan f. (Level 7)
func Map[T, U any](xs []T, f func(T) U) []U {
	result := make([]U, 0)
	for _, v := range xs {
		result = append(result, f(v))
	}
	return result
}

// Contains melaporkan apakah v ada di xs. (Level 8)
func Contains[T comparable](xs []T, v T) bool {
	panic("belum diimplementasikan")
}

// Max mengembalikan elemen terbesar di xs. (Level 8)
func Max[T cmp.Ordered](xs []T) (T, error) {
	panic("belum diimplementasikan")
}

// Gabung menyambung String() tiap elemen dengan pemisah sep. (Level 10)
func Gabung[T fmt.Stringer](xs []T, sep string) string {
	panic("belum diimplementasikan")
}

func main() {
	fmt.Println("Task Manager v1 - pertemuan 4")

	var store TaskStore = NewMemoryStore()
	for _, judul := range []string{"Belajar struct", "Belajar interface"} {
		t, err := NewTask(judul)
		if err != nil {
			fmt.Println("gagal membuat tugas:", err)
			continue
		}
		if _, err := store.Add(t); err != nil {
			fmt.Println("gagal menyimpan tugas:", err)
		}
	}
	for _, t := range store.List() {
		fmt.Println(t)
	}
}
