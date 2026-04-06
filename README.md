# one-time-link

`one-time-link` là ứng dụng chia sẻ secret một lần, được thiết kế ưu tiên cho mục đích portfolio với khả năng nâng cấp lên production sau này.

## Trọng Tâm Hiện Tại

Repository hiện đang ở Milestone 1:

- Tài liệu sản phẩm và deployment đã hoàn thiện
- Frontend có shell React + TypeScript
- Backend có shell Go API
- Redis local được cấu hình qua Docker Compose
- Tài liệu deployment hướng tới `quorix.io.vn` với VPS primary chi phí thấp và Oracle Cloud standby

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
