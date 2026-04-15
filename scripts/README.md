# Scripts

Thư mục này chứa các script hỗ trợ development, testing, và deployment.

## Test Scripts

### test-create-secret.sh / test-create-secret.ps1

Script để test endpoint POST /api/secrets với nhiều test cases khác nhau.

**Bash (Linux/Mac):**
```bash
# Chạy với backend mặc định (localhost:8080)
./scripts/test-create-secret.sh

# Chạy với custom API URL
API_BASE_URL=http://localhost:3000 ./scripts/test-create-secret.sh
```

**PowerShell (Windows):**
```powershell
# Chạy với backend mặc định (localhost:8080)
.\scripts\test-create-secret.ps1

# Chạy với custom API URL
$env:API_BASE_URL="http://localhost:3000"
.\scripts\test-create-secret.ps1
```

**Test Cases:**
1. ✅ Valid request với TTL 1 giờ
2. ✅ Valid request với TTL 24 giờ
3. ❌ Invalid algorithm (expect 400)
4. ❌ Invalid TTL (expect 400)
5. ❌ Empty ciphertext (expect 400)
6. ❌ Invalid nonce length (expect 400)
7. ❌ Payload too large (expect 413)

## Development Workflow

Xem `DEVELOPMENT.md` ở root directory để biết chi tiết về:
- Cách chạy backend
- Cách chạy frontend
- Cách chạy Redis local
- Cách chạy tests

## Future Scripts

Các script sẽ được thêm trong tương lai:
- Local bootstrap scripts
- Build và release automation
- Smoke tests
- Failover drills
