# backend-5053241033

Repo tugas mata kuliah **Pengembangan Backend Dasar**, dibuat dari template [`webdev-if-its/backend-template`](https://github.com/webdev-if-its/backend-template). Ganti judul di atas jadi nama repo kalian sendiri (`backend-nrp`, contoh: `backend-5025201012`).

## Aturan Umum

- Tugas tiap pertemuan disimpan di folder `pertemuan-XX/` pada repo ini.
- Commit message wajib menyebut level yang dicapai: `pertemuan-XX: level N selesai`.
- Deadline push: sebelum pertemuan berikutnya dimulai.
- Semua level dicek otomatis lewat `go test` — baca `pertemuan-XX/SOAL.md` tiap minggu untuk detail levelnya.

## Mengambil Pertemuan Baru Tiap Minggu

Repo ini **tidak otomatis sinkron** dengan template dosen. Begitu ada pertemuan baru, jalankan (ganti `pertemuan-02` sesuai minggu berjalan):

```bash
git fetch https://github.com/webdev-if-its/backend-template.git main
git checkout FETCH_HEAD -- pertemuan-02
```

Perintah ini **aman dijalankan kapan pun** — tidak akan menimpa folder pertemuan lain yang sudah kalian kerjakan, karena hanya mengambil folder yang disebutkan. Setelah itu, commit folder barunya seperti biasa.

Kalau dosen memperbaiki sesuatu di pertemuan yang sudah dirilis (mis. ada bug di test), biasanya cukup ambil ulang file yang diperbaiki saja, bukan seluruh folder — akan diumumkan file mana yang berubah.

---

Bagian di bawah ini **isi bertahap** sesuai level yang sedang kalian kerjakan (lihat `pertemuan-01/SOAL.md`) — heading-nya dicek otomatis, jangan diganti namanya.

## Identitas
- Nama:Arya Ramadhan
- NRP:5053241033
- Kelas:M

## Commit vs Push
git commit menyimpan perubahan file ke dalam riwayat lokal komputer kita. Sedangkan git push dia mengirim commit ke repositori remote di GitHub. Jika seseorang commit tetapi lupa di push, perubahan kode tidak akan terlihat oleh rekan tim sehingga mereka tidak bisa melanjutkan pekerjaan.

## Reproducibility
Reproducibility adalah program yang kita buat bisa dijalankan oleh orang lain dan hasilnya sama seperti punya kita sendiri. Tapi jika anggota tim menggunakkan go yang beda versi, akan terjadi errror karena mungkin ada sintaks atau fungsi baru yang tidak didukung di versi go lama

## Catatan Merge Conflict
(tulis di sini)

## Kenapa .gitignore Penting
(tulis di sini)

## Refleksi
(tulis di sini)
