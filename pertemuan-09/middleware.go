package main

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// TODO: lihat SOAL.md untuk kontrak lengkap. Ganti setiap "panic" dengan
// implementasi yang benar; tambahkan import (log, time, runtime/debug, mime,
// ...) sendiri kalau dibutuhkan.

// Middleware membungkus sebuah http.Handler dengan perilaku tambahan.
type Middleware func(http.Handler) http.Handler

// Chain membungkus h dengan mws. Urutan: middleware PERTAMA di daftar adalah
// yang PALING LUAR (dijalankan paling dulu saat request masuk). (Level 3)
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// statusRecorder membungkus ResponseWriter untuk mengingat status code yang
// ditulis handler. Status awal (kalau handler tak pernah memanggil
// WriteHeader) adalah 200. (Level 1-2)
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

// WriteHeader mencatat status lalu meneruskannya ke ResponseWriter asli.
// Hanya panggilan PERTAMA yang dicatat. (Level 2)
func (sr *statusRecorder) WriteHeader(code int) {
	if !sr.wroteHeader {
		sr.status = code
		sr.wroteHeader = true
	}
	sr.ResponseWriter.WriteHeader(code)
}

// Logging mencatat SATU baris per request lewat log.Printf dengan format
// "%s %s %d %v" -> method, path, status, durasi. (Level 1-2)
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mulai := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		durasi := time.Since(mulai)

		log.Printf("%s %s %d %v",
			r.Method,
			r.URL.Path,
			recorder.status,
			durasi,
		)
	})
}

// Recover menangkap panic dari handler di bawahnya: catat ke log (pesan +
// stack trace), balas 500 JSON generik. (Level 4-6)
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v\n%s", rec, debug.Stack())
				writeJSON(w, http.StatusInternalServerError, errorBody{Error: "terjadi kesalahan internal"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RequireJSON menolak POST/PUT/PATCH yang Content-Type-nya bukan
// application/json dengan 415. (Level 9)
func RequireJSON(next http.Handler) http.Handler {
	panic("belum diimplementasikan")
}

// NewHandler menyusun API + seluruh middleware menjadi satu http.Handler yang
// siap dipakai server. (Level 10)
func NewHandler(store TaskStore) http.Handler {
	panic("belum diimplementasikan")
}
