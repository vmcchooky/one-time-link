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

### Product Documentation
- [docs/README.md](docs/README.md) - Tổng quan tài liệu
- [docs/product-spec/one-time-link-requirements.md](docs/product-spec/one-time-link-requirements.md) - Yêu cầu sản phẩm
- [docs/product-spec/one-time-link-architecture.md](docs/product-spec/one-time-link-architecture.md) - Kiến trúc hệ thống
- [docs/product-spec/one-time-link-milestones.md](docs/product-spec/one-time-link-milestones.md) - Lộ trình phát triển

### Development Documentation
- [DEVELOPMENT.md](DEVELOPMENT.md) - Hướng dẫn phát triển chi tiết
- [docs/contracts/public-http-api.md](docs/contracts/public-http-api.md) - API contract
- [docs/contracts/crypto-and-api-decisions.md](docs/contracts/crypto-and-api-decisions.md) - Crypto specifications

### Deployment Documentation
- [docs/deployment/deployment-guide.md](docs/deployment/deployment-guide.md) - Hướng dẫn deployment
- [docs/deployment/production-checklist.md](docs/deployment/production-checklist.md) - Production checklist

### Milestone Documentation
- [docs/MILESTONE_1_COMPLETION.md](docs/MILESTONE_1_COMPLETION.md) - Milestone 1 completion report
- [docs/MILESTONE_2_COMPLETION.md](docs/MILESTONE_2_COMPLETION.md) - Milestone 2 completion report
- [docs/MILESTONE_2_QUICK_REFERENCE.md](docs/MILESTONE_2_QUICK_REFERENCE.md) - Milestone 2 quick reference

## Milestone Progress

### ✅ Milestone 1: Foundation and Local Development (COMPLETE)

**Completed Features:**
- ✅ Monorepo structure với frontend/backend separation
- ✅ React + TypeScript frontend shell với Vite
- ✅ Go backend với internal service boundaries
- ✅ Docker Compose cho local Redis
- ✅ Health endpoint và basic HTTP routing
- ✅ Structured logging và CORS middleware
- ✅ Request ID tracking
- ✅ Request size limiting (15KB)
- ✅ IP và User-Agent hashing trong logs
- ✅ API contract documentation

**Documentation:** `docs/MILESTONE_1_COMPLETION.md`

### ✅ Milestone 2: Client-Side Encryption and Secret Creation (COMPLETE)

**Completed Features:**

**Backend:**
- ✅ POST /api/secrets endpoint
- ✅ Redis service với TTL auto-expiration
- ✅ Comprehensive validation (algorithm, TTL, nonce, ciphertext)
- ✅ Error handling (400, 413, 500, 201)
- ✅ UUID generation cho secret IDs
- ✅ JSON serialization
- ✅ 81.5% test coverage cho httpapi

**Frontend:**
- ✅ AES-GCM 256-bit encryption với Web Crypto API
- ✅ 12-byte nonce generation
- ✅ Base64url encoding (RFC 4648)
- ✅ Create secret form với TTL selection (1h, 24h, 7d)
- ✅ Plaintext size validation (10KB limit)
- ✅ Byte counter
- ✅ Secret link generation với key trong URL fragment
- ✅ Copy to clipboard functionality
- ✅ Vietnamese UI

**Testing:**
- ✅ 20+ test cases
- ✅ ~70% overall coverage
- ✅ Integration tests
- ✅ Manual test scripts

**Quality:**
- ✅ Code quality: 9.0/10
- ✅ Security: 10/10
- ✅ Specs compliance: 100%

**Documentation:** 
- `docs/MILESTONE_2_COMPLETION.md` - Full completion report
- `docs/MILESTONE_2_REVIEW.md` - Code review and quality assessment
- `docs/MILESTONE_2_QUICK_REFERENCE.md` - Quick reference guide

### ⏳ Milestone 3: Secret Reveal and Consumption (NEXT)

**Planned Features:**
- GET /api/secrets/{id}/status endpoint
- POST /api/secrets/{id}/consume endpoint
- Reveal page component
- Client-side decryption
- Already-used state tracking
- End-to-end testing

## Testing

### Run Backend Tests

```bash
cd backend
go test ./...

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run Integration Tests

```bash
# Requires backend running
cd backend
go test -v ./test
```

### Manual API Testing

```powershell
# PowerShell
.\scripts\test-create-secret.ps1
.\scripts\test-milestone2-comprehensive.ps1
```

```bash
# Bash
./scripts/test-create-secret.sh
```

## Current Features

### Working Features ✅

1. **Create Secret Flow**
   - Enter plaintext (max 10KB)
   - Select TTL (1 hour, 24 hours, 7 days)
   - Client-side AES-GCM encryption
   - Generate shareable link
   - Copy link to clipboard

2. **Backend API**
   - POST /api/secrets - Create encrypted secret
   - GET /healthz - Health check
   - Request validation
   - Error handling
   - Redis storage với TTL

3. **Security**
   - Client-side encryption (plaintext never sent to server)
   - Encryption key in URL fragment (never sent to server)
   - Input validation
   - Size limits (10KB plaintext, 15KB request)
   - CORS protection
   - IP/UA hashing in logs

### Pending Features ⏳

1. **Reveal Flow** (Milestone 3)
   - Secret status checking
   - Reveal gate UI
   - Client-side decryption
   - One-time consumption

2. **Advanced Features** (Future Milestones)
   - Rate limiting
   - Custom expiration times
   - Email delivery
   - Admin dashboard

## Technology Stack

### Backend
- **Language:** Go 1.26
- **Framework:** Standard library (net/http)
- **Database:** Redis 7
- **Dependencies:**
  - github.com/google/uuid - UUID generation
  - github.com/redis/go-redis/v9 - Redis client

### Frontend
- **Framework:** React 19
- **Language:** TypeScript 5.8
- **Build Tool:** Vite 7
- **Crypto:** Web Crypto API (native)

### Infrastructure
- **Development:** Docker Compose
- **Production Target:** VPS + Vercel
- **Reverse Proxy:** Caddy

## Project Structure

```
one-time-link/
├── backend/
│   ├── cmd/api/              # Application entry point
│   ├── internal/
│   │   ├── config/           # Configuration
│   │   ├── httpapi/          # HTTP handlers and middleware
│   │   ├── secret/           # Secret service and types
│   │   └── store/            # Redis client
│   └── test/                 # Integration tests
├── frontend/web-app/
│   └── src/
│       ├── components/       # React components
│       ├── lib/              # Utilities (crypto, api, types)
│       └── styles.css        # Styling
├── deploy/
│   ├── local/                # Docker Compose for local dev
│   └── prod/                 # Production configs
├── docs/
│   ├── contracts/            # API and crypto specs
│   ├── deployment/           # Deployment guides
│   ├── product-spec/         # Product requirements
│   └── MILESTONE_*.md        # Milestone reports
└── scripts/                  # Test and utility scripts
```

## Environment Variables

### Backend (.env)

```bash
APP_SERVICE_NAME=one-time-link-api
APP_HOST=0.0.0.0
APP_PORT=8080
ALLOWED_ORIGIN=http://localhost:5173
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

### Frontend (.env)

```bash
VITE_API_BASE_URL=http://localhost:8080
```

## Contributing

Xem [DEVELOPMENT.md](DEVELOPMENT.md) để biết chi tiết về:
- Development workflow
- Code style guidelines
- Testing requirements
- Commit conventions

## License

[Add license information]

## Contact

[Add contact information]
