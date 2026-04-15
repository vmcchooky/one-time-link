# one-time-link

`one-time-link` là ứng dụng chia sẻ secret một lần, được thiết kế ưu tiên cho mục đích portfolio với khả năng nâng cấp lên production sau này.

## Trọng Tâm Hiện Tại

**Milestone 2 đã hoàn thành! ✅**

Repository hiện đang ở giai đoạn:

- ✅ **Milestone 1:** Foundation and Local Development - COMPLETE
- ✅ **Milestone 2:** Client-Side Encryption and Secret Creation - COMPLETE
- ⏳ **Milestone 3:** Secret Reveal and Consumption - NEXT

### Milestone 2 Highlights

- ✅ Client-side encryption với AES-GCM 256-bit
- ✅ POST /api/secrets endpoint hoạt động
- ✅ Redis storage với TTL tự động
- ✅ Create secret form với TTL selection
- ✅ URL generation với encryption key trong fragment
- ✅ Comprehensive validation và error handling
- ✅ 20+ test cases với 70% coverage

Xem chi tiết: `docs/MILESTONE_2_COMPLETION.md`

## Cấu Trúc Repository

```text
backend/
  cmd/api/
  internal/
frontend/
  web-app/
deploy/
  local/
  prod/
docs/
scripts/
```

## Phát Triển Local

### 1. Khởi Động Redis

```bash
docker compose -f deploy/local/docker-compose.yml up -d
```

### 2. Chạy Go API

```bash
go run ./backend/cmd/api
```

API sẽ lắng nghe tại `http://localhost:8080` theo mặc định.

### 3. Chạy Frontend

```bash
cd frontend/web-app
npm install
npm run dev
```

Frontend sẽ chạy tại `http://localhost:5173` theo mặc định.

## Tài Liệu Chính

- [docs/README.md](docs/README.md) - Tổng quan tài liệu
- [docs/product-spec/one-time-link-requirements.md](docs/product-spec/one-time-link-requirements.md) - Yêu cầu sản phẩm
- [docs/product-spec/one-time-link-architecture.md](docs/product-spec/one-time-link-architecture.md) - Kiến trúc hệ thống
- [docs/product-spec/one-time-link-milestones.md](docs/product-spec/one-time-link-milestones.md) - Lộ trình phát triển
- [docs/deployment/deployment-guide.md](docs/deployment/deployment-guide.md) - Hướng dẫn deployment
- [DEVELOPMENT.md](DEVELOPMENT.md) - Hướng dẫn phát triển chi tiết

## Milestone 1 Đã Hoàn Thành

- Cấu trúc monorepo
- Shell frontend
- Shell backend
- File Docker Compose cho Redis local
- Tài liệu HTTP contract
- Health endpoint và structured request logging
- Request ID middleware
- CORS middleware
- Request size limiting
- IP và User-Agent hashing trong logs

## Bước Tiếp Theo

Milestone 2 sẽ triển khai luồng tạo secret với client-side encryption và endpoint `POST /api/secrets` thực sự.
