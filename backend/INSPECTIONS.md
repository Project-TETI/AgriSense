# Pemeriksaan foto dan integrasi ML

`POST /api/v1/inspections` memerlukan access token Bearer dan `multipart/form-data` dengan tepat satu `plant` (kode tanaman katalog) dan satu `file` (JPEG atau PNG, maksimal 5 MiB). Isi gambar diperiksa, termasuk batas 20 juta piksel. Batas awal adalah 20 request per menit per IP per proses; beberapa replika memerlukan pembatasan bersama di gateway.

Respons 201 berisi `inspection_id`, `plant`, `predicted_class`, `confidence_score`, `status`, `disease_info`, `catalog_available`, `image_url`, dan `created_at`. Hasil sehat dan `UNCERTAIN` tidak memiliki informasi penyakit; `UNCERTAIN` disimpan dengan kelas `unknown`. Jika kelas tidak ada di katalog, `catalog_available` bernilai false. Hasil gagal, gambar buram, dan error ML tidak disimpan sebagai pemeriksaan.

Contoh dari direktori `backend/`, setelah API dan ML berjalan dan `ACCESS_TOKEN` diisi:

```sh
curl -X POST http://localhost:8080/api/v1/inspections \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -F 'plant=tomato' -F 'file=@/path/daun.jpg'
```

`GET /api/v1/inspections/:id/image` memerlukan token pemilik. Jangan membuat folder foto menjadi static/public. Endpoint daftar, detail, dan hapus histori belum tersedia.

Error menggunakan `{"error":{"code":"...","message":"..."}}`: input salah 400, foto terlalu besar 413, buram 422 `BLURRY_IMAGE`, ML tidak tersedia 503, timeout ML 504, dan respons ML tidak valid 502. ML belum menyediakan deteksi non-leaf khusus; `UNCERTAIN` tidak menjamin foto bukan daun.

## Konfigurasi

- `ML_SERVICE_URL`: origin ML, default `http://127.0.0.1:8000`; timeout 15 detik. Dari root repository, jalankan ML dengan `Env/bin/python -m uvicorn ml.api.main:app --host 127.0.0.1 --port 8000` jika dependency tersedia.
- `INSPECTION_STORAGE_DIR`: folder privat lokal, default `data/inspections` relatif terhadap direktori kerja backend. Folder default diabaikan Git. Penyimpanan lokal untuk development satu mesin; backup file dan database bersama-sama.
- `AZURE_BLOB_CONTAINER_SAS_URL`: URL HTTPS container Blob privat beserta SAS dengan izin read, create/write, dan delete. Jangan masukkan nilainya ke Git atau log. Saat disetel, gambar disimpan di Azure. Nilai ini wajib saat `GIN_MODE=release`. Batasi izin dan masa berlaku SAS; pastikan container tidak terbuka untuk publik. Integrasi Azure diuji dengan mock, belum dengan akun Azure nyata.

Seed katalog tomat saat ini data dummy dan belum ditinjau untuk digunakan sebagai panduan bagi petani. `go run ./cmd/seed` dari `backend/` menjalankan AutoMigrate sebelum mengisi data.

## Skema dan keterbatasan

Kolom `predicted_class` nullable menyimpan kelas meskipun data katalog belum ada. `image_url` di database menyimpan key internal, bukan URL publik. Untuk database lama tersedia `migrations/003_inspection_predicted_class.sql`; server juga menjalankan AutoMigrate saat startup. Backup dan uji perubahan pada salinan database terlebih dahulu.

Kegagalan insert yang teridentifikasi akan membersihkan foto. Jika hasil penulisan database atau upload Azure tidak dapat dipastikan, file dipertahankan untuk menghindari kehilangan foto dari record yang mungkin sudah tersimpan. Log mencatat key yang perlu diperiksa. Gangguan proses atau kegagalan cleanup dapat meninggalkan file yatim; rekonsiliasi diperlukan sebelum production. Hindari retry otomatis jika respons pemeriksaan hilang karena request belum idempoten.

Tes integrasi pemeriksaan memakai `INSPECTION_TEST_DATABASE_DSN` untuk PostgreSQL khusus tes. Tanpa variabel tersebut tes integrasi dilewati; jangan gunakan database development atau production.
