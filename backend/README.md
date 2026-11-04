# Backend AgriSense

Jalankan perintah Go dari direktori `backend/`.

## Status perbaikan auth

Tahap 1–2 menyediakan konfigurasi JWT dan validasi algoritma/claims. Tahap 3–5 menambahkan:

- Error domain serta respons error terstruktur.
- Hasil service internal `AuthSession`; hanya `Response` ditulis ke JSON.
- Transaksi register: user dan refresh token tersimpan bersama atau keduanya dibatalkan.
- Rotasi refresh dengan PostgreSQL `FOR UPDATE`: satu token lama hanya dapat dikonsumsi sekali.
- Rollback pencabutan jika penyimpanan token pengganti gagal.
- Expiry refresh absolut sejak register/login, mengikuti `JWT_REFRESH_EXPIRATION`; rotasi tidak memperpanjangnya.
- Expiry JWT dibatasi akhir sesi, dan `expires_in` mengikuti umur token yang benar-benar diterbitkan.
- Cookie memakai expiry aktual sesi. Kegagalan refresh tidak menghapus cookie yang mungkin sudah diganti oleh request paralel.

`InitDB()` sekarang mengembalikan error. Startup dihentikan ketika koneksi atau migration gagal, sehingga repository tidak menerima pointer database nil.

## Kontrak endpoint saat ini

| Endpoint | Kredensial | Sukses |
|---|---|---|
| `POST /api/v1/auth/register` | JSON nama/email/password | 201 + access token JSON + refresh cookie |
| `POST /api/v1/auth/login` | JSON email/password | 200 + access token JSON + refresh cookie |
| `POST /api/v1/auth/refresh` | Cookie `refresh_token` | 200 + access token JSON + cookie pengganti |
| `POST /api/v1/auth/logout` | Access token Bearer dan refresh cookie | 200 + penghapusan cookie |
| `GET /api/v1/auth/me` | Access token Bearer | 200 + profil pengguna |

Logout masih memerlukan access token dan masih mengembalikan 200; perubahan route/logout 204, kebijakan Origin/CSRF/CORS, dan konfigurasi cookie production lengkap adalah tahap berikutnya. Penanganan error logout sudah tidak menelan kegagalan pencabutan. JWT yang sudah terbit tetap berlaku sampai expiry; logout bukan pencabutan access token instan atau seluruh keluarga sesi.

Respons error auth berbentuk:

```json
{"error":{"code":"INVALID_CREDENTIALS","message":"Email atau kata sandi tidak valid"}}
```

| Kondisi | Status/code |
|---|---|
| Input tidak valid, termasuk password registrasi >72 byte | 400 / `INVALID_INPUT` |
| Email sudah ada | 409 / `EMAIL_ALREADY_REGISTERED` |
| Kredensial salah | 401 / `INVALID_CREDENTIALS` |
| Access token tidak sah | 401 / `INVALID_ACCESS_TOKEN` |
| Refresh token hilang/rusak/kedaluwarsa/dicabut | 401 / `INVALID_REFRESH_TOKEN` |
| Pengguna tidak ditemukan | 404 / `USER_NOT_FOUND` |
| Gangguan koneksi database sementara yang dikenali | 503 / `AUTH_SERVICE_UNAVAILABLE` |
| Kegagalan internal lainnya | 500 / `INTERNAL_ERROR` |

Client harus membaca code/status, bukan mencocokkan pesan. Semua respons auth memakai `Cache-Control: no-store`. Detail database tidak dikirim ke client. Refresh gagal tidak mengubah cookie, termasuk ketika menerima 401. Koordinasikan refresh dalam client agar request paralel tidak saling mengganggu; jangan menghapus cookie dari respons refresh yang kalah race.

## Unique index refresh token

Model GORM sekarang mempunyai unique index pada `refresh_tokens.token_hash`. Pada database baru, `AutoMigrate` membuatnya. Untuk database existing, buat backup dan jalankan pemeriksaan/migration terkontrol sebelum menjalankan versi baru:

```bash
psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f migrations/001_refresh_token_hash_unique.sql
```

Migration berhenti jika ada hash duplikat dan tidak menghapus record otomatis. Periksa konflik melalui prosedur administrasi yang sesuai; jangan mencetak nilai token/hash ke log publik. Jika migration startup gagal, aplikasi berhenti dan tidak menerima trafik. Tidak ada migration yang diterapkan ke database aplikasi saat pengembangan perubahan ini.

## Pengujian

Unit/handler test serta kompilasi semua paket:

```bash
go test ./...
```

Integration test PostgreSQL dilewati jika `AUTH_TEST_DATABASE_DSN` tidak diisi. Untuk menjalankannya, sediakan PostgreSQL khusus pengujian, kemudian:

```bash
export AUTH_TEST_DATABASE_DSN='host=127.0.0.1 port=5432 user=testuser password=testpassword dbname=agrisense_auth_test sslmode=disable'
go test ./internal/service -run TestAuth -count=1 -v
go test -race ./...
go vet ./...
```

Nilai pada contoh adalah placeholder untuk lingkungan test, bukan kredensial production. Test membuat schema unik `auth_test_<uuid>` per kasus, lalu menghapus hanya schema tersebut. User database test harus boleh membuat schema/function/trigger dan membaca aktivitas koneksinya sendiri. Jangan gunakan DSN database production.

Integration test memeriksa rollback register/refresh, email register bersamaan, dua refresh yang benar-benar menunggu row lock, expiry setelah menunggu lock, alur login/profil/logout, batas expiry sesi, unique hash, serta migration yang menolak duplikasi tanpa menghapus data. `go test -race` memeriksa race memori Go; bukti locking berasal dari integration test PostgreSQL.
