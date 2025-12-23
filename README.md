# Store Backend (Gin + Midtrans)

Backend Go untuk toko top-up game dengan Midtrans Snap, invoice PDF, worker email, dan observability bawaan.

## Sorotan
- API RESTful berbasis Gin, JWT untuk admin, CORS siap untuk frontend (`http://localhost:3000` dan `https://whastore.my.id`).
- Postgres + GORM auto migrate, seeding admin awal dan sample game/package.
- Pembayaran Midtrans Snap: generate snap token, verifikasi signature callback, map status transaksi.
- Invoice PDF memakai gofpdf, antrian email via Redis, worker background mengirim email sukses pembayaran dengan lampiran.
- Unggah gambar game ke Cloudinary.
- Monitoring: endpoint `/metrics` Prometheus, dashboard Grafana, cAdvisor, node-exporter melalui Docker Compose.
- Health check `/health` untuk readiness.

## Struktur Cepat
```
.
|- cmd/server/main.go          # bootstrap server + CORS
|- internal/config             # load & validate environment
|- internal/database           # Postgres init + auto migrate + seed
|- internal/routes             # routing publik, admin, callback
|- internal/handlers           # HTTP handler publik/admin/payment
|- internal/middlewares        # JWT auth admin
|- internal/services           # Midtrans, Redis/email queue, invoice PDF, Cloudinary, JWT
|- monitoring/prometheus.yml   # scraping config
`- docker-compose.yml          # app + Postgres + Redis + monitoring stack
```

## Prasyarat
- Docker dan Docker Compose (untuk jalur termudah).
- Atau Go (sesuai `go.mod`) bila ingin jalan langsung.
- Akses ke Postgres, Redis, SMTP, Midtrans, dan Cloudinary.

## Setup Cepat (Docker Compose)
1. Salin `.env.example` menjadi `.env` lalu isi nilai berikut (contoh untuk stack compose):
```
APP_ENV=development
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=storeDB
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
JWT_SECRET=super-secret-jwt
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=you@mail.com
SMTP_PASSWORD=app-password
MIDTRANS_BASE_URL=https://app.sandbox.midtrans.com
MIDTRANS_SERVER_KEY=SB-Mid-server-xxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxx
CLOUDINARY_CLOUD_NAME=your-cloud
CLOUDINARY_API_KEY=your-key
CLOUDINARY_API_SECRET=your-secret
```
2. Jalankan stack: `docker compose up -d`.
3. Aplikasi siap di `http://localhost:8080` (backend), Prometheus `:9090`, Grafana `:3001`, cAdvisor `:8081`, node-exporter `:9100`.
4. Database dan data awal otomatis dimigrasi serta disediakan.

## Jalankan Langsung (tanpa Docker)
1. Siapkan Postgres dan Redis lokal, set env sesuai kebutuhan (bisa reuse `.env`).
2. Instal dependency Go: `go mod download`.
3. Jalankan: `go run ./cmd/server`.
4. Server start di `:8080`, auto migrate + seed berjalan di bootstrap.

## Akun Seed
- Admin default: `admin@mail.com` / `admin123` (ubah segera di production).

## API Kilat
- Public:
  - `GET /health` pengecekan layanan.
  - `GET /api/v1/games` daftar game aktif.
  - `GET /api/v1/games/:id/packages` paket top-up per game.
  - `POST /api/v1/checkout` body `{"game_id":1,"package_id":2,"game_user_id":"123","email":"user@mail.com"}` -> balikan `order_id`, `snap_token`, `redirect_url`.
  - `GET /api/v1/transactions/:order_id` status transaksi.
- Pembayaran:
  - `POST /api/v1/payments/midtrans/callback` dipakai Midtrans, sudah verifikasi signature dan idempotent.
- Admin (Authorization: Bearer <token> dari `/admin/login`):
  - `POST /api/v1/admin/login`
  - `GET /api/v1/admin/transactions?status=PENDING|PAID|FAILED`
  - CRUD game: `GET/POST/PUT/DELETE /api/v1/admin/games` dan `PUT /api/v1/admin/games/:id`
  - CRUD paket: `GET/POST /api/v1/admin/games/:id/packages`, `PUT/DELETE /api/v1/admin/packages/:id`
  - `GET /api/v1/admin/invoices/:order_id` unduh invoice PDF

## Alur Pembayaran
1. Frontend memanggil `POST /api/v1/checkout` untuk mendapat `snap_token`.
2. User bayar di Midtrans; Midtrans memanggil callback backend.
3. Backend verifikasi signature + nominal, update status transaksi dan payment.
4. Jika sukses, backend generate invoice PDF (`invoices/<order_id>.pdf`) dan enqueue email di Redis.
5. Worker background menarik queue dan mengirim email sukses dengan lampiran invoice.

## Monitoring & Observability
- Prometheus scrape `/metrics` (default dari Gin + promhttp).
- Grafana bawaan Compose bisa langsung dipakai untuk membuat dashboard.
- cAdvisor + node-exporter siap untuk metrik container dan host.

## Catatan
- CORS sudah mengizinkan `http://localhost:3000` dan `https://whastore.my.id`.
- Invoice dan file sementara disimpan di disk container (`invoices/`, `tmp_logo.png`).
- Upload gambar game menggunakan Cloudinary folder `games`.
