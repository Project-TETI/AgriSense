# AgriSense

## Kelompok Gokils

### Project Senior Project TI

**Instansi:**  
Departemen Teknologi Elektro dan Teknologi Informasi,  
Fakultas Teknik,  
Universitas Gadjah Mada

---

## Anggota Kelompok

| Nama | NIM | Peran |
|---|---|---|
| Nathanael Satya Putra | 24/534424/TK/59236 | Cloud Engineer, AI Engineer |
| Zahra Elfatima | 24/535709/TK/59448 | UI/UX Designer, Software Engineer |
| Ninda Alifa Rachmayanti | 24/545484/TK/60679 | Project Manager |

---

# Modul 1 – Pembentukan Kelompok & Perumusan Masalah

## 1. Nama Produk

**AgriSense**

## 2. Jenis Produk

**Aplikasi web berbasis Agricultural Intelligence** yang memanfaatkan kecerdasan buatan untuk membantu identifikasi awal penyakit tanaman melalui citra daun.

Pada tahap awal, implementasi difokuskan pada **deteksi penyakit daun tanaman tomat** menggunakan dataset publik. 3

## 3. Latar Belakang & Permasalahan

Identifikasi awal penyakit tanaman secara manual masih lambat, subjektif, dan membutuhkan keahlian khusus. Kondisi tersebut membuat pemantauan kesehatan tanaman menjadi kurang optimal, terutama karena keterbatasan tenaga ahli atau penyuluh dibandingkan dengan luas lahan yang harus dipantau.

AgriSense dikembangkan untuk membantu proses identifikasi awal penyakit tanaman secara lebih cepat dan konsisten melalui citra daun. Pada tahap awal, sistem berfokus pada deteksi penyakit daun tomat menggunakan dataset publik. 4

Permasalahan yang dirumuskan dalam pengembangan AgriSense meliputi:

1. Bagaimana membangun model AI yang dapat mengidentifikasi penyakit tanaman berdasarkan citra daun dengan arsitektur yang dapat diterapkan pada berbagai jenis tanaman?
2. Bagaimana menyediakan hasil inferensi AI melalui aplikasi web menggunakan REST API?
3. Bagaimana menyajikan hasil prediksi, confidence score, dan histori pemeriksaan dalam dashboard?
4. Bagaimana melakukan deployment aplikasi dan model menggunakan komputasi awan agar sistem dapat digunakan secara terintegrasi dan skalabel? 5

## 4. Ide Solusi

AgriSense merupakan aplikasi web berbasis Agricultural Intelligence yang dirancang untuk membantu identifikasi awal penyakit tanaman secara cepat melalui foto daun.

Sistem memanfaatkan model kecerdasan buatan berbasis **Convolutional Neural Network (CNN) atau transfer learning** untuk mengklasifikasikan kondisi kesehatan daun, dengan fokus awal pada tanaman tomat agar cakupan pengembangan tetap terukur.

Sistem mengintegrasikan kecerdasan buatan, jaringan komputer, dan komputasi awan dalam arsitektur modular. Proses klasifikasi dijalankan pada backend melalui cloud deployment dan REST API, kemudian hasilnya ditampilkan pada dashboard web. 6

### Rancangan Fitur

- **Deteksi Penyakit Tanaman**  
  Pengguna dapat memfoto penyakit tanaman dan aplikasi menghasilkan hasil deteksi penyakit tanaman.
- **Informasi Penyakit Tanaman**  
  Aplikasi menampilkan informasi mengenai penyakit dan teknik penanganannya.
- **Histori**  
  Pengguna dapat mengakses histori hasil deteksi penyakit tanaman. 7

## 5. Analisis Kompetitor

### Kompetitor 1 — Plantix

**Jenis:** Direct Competitor  
**Produk:** Aplikasi identifikasi penyakit tanaman berbasis foto dan informasi pertanian.  
**Target:** Petani dan pengguna yang membutuhkan identifikasi awal masalah tanaman.

**Kelebihan:**
- Mendukung deteksi berbagai jenis tanaman dan penyakit secara luas.
- Memiliki tampilan aplikasi mobile yang intuitif.
- Mudah digunakan di lapangan.

**Kekurangan:**
- Berfokus pada aplikasi mobile sehingga kurang fleksibel untuk kebutuhan analitik dashboard berbasis web.
- Banyak fitur lanjutan memerlukan koneksi.

**Keunggulan kompetitif:** Plantix menawarkan ekosistem all-in-one dengan database penyakit tanaman yang luas dan komunitas global, tetapi terbatas pada aplikasi mobile dan kurang fleksibel untuk integrasi API. 8

### Kompetitor 2 — Agrio

**Jenis:** Direct Competitor  
**Produk:** Platform pertanian digital dengan identifikasi penyakit tanaman dan fitur monitoring berbasis citra.  
**Target:** Petani, agronomis, dan pengguna profesional di bidang pertanian.

**Kelebihan:**
- Menyediakan pemantauan berbasis citra satelit (NDVI).
- Memiliki sistem peringatan dini penyebaran hama berbasis lokasi geografis.

**Kekurangan:**
- Fitur analitik dan monitoring utama menggunakan layanan freemium/berbayar.
- Ekosistem aplikasi cukup kompleks sehingga membutuhkan learning curve bagi pengguna pemula.

**Keunggulan kompetitif:** Agrio merupakan platform precision agriculture dengan fitur pemantauan satelit dan peringatan dini penyebaran hama, tetapi memiliki layanan berbayar dan relatif kompleks bagi pengguna awam. 9

### Kompetitor 3 — Tomato Disease Detection

**Jenis:** Direct Competitor  
**Produk:** Mobile Apps  
**Target:** Deteksi penyakit tanaman tomat.

**Kelebihan:**
- Dapat dijalankan secara offline pada perangkat mobile.
- Model dapat berjalan dengan cepat pada smartphone.

**Kekurangan:**
- Hanya berfokus pada tanaman tomat.

**Keunggulan kompetitif:** Dapat berjalan secara offline pada platform smartphone dengan resource yang minim. 10
