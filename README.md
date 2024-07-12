To-Do Golang App
Aplikasi manajemen tugas sederhana dengan backend Go dan frontend Next.js.

Fitur

Otentikasi: Login dengan JWT token.

Manajemen Tugas:
Tampilan daftar tugas.
Buat, edit, dan hapus tugas.
Tandai tugas sebagai selesai.

Work Log Harian:
Tampilan daftar tugas harian.
Tandai tugas harian sebagai selesai atau belum selesai.
Refresh otomatis setiap hari pada tengah malam.

Riwayat Work Log:
Tampilan riwayat work log berdasarkan tanggal.


Prasyarat
Go: Pastikan Anda sudah menginstal Go di sistem Anda. Anda bisa mengunduh Go dari https://golang.org/dl/.
Node.js dan npm: Pastikan Anda sudah menginstal Node.js dan npm. Anda bisa mengunduh Node.js dari https://nodejs.org/.
MySQL: Pastikan Anda sudah menginstal dan menjalankan server MySQL.
HeidiSQL atau MySQL Client Lain: Anda akan membutuhkan alat ini untuk mengelola database MySQL.

Cara Menjalankan Proyek
Clone Repositori:
Bash
git clone https://github.com/username-anda/todo-golang.git

Setup Backend (Go):
Buat Database:
Buka HeidiSQL atau MySQL client lainnya.
Buat database baru dengan nama todo_golang.
Import file todo_golang.sql yang ada di dalam folder backend untuk membuat tabel-tabel yang dibutuhkan.

Konfigurasi Koneksi Database:
Buka file backend/database/database.go.
Sesuaikan nilai dbUser, dbPass, dan dbName dengan kredensial MySQL Anda.

Instal Dependensi:
cd backend
go mod download

Jalankan Server Backend:
go run main.go

Server backend akan berjalan di http://localhost:8080.

Setup Frontend (Next.js):
Install Dependensi:
cd frontend
npm install

Jalankan Server Frontend:
npm run dev

Server frontend akan berjalan di http://localhost:3000.

Akses Aplikasi:
Buka browser Anda dan akses http://localhost:3000.
Anda akan diarahkan ke halaman login. Masukkan username dan password yang sudah ditentukan di backend (admin dan admin).
Setelah login, Anda bisa mengakses halaman-halaman lain seperti dashboard, work log, dan work log history.
Catatan:

Pastikan port yang digunakan oleh backend dan frontend tidak bentrok. Jika terjadi bentrok, Anda bisa mengubah port di konfigurasi backend atau frontend.
Pastikan Anda sudah menjalankan server MySQL sebelum menjalankan backend Go.
Troubleshooting:

Jika terjadi error "parsing time" pada saat mengambil data work log, pastikan kolom date pada tabel work_log dan work_log_history di database Anda bertipe data DATE.
Semoga panduan ini membantu Anda menjalankan aplikasi To-Do Golang dengan lancar!