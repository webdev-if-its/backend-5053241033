package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // mendaftarkan driver "pgx" untuk database/sql
)

// TODO: lihat SOAL.md untuk kontrak lengkap. Ganti setiap "panic" dengan
// implementasi yang benar; tambahkan import (context, errors, fmt, strings, ...)
// sendiri kalau dibutuhkan.
//
// Aturan yang berlaku untuk SEMUA method:
//   - nama tabel ditulis TANPA awalan schema ("tasks", bukan "public.tasks"):
//     test menjalankan tiap level di schema terpisah lewat search_path;
//   - nilai dari luar SELALU lewat placeholder ($1, $2, ...), tidak pernah
//     disambung ke string SQL;
//   - tiap operasi memakai context dengan timeout dbTimeout;
//   - error dari database dibungkus dengan %w dan konteks ("add task: ...").

// dbTimeout adalah batas waktu tiap operasi database. Sengaja `var` (bukan
// const) supaya test bisa mengecilkannya.
var dbTimeout = 3 * time.Second

// Connect membuka pool koneksi ke PostgreSQL dan MEMASTIKAN koneksinya benar-
// benar bisa dipakai. (Level 1)
func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal open dsn: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("gagal ping database: %w", err)
	}

	return db, nil
}

// PostgresStore adalah TaskStore yang menyimpan data di PostgreSQL.
type PostgresStore struct {
	db *sql.DB
}

// Pastikan *PostgresStore memenuhi TaskStore.
var _ TaskStore = (*PostgresStore)(nil)

// NewPostgresStore membungkus db menjadi TaskStore.
func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// Migrate membuat tabel tasks dan task_logs kalau belum ada. Aman dipanggil
// berkali-kali. (Level 2)
func (s *PostgresStore) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id      INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		judul   TEXT NOT NULL,
		selesai BOOLEAN NOT NULL DEFAULT FALSE
	);
	CREATE TABLE IF NOT EXISTS task_logs (
		id      INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		task_id INT NOT NULL,
		aksi    TEXT NOT NULL,
		waktu   TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("migrate gagal: %w", err)
	}
	return nil
}

// Add menyimpan tugas baru; ID dibuat oleh database. (Level 3)
func (s *PostgresStore) Add(t Task) (Task, error) {
	panic("belum diimplementasikan")
}

// Get mengambil satu tugas menurut ID. (Level 4)
func (s *PostgresStore) Get(id int) (Task, error) {
	panic("belum diimplementasikan")
}

// List mengambil semua tugas berurutan menurut ID. (Level 4)
func (s *PostgresStore) List() ([]Task, error) {
	panic("belum diimplementasikan")
}

// Update mengganti judul dan status selesai tugas ber-ID t.ID. (Level 5)
func (s *PostgresStore) Update(t Task) (Task, error) {
	panic("belum diimplementasikan")
}

// TambahBanyak menyimpan banyak tugas sekaligus memakai satu prepared
// statement. (Level 6 -- atomik di Level 8)
func (s *PostgresStore) TambahBanyak(judul []string) ([]Task, error) {
	panic("belum diimplementasikan")
}

// Delete menghapus tugas dan mencatat aksinya ke task_logs DALAM SATU
// TRANSAKSI. (Level 7-8)
func (s *PostgresStore) Delete(id int) error {
	panic("belum diimplementasikan")
}
