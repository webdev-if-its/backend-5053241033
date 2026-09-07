package main

import (
	"fmt"
	"os"
	"runtime"
)

// TODO(Level 4): ganti dengan NRP kalian sendiri, contoh: "5025201012"
const NRP = "5053241033"

// TODO(Level 3): kembalikan args[0] kalau ada isinya, kalau tidak kembalikan fallback.
func ResolveNama(args []string, fallback string) string {
	if len(args) > 0 {
		return args[0]
	}
	return fallback
}

// TODO(Level 7): kembalikan kalimat sapaan untuk nama yang diberikan,
// contoh: "Halo, Budi! Selamat datang di kelas Backend."
// Baru diimplementasikan saat mengerjakan level 7-9, lihat SOAL.md.
func Sapa(nama string) string {
	return fmt.Sprintf("Halo, %s! Selamat datang di kelas Backend.", nama)
}

// TODO(Level 5): gabungkan Nama, NRP, dan hasil runtime.Version() jadi satu
// string siap cetak (lihat contoh format di SOAL.md).
func CetakInfo(nama string) string {
	return fmt.Sprintf("Nama: %s\nNRP: %s\n%s", nama, NRP, runtime.Version())
}

func main() {
	nama := ResolveNama(os.Args[1:], NRP)
	fmt.Println(CetakInfo(nama))
	fmt.Println(Sapa(nama))
}
