package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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
	t.Judul = strings.TrimSpace(t.Judul)
	if t.Judul == "" {
		return Task{}, ErrInputKosong
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `INSERT INTO tasks (judul, selesai) VALUES ($1, $2) RETURNING id`
	err := s.db.QueryRowContext(ctx, query, t.Judul, t.Selesai).Scan(&t.ID)
	if err != nil {
		return Task{}, fmt.Errorf("add task: %w", err)
	}

	return t, nil
}

// Get mengambil satu tugas menurut ID. (Level 4)
func (s *PostgresStore) Get(id int) (Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var t Task
	query := `SELECT id, judul, selesai FROM tasks WHERE id = $1`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Judul, &t.Selesai)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrTugasTidakDitemukan
		}
		return Task{}, fmt.Errorf("get task: %w", err)
	}

	return t, nil
}

// List mengambil semua tugas berurutan menurut ID. (Level 4)
func (s *PostgresStore) List() ([]Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, judul, selesai FROM tasks ORDER BY id`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	tugas := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Judul, &t.Selesai); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tugas = append(tugas, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return tugas, nil
}

// Update mengganti judul dan status selesai tugas ber-ID t.ID. (Level 5)
func (s *PostgresStore) Update(t Task) (Task, error) {
	t.Judul = strings.TrimSpace(t.Judul)
	if t.Judul == "" {
		return Task{}, ErrInputKosong
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `UPDATE tasks SET judul = $1, selesai = $2 WHERE id = $3`
	result, err := s.db.ExecContext(ctx, query, t.Judul, t.Selesai, t.ID)
	if err != nil {
		return Task{}, fmt.Errorf("update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Task{}, fmt.Errorf("update task rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return Task{}, ErrTugasTidakDitemukan
	}

	return t, nil
}

// TambahBanyak menyimpan banyak tugas sekaligus memakai satu prepared
// statement. (Level 6 -- atomik di Level 8)
func (s *PostgresStore) TambahBanyak(judul []string) ([]Task, error) {
	if len(judul) == 0 {
		return []Task{}, nil
	}

	// Validasi awal semua judul
	trimmed := make([]string, len(judul))
	for i, j := range judul {
		j = strings.TrimSpace(j)
		if j == "" {
			return nil, ErrInputKosong
		}
		trimmed[i] = j
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, `INSERT INTO tasks (judul) VALUES ($1) RETURNING id, judul, selesai`)
	if err != nil {
		return nil, fmt.Errorf("prepare tambah banyak: %w", err)
	}
	defer stmt.Close()

	var hasil []Task
	for _, j := range trimmed {
		var t Task
		err := stmt.QueryRowContext(ctx, j).Scan(&t.ID, &t.Judul, &t.Selesai)
		if err != nil {
			return nil, fmt.Errorf("exec tambah banyak: %w", err)
		}
		hasil = append(hasil, t)
	}

	if hasil == nil {
		return []Task{}, nil
	}

	return hasil, nil
}

// Delete menghapus tugas dan mencatat aksinya ke task_logs DALAM SATU
// TRANSAKSI. (Level 7-8)
func (s *PostgresStore) Delete(id int) error {
	panic("belum diimplementasikan")
}
