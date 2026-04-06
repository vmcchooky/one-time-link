# Hướng Dẫn Phát Triển

Tài liệu này cung cấp hướng dẫn nhanh để thiết lập môi trường phát triển local cho dự án `one-time-link`.

## Yêu Cầu Hệ Thống

- **Go**: 1.21 trở lên
- **Node.js**: 18 trở lên
- **Docker**: để chạy Redis local
- **Git**: để quản lý source code

## Thiết Lập Môi Trường Local

### 1. Clone Repository

```bash
git clone https://github.com/vmcchooky/one-time-link.git
cd one-time-link
```

### 2. Khởi Động Redis

```bash
docker compose -f deploy/local/docker-compose.yml up -d
```

Kiểm tra Redis đang chạy:

```bash
docker compose -f deploy/local/docker-compose.yml ps
```

### 3. Chạy Backend API

```bash
# Từ thư mục gốc của repository
go run ./backend/cmd/api
```

Backend sẽ lắng nghe tại `http://localhost:8080`

Kiểm tra health endpoint:

```bash
curl http://localhost:8080/healthz
```

### 4. Chạy Frontend

```bash
cd frontend/web-app
npm install
npm run dev
```

Frontend sẽ chạy tại `http://localhost:5173`

## Cấu Hình Môi Trường

### Backend

Tạo file `.env` trong thư mục `backend/` (tùy chọn):

```bash
cp backend/.env.example backend/.env
```

Các biến môi trường có sẵn:
- `APP_SERVICE_NAME`: Tên service (mặc định: `one-time-link-api`)
- `APP_HOST`: Host để bind (mặc định: `0.0.0.0`)
- `APP_PORT`: Port để lắng nghe (mặc định: `8080`)
- `ALLOWED_ORIGIN`: CORS origin cho phép (mặc định: `http://localhost:5173`)

### Frontend

Tạo file `.env` trong thư mục `frontend/web-app/` (tùy chọn):

```bash
cp frontend/web-app/.env.example frontend/web-app/.env
```

## Chạy Tests

### Backend Tests

```bash
# Chạy tất cả tests
go test ./backend/...

# Chạy tests với coverage
go test -cover ./backend/...

# Chạy tests với output chi tiết
go test -v ./backend/...
```

### Frontend Verification

Frontend test runner has not been added yet in Milestone 1.

Use build verification instead:

```bash
cd frontend/web-app
npm install
npm run build
```

## Kiểm Tra Nhanh

Sau khi khởi động tất cả services, kiểm tra các endpoint sau:

1. **Backend Health**: `http://localhost:8080/healthz`
2. **Frontend**: `http://localhost:5173`
3. **Redis**: `docker compose -f deploy/local/docker-compose.yml exec redis redis-cli ping`

## Dừng Services

### Dừng Backend
Nhấn `Ctrl+C` trong terminal đang chạy backend

### Dừng Frontend
Nhấn `Ctrl+C` trong terminal đang chạy frontend

### Dừng Redis
```bash
docker compose -f deploy/local/docker-compose.yml down
```

## Cấu Trúc Thư Mục

```
one-time-link/
├── backend/           # Go backend API
│   ├── cmd/api/      # Main application entry point
│   └── internal/     # Internal packages
├── frontend/         # React frontend
│   └── web-app/     # Main web application
├── deploy/          # Deployment configurations
│   ├── local/       # Local development (Docker Compose)
│   └── prod/        # Production configurations
└── docs/            # Documentation
```

## Troubleshooting

### Backend không khởi động được

- Kiểm tra port 8080 có bị chiếm không: `lsof -i :8080` (macOS/Linux) hoặc `netstat -ano | findstr :8080` (Windows)
- Kiểm tra Go version: `go version`

### Frontend không khởi động được

- Xóa `node_modules` và cài lại: `rm -rf node_modules && npm install`
- Kiểm tra Node version: `node --version`

### Redis không kết nối được

- Kiểm tra Docker đang chạy: `docker ps`
- Kiểm tra logs: `docker compose -f deploy/local/docker-compose.yml logs redis`
- Restart Redis: `docker compose -f deploy/local/docker-compose.yml restart redis`

## Tài Liệu Bổ Sung

- [README.md](README.md) - Tổng quan dự án
- [docs/README.md](docs/README.md) - Tài liệu chi tiết
- [docs/contracts/public-http-api.md](docs/contracts/public-http-api.md) - API contract
- [docs/product-spec/one-time-link-milestones.md](docs/product-spec/one-time-link-milestones.md) - Lộ trình phát triển

## Đóng Góp

Hiện tại dự án đang ở giai đoạn Milestone 1. Vui lòng tham khảo [docs/product-spec/one-time-link-milestones.md](docs/product-spec/one-time-link-milestones.md) để biết thêm chi tiết về các milestone và nhiệm vụ hiện tại.
