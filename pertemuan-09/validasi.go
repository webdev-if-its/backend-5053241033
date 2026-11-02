package main

import (
	"strings"
	"unicode/utf8"
)

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "validasi gagal"
}

func ValidasiTugas(in taskInput) error {
	fields := make(map[string]string)

	judul := strings.TrimSpace(in.Judul)

	if judul == "" {
		fields["judul"] = "tidak boleh kosong"
	} else if utf8.RuneCountInString(judul) > 100 {
		fields["judul"] = "maksimal 100 karakter"
	}

	if len(fields) > 0 {
		return &ValidationError{
			Fields: fields,
		}
	}

	return nil
}
