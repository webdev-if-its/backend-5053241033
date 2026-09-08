package main

import "fmt"

// TODO(Level 1): lihat SOAL.md untuk kontrak lengkap tiap fungsi di bawah.
// Ganti setiap "panic" dengan implementasi yang benar.

func HitungSubtotal(qty int, hargaSatuan float64) float64 {
	return float64(qty) * hargaSatuan
}

func HitungTotalPesanan(qty []int, hargaSatuan []float64) float64 {
	if len(qty) != len(hargaSatuan) {
		return 0
	}
	var total float64
	for i := 0; i < len(qty); i++ {
		total += HitungSubtotal(qty[i], hargaSatuan[i])
	}
	return total
}

func TerapkanPajak(total float64, tarifPajak float64) float64 {
	return total * (1 + tarifPajak)
}

func HitungDiskon(total float64) float64 {
	if total >= 1000000 {
		return total * 0.10
	}
	if total >= 500000 {
		return total * 0.05
	}
	return 0
}

func TotalSetelahDiskon(qty []int, hargaSatuan []float64, tarifPajak float64) float64 {
	subtotal := HitungTotalPesanan(qty, hargaSatuan)
	diskon := HitungDiskon(subtotal)
	SetelahDiskon := subtotal - diskon
	return TerapkanPajak(SetelahDiskon, tarifPajak)
}

func ValidasiPesanan(qty []int, hargaSatuan []float64) (bool, string) {
	if len(qty) != len(hargaSatuan) {
		return false, "Jumlah item di daftar qty dan hargasatuan tidak sama"
	}
	for i := 0; i < len(qty); i++ {
		if qty[i] <= 0 {
			return false, fmt.Sprintf("qty pada index %d harus bilangan positif", i)
		}
		if hargaSatuan[i] <= 0 {
			return false, fmt.Sprintf("hargaSatuan pada index %d harus bilangan positif", i)
		}
	}
	return true, ""
}

func TentukanStatus(total float64) string {
	if total > 1000000 {
		return "Prioritas"
	}
	if total > 100000 {
		return "Reguler"
	}
	return "Hemat"
}

func RingkasanPesanan(qty []int, hargaSatuan []float64, tarifPajak float64) string {
	subtotal := HitungTotalPesanan(qty, hargaSatuan)
	diskon := HitungDiskon(subtotal)
	totalakhir := TotalSetelahDiskon(qty, hargaSatuan, tarifPajak)
	status := TentukanStatus(totalakhir)

	return fmt.Sprintf(
		"Subtotal   : Rp%.2f\nDiskon     : Rp%.2f\nTotal Akhir: Rp%.2f\nStatus     : %s\n",
		subtotal, diskon, totalakhir, status,
	)

}

// TODO(Level 9): signature ini SUDAH benar (cari tahu sendiri kenapa
// bentuknya begini - lihat SOAL.md) - tinggal implementasikan isinya.
func Total(harga ...float64) float64 {
	var total float64
	for i := 0; i < len(harga); i++ {
		total += harga[i]
	}
	return total
}

// TODO(Level 10, bonus): signature ini SUDAH benar (cari tahu sendiri
// kenapa ada dua nilai balik - lihat SOAL.md) - tinggal implementasikan isinya.
func HitungOngkosKirim(beratKg float64, jarakKm float64) (float64, error) {
}

func main() {
	fmt.Println("Sales Order Processor - pertemuan 2")
}
