# Article API

Backend test Sharing Vision 2023. Go + MySQL. Frontend ada di [article-frontend](https://github.com/andiahyaul/article-frontend).

## Menjalankan

Prasyarat: Go 1.22+ dan Docker.

```bash
docker compose up -d     # MySQL di port 3307, database `article`
go run ./cmd/server      # API di http://localhost:8080
```

Tabel `posts` dibuat otomatis lewat migrasi (`migrations/`) saat server start. Koneksi dan port bisa diubah lewat env `DB_DSN` dan `ADDR`.

## Endpoint

| Method | URL | |
|---|---|---|
| POST | `/article/` | buat artikel |
| GET | `/article/{limit}/{offset}` | list dengan paging |
| GET | `/article/{id}` | detail |
| PUT, PATCH, POST | `/article/{id}` | ubah |
| DELETE | `/article/{id}` | hapus |

Body create dan update:

```json
{ "title": "", "content": "", "category": "", "status": "publish" }
```

Validasi: title min 20 karakter, content min 200, category min 3, status `publish`, `draft`, atau `thrash`. Jika gagal, API mengembalikan 400 beserta field yang salah.

Catatan:

- POST `/article/{id}` dipakai untuk update. Delete via POST tersedia di `/article/{id}/delete`.
- `GET /article/{limit}/{offset}` menerima `?status=` untuk filter, dan mengirim total baris di header `X-Total-Count`. Dipakai frontend untuk tab dan pagination.

Postman collection: `postman/article-api.postman_collection.json`.
